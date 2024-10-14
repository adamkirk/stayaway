package dbmigrate

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type Migrator interface {
	Up(to string) error
	Down(to string) error
}

type Action struct {
	migrator Migrator
	sh       fx.Shutdowner
	cmd      *cobra.Command
	args     []string
}

type ActionInput struct {
	cmd  *cobra.Command
	args []string
}

func newAction(
	lc fx.Lifecycle,
	sh fx.Shutdowner,
	migrator Migrator,
	input *ActionInput,
) *Action {
	act := &Action{
		sh:       sh,
		migrator: migrator,
		cmd:      input.cmd,
		args:     input.args,
	}

	lc.Append(fx.Hook{
		OnStart: act.start,
		OnStop:  act.stop,
	})

	return act
}

func (act *Action) start(ctx context.Context) error {
	go act.run()
	return nil
}

func (act *Action) stop(ctx context.Context) error {
	return nil
}

func (act *Action) run() {

	var err error
	var target = ""

	if len(act.args) > 1 {
		target = act.args[1]
	}

	if act.args[0] == "down" {
		err = act.migrator.Down(target)
	} else {
		err = act.migrator.Up(target)
	}

	if err != nil {
		slog.Error("failed to migrate db", "error", err)
		act.sh.Shutdown(fx.ExitCode(1))
		return
	}

	fmt.Println("Successfully migrated!")
	act.sh.Shutdown()
}

func Handler(opts []fx.Option, cmd *cobra.Command, args []string) {
	opts = append(opts, []fx.Option{
		// Prevents all the logging noise when building the service container
		fx.NopLogger,
		fx.Provide(func() *ActionInput {
			return &ActionInput{
				cmd:  cmd,
				args: args,
			}
		}),
		fx.Provide(newAction),
		fx.Invoke(func(*Action) {}),
	}...)

	fx.New(
		opts...,
	).Run()
}
