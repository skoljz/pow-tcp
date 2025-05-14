package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/skoljz/pow_tcp_client/internal/client"
	"github.com/skoljz/pow_tcp_client/internal/pow"
)

func RequestCommand(svc client.Client, solver pow.Solver) *cli.Command {
	return &cli.Command{
		Name:    "request",
		Aliases: []string{"req"},
		Usage:   "solve PoW and fetch a quote",
		Action:  createRequestAction(svc, solver),
	}
}

func createRequestAction(svc client.Client, solver pow.Solver) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		conn, err := svc.Connect(ctx.Context)
		if err != nil {
			return err
		}
		defer conn.Close()

		challenge, err := svc.RequestChallenge(ctx.Context, conn)
		if err != nil {
			return err
		}

		nonce, err := solver.Solve(ctx.Context, challenge)
		if err != nil {
			return err
		}

		quote, err := svc.SubmitSolution(ctx.Context, conn, nonce)
		if err != nil {
			return err
		}

		fmt.Println("> quote:", quote)
		return nil
	}
}

func Exists(name string, cmds []*cli.Command) bool {
	for _, cmd := range cmds {
		if cmd.Name == name {
			return true
		}
		for _, alias := range cmd.Aliases {
			if alias == name {
				return true
			}
		}
	}
	return false
}
