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
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "modernc.org/sqlite" // register database/sql driver
	sqlite3 "modernc.org/sqlite"
	sqlite3lib "modernc.org/sqlite/lib"

	"github.com/FerretDB/FerretDB/build/version"
	mysqlpool "github.com/FerretDB/FerretDB/internal/backends/mysql/metadata/pool"
	"github.com/FerretDB/FerretDB/internal/backends/postgresql/metadata/pool"
	"github.com/FerretDB/FerretDB/internal/util/ctxutil"
	"github.com/FerretDB/FerretDB/internal/util/debug"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/logging"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/util/state"
)

var (
	//go:embed error.tmpl
	errorTemplateB []byte

	//go:embed test_user.json
	testUser string

	//go:embed test_system_users.json
	testSystemUsers string

	// Parsed error template.
	errorTemplate = template.Must(template.New("error").Option("missingkey=error").Parse(string(errorTemplateB)))
)

// versionFile contains version information with leading v.
const versionFile = "build/version/version.txt"

// waitForPort waits for the given port to be available until ctx is canceled.
func waitForPort(ctx context.Context, logger *zap.SugaredLogger, port uint16) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	logger.Infof("Waiting for %s to be up...", addr)

	var retry int64
	for ctx.Err() == nil {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			return nil
		}

		logger.Infof("%s: %s", addr, err)

		retry++
		ctxutil.SleepWithJitter(ctx, time.Second, retry)
	}

	return fmt.Errorf("failed to connect to %s", addr)
}

// setupAnyPostgres configures given PostgreSQL.
func setupAnyPostgres(ctx context.Context, logger *zap.SugaredLogger, uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	port, err := strconv.ParseUint(u.Port(), 10, 16)
	if err != nil {
		return err
	}

	if u.User == nil {
		return lazyerrors.New("No username specified")
	}

	if err = waitForPort(ctx, logger, uint16(port)); err != nil {
		return err
	}

	sp, err := state.NewProvider("")
	if err != nil {
		return err
	}

	p, err := pool.New(uri, logger.Desugar(), sp)
	if err != nil {
		return err
	}

	defer p.Close()

	username := u.User.Username()
	password, _ := u.User.Password()

	var retry int64
	for ctx.Err() == nil {
		if _, err = p.Get(username, password); err == nil {
			break
		}

		logger.Infof("%s: %s", uri, err)

		retry++
		ctxutil.SleepWithJitter(ctx, time.Second, retry)
	}

	return ctx.Err()
}

// setupPostgres configures `postgres` container.
//
// It also creates a user in ferretdb and template1 databases so all subsequent
// databases created from template1 have the user.
// See also https://www.postgresql.org/docs/current/manage-ag-templatedbs.html.
func setupPostgres(ctx context.Context, logger *zap.SugaredLogger) error {
	l := logger.Named("postgres")

	// user `username` must exist, but password may be any, even empty
	err := setupAnyPostgres(ctx, l, "postgres://username@127.0.0.1:5432/ferretdb")
	if err != nil {
		return err
	}

	err = setupUserInPostgres(ctx, logger, "postgres://username@127.0.0.1:5432/ferretdb")
	if err != nil {
		return err
	}

	return setupUserInPostgres(ctx, logger, "postgres://username@127.0.0.1:5432/template1")
}

// setupPostgresSecured configures `postgres_secured` container.
//
// It also creates a user in ferretdb and template1 databases so all subsequent
// databases created from template1 have the user.
// See also https://www.postgresql.org/docs/current/manage-ag-templatedbs.html.
func setupPostgresSecured(ctx context.Context, logger *zap.SugaredLogger) error {
	l := logger.Named("postgres_secured")

	err := setupAnyPostgres(ctx, l, "postgres://username:password@127.0.0.1:5433/ferretdb")
	if err != nil {
		return err
	}

	err = setupUserInPostgres(ctx, logger, "postgres://username:password@127.0.0.1:5433/ferretdb")
	if err != nil {
		return err
	}

	return setupUserInPostgres(ctx, logger, "postgres://username:password@127.0.0.1:5433/template1")
}

// setupSQLite configures sqlite local file.
func setupSQLite(ctx context.Context, logger *zap.SugaredLogger) error {
	return setupUserInSQLite(ctx, "file:../../tmp/sqlite/ferretdb.sqlite", "ferretdb")
}

