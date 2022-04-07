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

// Package clientconn provides client connection implementation.
package clientconn

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/pmezard/go-difflib/difflib"
	"go.uber.org/zap"

	"github.com/FerretDB/FerretDB/internal/handlers/pg"
	"github.com/FerretDB/FerretDB/internal/handlers/pg/pgdb"
	"github.com/FerretDB/FerretDB/internal/handlers/proxy"
	"github.com/FerretDB/FerretDB/internal/handlers/tigris"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/wire"
)

// Mode represents FerretDB mode of operation.
type Mode string

const (
	// NormalMode only handles requests.
	NormalMode Mode = "normal"
	// ProxyMode only proxies requests to another wire protocol compatible service.
	ProxyMode Mode = "proxy"
	// DiffNormalMode both handles requests and proxies them, then logs the diff.
	// Only the FerretDB response is sent to the client.
	DiffNormalMode Mode = "diff-normal"
	// DiffProxyMode both handles requests and proxies them, then logs the diff.
	// Only the proxy response is sent to the client.
	DiffProxyMode Mode = "diff-proxy"
)

// AllModes includes all operation modes, with the first one being the default.
var AllModes = []Mode{NormalMode, ProxyMode, DiffNormalMode, DiffProxyMode}

// conn represents client connection.
type conn struct {
	netConn net.Conn
	mode    Mode
	backend Backend
	l       *zap.SugaredLogger
	pgh     *pg.Handler
	tgh     *tigris.Handler
	proxy   *proxy.Handler
}

// newConnOpts represents newConn options.
type newConnOpts struct {
	netConn   net.Conn
	mode      Mode
	backend   Backend
	l         *zap.Logger
	pgPool    *pgdb.Pool
	pgMetrics *pg.Metrics
	tgConn    *tigris.Client
	tgMetrics *tigris.Metrics
	proxyAddr string
	startTime time.Time
}

// newConn creates a new client connection for given net.Conn.
func newConn(opts *newConnOpts) (*conn, error) {
	prefix := fmt.Sprintf("// %s -> %s ", opts.netConn.RemoteAddr(), opts.netConn.LocalAddr())
	l := opts.l.Named(prefix)

	peerAddr := opts.netConn.RemoteAddr().String()

	var p *proxy.Handler
	if opts.mode != NormalMode {
		var err error
		if p, err = proxy.New(opts.proxyAddr); err != nil {
			return nil, lazyerrors.Error(err)
		}
	}

	switch opts.backend {
	case Postgres:
		pgHandlerOpts := &pg.NewOpts{
			PgPool:    opts.pgPool,
			L:         l,
			PeerAddr:  peerAddr,
			Metrics:   opts.pgMetrics,
			StartTime: opts.startTime,
		}
		return &conn{
			backend: Postgres,
			netConn: opts.netConn,
			mode:    opts.mode,
			l:       l.Sugar(),
			pgh:     pg.New(pgHandlerOpts),
			proxy:   p,
		}, nil

	case Tigris:
		tgHandlerOpts := &tigris.NewOpts{
			Conn:      opts.tgConn,
			L:         l,
			Metrics:   opts.tgMetrics,
			StartTime: opts.startTime,
		}
		return &conn{
			backend: Tigris,
			netConn: opts.netConn,
			mode:    opts.mode,
			l:       l.Sugar(),
			tgh:     tigris.New(tgHandlerOpts),
			proxy:   p,
		}, nil

	default:
		return nil, fmt.Errorf("%s is not supported", opts.backend)
	}
}

