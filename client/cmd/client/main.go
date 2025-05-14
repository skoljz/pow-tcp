package main

import (
	"log"
	"os"

	"github.com/skoljz/pow_tcp_client/internal/cli"
	"github.com/skoljz/pow_tcp_client/internal/client"
	"github.com/skoljz/pow_tcp_client/internal/config"
	"github.com/skoljz/pow_tcp_client/internal/pow"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("pow-cli > failed to init config: %v", err)
	}

	powClient := client.New(cfg)
	solver := pow.NewSolver(cfg.TargetSize)

	app := cli.New(cfg, powClient, solver)
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("pow-cli > failed to start pow-cli app: %v", err)
	}
}
