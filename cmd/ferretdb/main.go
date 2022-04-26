// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"go.uber.org/zap"

	"github.com/FerretDB/FerretDB/internal/clientconn"
	"github.com/FerretDB/FerretDB/internal/handlers/common"
	"github.com/FerretDB/FerretDB/internal/handlers/dummy"
	"github.com/FerretDB/FerretDB/internal/handlers/pg"
	"github.com/FerretDB/FerretDB/internal/handlers/pg/pgdb"
	"github.com/FerretDB/FerretDB/internal/util/debug"
	"github.com/FerretDB/FerretDB/internal/util/logging"
	"github.com/FerretDB/FerretDB/internal/util/version"
)

//nolint:gochecknoglobals // flags are defined there to be visible in `bin/ferretdb-testcover -h` output
var (
	debugAddrF       = flag.String("debug-addr", "127.0.0.1:8088", "debug address")
	listenAddrF      = flag.String("listen-addr", "127.0.0.1:27017", "listen address")
	modeF            = flag.String("mode", string(clientconn.AllModes[0]), fmt.Sprintf("operation mode: %v", clientconn.AllModes))
	postgresqlURLF   = flag.String("postgresql-url", "postgres://postgres@127.0.0.1:5432/ferretdb", "PostgreSQL URL")
	proxyAddrF       = flag.String("proxy-addr", "127.0.0.1:37017", "")
	versionF         = flag.Bool("version", false, "print version to stdout (full version, commit, branch, dirty flag) and exit")
	testConnTimeoutF = flag.Duration("test-conn-timeout", 0, "test: set connection timeout")
	handlerF         = flag.String("handler", "pg", "handler: can be pg or dummy")
)

func main() {
	logging.Setup(zap.DebugLevel)
	logger := zap.L()
	flag.Parse()

	info := version.Get()

	if *versionF {
		fmt.Fprintln(os.Stdout, "version:", info.Version)
		fmt.Fprintln(os.Stdout, "commit:", info.Commit)
		fmt.Fprintln(os.Stdout, "branch:", info.Branch)
		fmt.Fprintln(os.Stdout, "dirty:", info.Dirty)
		return
	}

	logger.Info(
		"Starting FerretDB "+info.Version+"...",
		zap.String("version", info.Version),
		zap.String("commit", info.Commit),
		zap.String("branch", info.Branch),
		zap.Bool("dirty", info.Dirty),
	)

	var found bool
	for _, m := range clientconn.AllModes {
		if *modeF == string(m) {
			found = true
			break
		}
	}
	if !found {
		logger.Sugar().Fatalf("Unknown mode %q.", *modeF)
	}

	ctx, stop := notifyAppTermination(context.Background())
	go func() {
		<-ctx.Done()
		logger.Info("Stopping...")
		stop()
	}()

	go debug.RunHandler(ctx, *debugAddrF, logger.Named("debug"))

	var h common.Handler

	startTime := time.Now()
	switch *handlerF {
	case "pg":
		pgPool, err := pgdb.NewPool(ctx, *postgresqlURLF, logger, false)
		if err != nil {
			logger.Fatal(err.Error())
		}
		handlerOpts := &pg.NewOpts{
			PgPool:    pgPool,
			L:         logger,
			StartTime: startTime,
		}
		h = pg.New(handlerOpts)

	case "dummy":
		h = dummy.New()

	default:
		panic("unknown handler")
	}

	defer h.Close()

	l := clientconn.NewListener(&clientconn.NewListenerOpts{
		ListenAddr:      *listenAddrF,
		ProxyAddr:       *proxyAddrF,
		Mode:            clientconn.Mode(*modeF),
		Handler:         h,
		Logger:          logger,
		TestConnTimeout: *testConnTimeoutF,
	})

	prometheus.DefaultRegisterer.MustRegister(l)

	var err error
	err = l.Run(ctx)
	if err == nil || err == context.Canceled {
		logger.Info("Listener stopped")
	} else {
		logger.Error("Listener stopped", zap.Error(err))
	}

	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		panic(err)
	}
	for _, mf := range mfs {
		if _, err := expfmt.MetricFamilyToText(os.Stderr, mf); err != nil {
			panic(err)
		}
	}
}
