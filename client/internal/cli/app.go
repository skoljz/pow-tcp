package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/skoljz/pow_tcp_client/internal/client"
	"github.com/skoljz/pow_tcp_client/internal/config"
	"github.com/skoljz/pow_tcp_client/internal/pow"
)

const (
	AppName = "pow-cli"
	Usage   = "Proof-of-Work TCP client (enter q or quit to exit)"
)

func New(
	cfg *config.Config,
	svc client.Client,
	solver pow.Solver,
) *cli.App {
	requestCmd := RequestCommand(svc, solver)
	commands := []*cli.Command{requestCmd}

	app := &cli.App{
		Name:     AppName,
		Usage:    Usage,
		Commands: commands,
		Action:   createREPL(AppName, commands),
	}

	return app
}

func createREPL(appName string, commands []*cli.Command) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Print("pow-cli > ")
			line, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			input := strings.TrimSpace(line)

			if input == "" {
				continue
			}
			if input == "q" || input == "quit" {
				return nil
			}

			parts := strings.Fields(input)
			name := parts[0]
			if !Exists(name, commands) {
				fmt.Fprintf(os.Stderr, "command not found: %q\n", name)
				continue
			}

			args := append([]string{appName}, parts...)
			if err := ctx.App.Run(args); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
		}
	}
}
