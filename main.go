package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/MontFerret/worker/internal/server"
	"github.com/MontFerret/worker/pkg/worker"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/waitfor/pkg/runner"
	waitrunner "github.com/ziflex/waitfor/pkg/runner"
)

var (
	version string

	port = flag.Uint64("port", 8080, "port to listen")

	chromeIP = flag.String("chrome-ip", "127.0.0.1", "Google Chrome remote IP address")

	chromeDebuggingPort = flag.Uint64("chrome-port", 9222, "Google Chrome remote debugging port")

	showVersion = flag.Bool(
		"version",
		false,
		"show REPL version",
	)

	help = flag.Bool(
		"help",
		false,
		"show this list",
	)
)

func main() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
		return
	}

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
		return
	}

	cdp := worker.CDPSettings{
		Host: *chromeIP,
		Port: *chromeDebuggingPort,
	}
	err := waitForChrome(cdp)

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("wait for Chrome")
	}

	srv := server.New(server.Settings{
		CDP: cdp,
	})

	err = srv.Run(*port)

	log.Err(err).
		Timestamp().
		Msg("listen and server")
}

func waitForChrome(cdp worker.CDPSettings) error {
	return waitrunner.Test(context.Background(), []string{
		cdp.BaseURL(),
	}, runner.WithAttempts(10))
}
