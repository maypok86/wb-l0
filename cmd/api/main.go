package main

import (
	"context"
	"flag"
	"log"

	"github.com/maypok86/wb-l0/internal/app/api"
	"github.com/maypok86/wb-l0/internal/config"
	"github.com/maypok86/wb-l0/pkg/logger"
)

const versionCommand = "version"

func main() {
	flag.Parse()

	if flag.Arg(0) == versionCommand {
		printVersion()
		return
	}

	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg := config.Get()
	logger.Init(cfg.Logger.Level)

	app, err := api.New(ctx)
	if err != nil {
		return err
	}

	return app.Run()
}