// setupMySQL configures `mysql` container.
func setupMySQL(ctx context.Context, logger *zap.SugaredLogger) error {
	uri := "mysql://root:password@127.0.0.1:3306/ferretdb"

	sp, err := state.NewProvider("")
	if err != nil {
		return err
	}

	if err := waitForPort(ctx, logger.Named("mysql"), 3306); err != nil {
		return err
	}

	p, err := mysqlpool.New(uri, logger.Desugar(), sp)
	if err != nil {
		return err
	}

	defer p.Close()

	var retry int64
	for ctx.Err() == nil {
		db, err := p.Get("root", "password")
		if err == nil {
			if _, err = db.ExecContext(ctx, "GRANT ALL PRIVILEGES ON *.* TO 'username'@'%';"); err != nil {
				return lazyerrors.Error(err)
			}

			break
		}

		logger.Infof("%s: %s", uri, err)

		retry++
		ctxutil.SleepWithJitter(ctx, time.Second, retry)
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}

// setupMongodb configures `mongodb` container.
func setupMongodb(ctx context.Context, logger *zap.SugaredLogger) error {
	if err := waitForPort(ctx, logger.Named("mongodb"), 47017); err != nil {
		return err
	}

	eval := `'rs.initiate({_id: "rs0", members: [{_id: 0, host: "localhost:47017" }]})'`
	args := []string{"compose", "exec", "-T", "mongodb", "mongosh", "--port=47017", "--eval", eval}

	var buf bytes.Buffer
	var retry int64

	for ctx.Err() == nil {
		buf.Reset()

		err := runCommand("docker", args, &buf, logger)
		if err == nil {
			break
		}

		logger.Infof("%s:\n%s", err, buf.String())

		retry++
		ctxutil.SleepWithJitter(ctx, time.Second, retry)
	}

	return ctx.Err()
}

// setupUserInPostgres creates a user with username/password credential in admin database
// with supported mechanisms.
// It creates admin database (PostgreSQL admin schema), if it does not exist.
// It also creates system.users collection, if it does not exist.
//
// Without this, once the first user is created, the authentication fails
// as username/password does not exist in admin.system.users collection.
func setupUserInPostgres(ctx context.Context, logger *zap.SugaredLogger, uri string) error {
	sp, err := state.NewProvider("")
	if err != nil {
		return err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	p, err := pool.New(uri, logger.Desugar(), sp)
	if err != nil {
		return err
	}

	defer p.Close()

	username := u.User.Username()
	password, _ := u.User.Password()

	dbPool, err := p.Get(username, password)
	if err != nil {
		return err
	}

	defer dbPool.Close()

	var pgErr *pgconn.PgError

	q := `CREATE SCHEMA admin`
	if _, err = dbPool.Exec(ctx, q); err != nil && (!errors.As(err, &pgErr) || pgErr.Code != pgerrcode.DuplicateSchema) {
		return err
	}

	if err == nil {
		q = `CREATE TABLE admin._ferretdb_database_metadata (_jsonb jsonb)`
		if _, err = dbPool.Exec(ctx, q); err != nil {
			return err
		}

		q = `CREATE UNIQUE INDEX _ferretdb_database_metadata_id_idx ON admin._ferretdb_database_metadata (((_jsonb->'_id')))`
		if _, err = dbPool.Exec(ctx, q); err != nil {
			return err
		}

		q = `CREATE UNIQUE INDEX _ferretdb_database_metadata_table_idx ON admin._ferretdb_database_metadata (((_jsonb->'table')))`
		if _, err = dbPool.Exec(ctx, q); err != nil {
			return err
		}
	}

	q = `CREATE TABLE admin.system_users_aff2f7ce (_jsonb jsonb)`
	if _, err = dbPool.Exec(ctx, q); err != nil && (!errors.As(err, &pgErr) || pgErr.Code != pgerrcode.DuplicateTable) {
		return err
	}

	if err == nil {
		q = `CREATE UNIQUE INDEX "system_users_aff2f7ce__id__67399184_idx" ON admin.system_users_aff2f7ce (((_jsonb->'_id')))`
		if _, err = dbPool.Exec(ctx, q); err != nil {
			return err
		}

		q = `INSERT INTO admin._ferretdb_database_metadata (_jsonb) VALUES ($1)`
		if _, err = dbPool.Exec(ctx, q, testSystemUsers); err != nil {
			return err
		}
	}

	q = `INSERT INTO admin.system_users_aff2f7ce (_jsonb) VALUES ($1)`

	_, err = dbPool.Exec(ctx, q, testUser)
	if err != nil && (!errors.As(err, &pgErr) || pgErr.Code != pgerrcode.UniqueViolation) {
		return err
	}

	return nil
}

// setupUserInSQLite creates a user with username/password credential in admin database
// with supported mechanisms.
// It creates admin database, if it does not exist.
// It also creates system.users collection, if it does not exist.
//
// Without this, once the first user is created, the authentication fails
// as username/password does not exist in admin.system.users collection.
func setupUserInSQLite(ctx context.Context, uri string, dbName string) error {
	dbURI, err := url.Parse(uri)
	if err != nil {
		return err
	}

	dbURI.Path = path.Join(dbURI.Path, dbName+".sqlite")
	dbURI.Opaque = dbURI.Path
	db, err := sql.Open("sqlite", dbURI.String())
	if err != nil {
		return err
	}

	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	var se *sqlite3.Error

	q := `
		CREATE TABLE "_ferretdb_collections"
		(
			name TEXT NOT NULL UNIQUE CHECK (name != ''),
			table_name TEXT NOT NULL UNIQUE CHECK (table_name != ''),
			settings TEXT NOT NULL CHECK (settings != '')
		) STRICT`
	if _, err = db.ExecContext(ctx, q); err != nil && (!errors.As(err, &se) || !strings.Contains(se.Error(), "already exists")) {
		return err
	}

	q = `CREATE TABLE "system.users_aff2f7ce" (_ferretdb_sjson TEXT NOT NULL CHECK(_ferretdb_sjson != '')) STRICT`
	if _, err = db.ExecContext(ctx, q); err != nil && (!errors.As(err, &se) || !strings.Contains(se.Error(), "already exists")) {
		return err
	}

	if err == nil {
		q = `CREATE UNIQUE INDEX "system.users_aff2f7ce__id_" ON "system.users_aff2f7ce" (_ferretdb_sjson->"_id")`
		if _, err = db.ExecContext(ctx, q); err != nil {
			return err
		}

		q = `INSERT INTO "_ferretdb_collections" (name, table_name, settings) VALUES ('system.users', 'system.users_aff2f7ce',?)`
		arg := `{"uuid":"6f3bdba1-06a5-453b-b71d-c60854a51415","indexes":[{"name":"_id_","key":[{"field":"_id","descending":false}],"unique":true}],"cappedSize":0,"cappedDocuments":0}`
		if _, err = db.ExecContext(ctx, q, arg); err != nil {
			return err
		}
	}

	q = `INSERT INTO "system.users_aff2f7ce" (_ferretdb_sjson) VALUES (?)`

	_, err = db.ExecContext(ctx, q, testUser)
	if err != nil && (!errors.As(err, &se) || se.Code() != sqlite3lib.SQLITE_CONSTRAINT_UNIQUE) {
		return err
	}

	return nil
}

// setupMongodbSecured configures `mongodb_secured` container.
func setupMongodbSecured(ctx context.Context, logger *zap.SugaredLogger) error {
	if err := waitForPort(ctx, logger.Named("mongodb_secured"), 47018); err != nil {
		return err
	}

	eval := `'rs.initiate({_id: "rs0", members: [{_id: 0, host: "localhost:47018" }]})'`
	shell := `mongodb://username:password@127.0.0.1:47018/?tls=true&tlsCertificateKeyFile=/etc/certs/client.pem&tlsCaFile=/etc/certs/rootCA-cert.pem` //nolint:lll // for readability
	args := []string{"compose", "exec", "-T", "mongodb_secured", "mongosh", "--eval", eval, "--shell", shell}

	var buf bytes.Buffer
	var retry int64

	for ctx.Err() == nil {
		buf.Reset()

		err := runCommand("docker", args, &buf, logger)
		if err == nil {
			break
		}

		logger.Infof("%s:\n%s", err, buf.String())

		retry++
		ctxutil.SleepWithJitter(ctx, time.Second, retry)
	}

	return ctx.Err()
}

// setup runs all setup commands.
func setup(ctx context.Context, logger *zap.SugaredLogger) error {
	go debug.RunHandler(ctx, "127.0.0.1:8089", prometheus.DefaultRegisterer, logger.Named("debug").Desugar())

	for _, f := range []func(context.Context, *zap.SugaredLogger) error{
		setupPostgres,
		setupPostgresSecured,
		setupSQLite,
		setupMySQL,
		setupMongodb,
		setupMongodbSecured,
	} {
		if err := f(ctx, logger); err != nil {
			return err
		}
	}

	logger.Info("Done.")
	return nil
}

// runCommand runs command with given arguments.
func runCommand(command string, args []string, stdout io.Writer, logger *zap.SugaredLogger) error {
	bin, err := exec.LookPath(command)
	if err != nil {
		return err
	}
	cmd := exec.Command(bin, args...)
	logger.Debugf("Running %s", strings.Join(cmd.Args, " "))

	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %s", strings.Join(cmd.Args, " "), err)
	}

	return nil
}

// printDiagnosticData prints diagnostic data and error template on stdout.
func printDiagnosticData(w io.Writer, setupError error, logger *zap.SugaredLogger) error {
	_ = runCommand("docker", []string{"compose", "logs"}, w, logger)

	_ = runCommand("docker", []string{"compose", "ps", "--all"}, w, logger)

	_ = runCommand("docker", []string{"stats", "--all", "--no-stream"}, w, logger)

	var buf bytes.Buffer

	var gitVersion string
	if err := runCommand("git", []string{"version"}, &buf, logger); err != nil {
		gitVersion = err.Error()
	} else {
		gitVersion = buf.String()
	}

	buf.Reset()

	var dockerVersion string
	if err := runCommand("docker", []string{"version"}, &buf, logger); err != nil {
		dockerVersion = err.Error()
	} else {
		dockerVersion = buf.String()
	}

	buf.Reset()

	var composeVersion string
	if err := runCommand("docker", []string{"compose", "version"}, &buf, logger); err != nil {
		composeVersion = err.Error()
	} else {
		composeVersion = buf.String()
	}

	info := version.Get()

	return errorTemplate.Execute(w, map[string]any{
		"Error": setupError,

		"GOOS":   runtime.GOOS,
		"GOARCH": runtime.GOARCH,

		"Version":    info.Version,
		"Commit":     info.Commit,
		"Branch":     info.Branch,
		"Dirty":      info.Dirty,
		"Package":    info.Package,
		"DebugBuild": info.DebugBuild,

		"GoVersion":      runtime.Version(),
		"GitVersion":     strings.TrimSpace(gitVersion),
		"DockerVersion":  strings.TrimSpace(dockerVersion),
		"ComposeVersion": strings.TrimSpace(composeVersion),

		"NewIssueURL": "https://github.com/FerretDB/FerretDB/issues/new/choose",
	})
}

// shellMkDir creates all directories from given paths.
func shellMkDir(paths ...string) error {
	var errs error

	for _, path := range paths {
		if err := os.MkdirAll(path, 0o777); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

// shellRmDir removes all directories from given paths.
func shellRmDir(paths ...string) error {
	var errs error

	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

// shellRead will show the content of a file.
func shellRead(w io.Writer, paths ...string) error {
	for _, path := range paths {
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		fmt.Fprint(w, string(b))
	}

	return nil
}

// packageVersion will print out FerretDB's package version (omitting leading v).
func packageVersion(w io.Writer, file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	v := string(b)
	v = strings.TrimPrefix(v, "v")

	_, err = fmt.Fprint(w, v)

	return err
}

// cli struct represents all command-line commands, fields and flags.
// It's used for parsing the user input.
//
//nolint:vet // for readability
var cli struct {
	Debug bool `help:"Enable debug mode."`

	Setup struct{} `cmd:"" help:"Setup development environment."`

	PackageVersion struct{} `cmd:"" help:"Print package version."`

	Shell struct {
		Mkdir struct {
			Paths []string `arg:"" name:"path" help:"Paths to create." type:"path"`
		} `cmd:"" help:"Create directories if they do not already exist."`
		Rmdir struct {
			Paths []string `arg:"" name:"path" help:"Paths to remove." type:"path"`
		} `cmd:"" help:"Remove directories."`
		Read struct {
			Paths []string `arg:"" name:"path" help:"Paths to read." type:"path"`
		} `cmd:"" help:"Read files."`
	} `cmd:""`

	Tests struct {
		Run struct {
			ShardIndex uint   `help:"Shard index, starting from 1."`
			ShardTotal uint   `help:"Total number of shards."`
			Run        string `help:"Run only tests matching the regexp."`
			Skip       string `help:"Skip tests matching the regexp."`

			Args []string `arg:"" help:"Other arguments and flags for 'go test'." passthrough:""`
		} `cmd:"" help:"Run tests."`
	} `cmd:""`

	Fuzz struct {
		Corpus struct {
			Src string `arg:"" help:"Source, one of: 'seed', 'generated', or collected corpus' directory."`
			Dst string `arg:"" help:"Destination, one of: 'seed', 'generated', or collected corpus' directory."`
		} `cmd:"" help:"Sync fuzz corpora."`
	} `cmd:""`
}

// makeLogger returns a human-friendly logger.
func makeLogger(level zapcore.Level, output []string) (*zap.Logger, error) {
	start := time.Now()

	return zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       true,
		DisableCaller:     true,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "M",
			LevelKey:      zapcore.OmitKey,
			TimeKey:       "T",
			NameKey:       "N",
			CallerKey:     zapcore.OmitKey,
			FunctionKey:   zapcore.OmitKey,
			StacktraceKey: zapcore.OmitKey,
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(fmt.Sprintf("%7.2fs", t.Sub(start).Seconds()))
			},
			EncodeDuration:      zapcore.StringDurationEncoder,
			EncodeCaller:        zapcore.ShortCallerEncoder,
			EncodeName:          nil,
			NewReflectedEncoder: nil,
			ConsoleSeparator:    "  ",
		},
		OutputPaths:      output,
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    nil,
	}.Build()
}

func main() {
	kongCtx := kong.Parse(&cli)

	// https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables
	if t, _ := strconv.ParseBool(os.Getenv("CI")); t {
		cli.Debug = true
	}

	level := zap.InfoLevel
	if cli.Debug {
		level = zap.DebugLevel
	}

	logging.SetupWithZapLogger(must.NotFail(makeLogger(level, []string{"stderr"})))

	logger := zap.S()

	cmd := kongCtx.Command()
	logger.Debugf("Command: %q", cmd)

	var err error

	switch cmd {
	case "setup":
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		err = setup(ctx, logger)

	case "package-version":
		err = packageVersion(os.Stdout, versionFile)

	case "shell mkdir <path>":
		err = shellMkDir(cli.Shell.Mkdir.Paths...)
	case "shell rmdir <path>":
		err = shellRmDir(cli.Shell.Rmdir.Paths...)
	case "shell read <path>":
		err = shellRead(os.Stdout, cli.Shell.Read.Paths...)

	case "tests run <args>":
		ctx, stop := ctxutil.SigTerm(context.Background())
		defer stop()

		err = testsRun(
			ctx,
			cli.Tests.Run.ShardIndex, cli.Tests.Run.ShardTotal,
			cli.Tests.Run.Run, cli.Tests.Run.Skip, cli.Tests.Run.Args,
			logger,
		)

	case "fuzz corpus <src> <dst>":
		var seedCorpus, generatedCorpus string

		if seedCorpus, err = os.Getwd(); err != nil {
			logger.Fatal(err)
		}

		if generatedCorpus, err = fuzzGeneratedCorpus(); err != nil {
			logger.Fatal(err)
		}

		var src, dst string

		switch cli.Fuzz.Corpus.Src {
		case "seed":
			src = seedCorpus
		case "generated":
			src = generatedCorpus
		default:
			if src, err = filepath.Abs(cli.Fuzz.Corpus.Src); err != nil {
				logger.Fatal(err)
			}
		}

		switch cli.Fuzz.Corpus.Dst {
		case "seed":
			// Because we would need to add `/testdata/fuzz` back, and that's not very easy.
			logger.Fatal("Copying to seed corpus is not supported.")
		case "generated":
			dst = generatedCorpus
		default:
			dst, err = filepath.Abs(cli.Fuzz.Corpus.Dst)
			if err != nil {
				logger.Fatal(err)
			}
		}

		err = fuzzCopyCorpus(src, dst, logger)

	default:
		err = fmt.Errorf("unknown command: %s", cmd)
	}

	if err != nil {
		if cmd == "setup" {
			_ = printDiagnosticData(os.Stderr, err, logger)
		}

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			logger.Error(exitErr)
			os.Exit(exitErr.ExitCode())
		}

		logger.Fatal(err)
	}
}