// run runs the client connection until ctx is canceled, client disconnects,
// or fatal error or panic is encountered.
//
// The caller is responsible for closing the underlying net.Conn.
func (c *conn) run(ctx context.Context) (err error) {
	done := make(chan struct{})

	// handle ctx cancelation
	go func() {
		select {
		case <-done:
			// nothing, let goroutine exit
		case <-ctx.Done():
			// unblocks ReadMessage below; any non-zero past value will do
			if e := c.netConn.SetDeadline(time.Unix(0, 0)); e != nil {
				c.l.Warnf("Failed to set deadline: %s", e)
			}
		}
	}()

	defer func() {
		if p := recover(); p != nil {
			// Log human-readable stack trace there (included in the error level automatically).
			c.l.DPanicf("%v\n(err = %v)", p, err)
			err = errors.New("panic")
		}

		if err == nil {
			err = ctx.Err()
		}

		// let goroutine above exit
		close(done)
	}()

	bufr := bufio.NewReader(c.netConn)
	bufw := bufio.NewWriter(c.netConn)
	defer func() {
		if e := bufw.Flush(); err == nil {
			err = e
		}

		if c.proxy != nil {
			c.proxy.Close()
		}

		// c.netConn is closed by the caller
	}()

	for {
		var reqHeader *wire.MsgHeader
		var reqBody wire.MsgBody
		reqHeader, reqBody, err = wire.ReadMessage(bufr)
		if err != nil {
			return
		}

		// do not spend time dumping if we are not going to log it
		if c.l.Desugar().Core().Enabled(zap.DebugLevel) {
			c.l.Debugf("Request header: %s", reqHeader)
			c.l.Debugf("Request message:\n%s\n\n\n", reqBody)
		}

		// handle request unless we are in proxy mode
		var resHeader *wire.MsgHeader
		var resBody wire.MsgBody
		var resCloseConn bool
		if c.mode != ProxyMode {
			switch c.backend {
			case Postgres:
				resHeader, resBody, resCloseConn = c.pgh.Handle(ctx, reqHeader, reqBody)
			case Tigris:
				resHeader, resBody, resCloseConn = c.tgh.Handle(ctx, reqHeader, reqBody)
			default:
				panic(c.backend + ": unknown backend")
			}

			// do not spend time dumping if we are not going to log it
			if c.l.Desugar().Core().Enabled(zap.DebugLevel) {
				c.l.Debugf("Response header: %s", resHeader)
				c.l.Debugf("Response message:\n%s\n\n\n", resBody)
			}
		}

		// send request to proxy unless we are in normal mode
		var proxyHeader *wire.MsgHeader
		var proxyBody wire.MsgBody
		if c.mode != NormalMode {
			if c.proxy == nil {
				panic("proxy addr was nil")
			}

			proxyHeader, proxyBody, _ = c.proxy.Handle(ctx, reqHeader, reqBody)

			// do not spend time dumping if we are not going to log it
			if c.l.Desugar().Core().Enabled(zap.DebugLevel) {
				c.l.Debugf("Proxy header: %s", proxyHeader)
				c.l.Debugf("Proxy message:\n%s\n\n\n", proxyBody)
			}
		}

		// diff in diff mode
		if c.mode == DiffNormalMode || c.mode == DiffProxyMode {
			var diffHeader string
			diffHeader, err = difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
				A:        difflib.SplitLines(resHeader.String()),
				FromFile: "res header",
				B:        difflib.SplitLines(proxyHeader.String()),
				ToFile:   "proxy header",
				Context:  1,
			})
			if err != nil {
				return
			}

			var diffBody string
			diffBody, err = difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
				A:        difflib.SplitLines(resBody.String()),
				FromFile: "res body",
				B:        difflib.SplitLines(proxyBody.String()),
				ToFile:   "proxy body",
				Context:  1,
			})
			if err != nil {
				return
			}

			c.l.Infof("Header diff:\n%s\nBody diff:\n%s\n\n", diffHeader, diffBody)
		}

		// replace response with one from proxy in proxy and diff-proxy modes
		if c.mode == ProxyMode || c.mode == DiffProxyMode {
			resHeader = proxyHeader
			resBody = proxyBody
		}

		if resHeader == nil || resBody == nil {
			c.l.Info("no response to send to client")
			return
		}

		if err = wire.WriteMessage(bufw, resHeader, resBody); err != nil {
			return
		}

		if err = bufw.Flush(); err != nil {
			return
		}

		if resCloseConn {
			err = errors.New("fatal error")
			return
		}
	}
}
