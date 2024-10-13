package dbping

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/adamkirk-stayaway/organisations/internal/db"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type Action struct {
	pinger db.Pinger
	sh     fx.Shutdowner
}

func newAction(
	lc fx.Lifecycle,
	sh fx.Shutdowner,
	pinger db.Pinger,
) *Action {
	act := &Action{
		sh:     sh,
		pinger: pinger,
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
	err := act.pinger.Ping()
	if err != nil {
		slog.Error("failed to ping db", "error", err)
		act.sh.Shutdown(fx.ExitCode(1))
		return
	}

	fmt.Println("Successfully pinged DB!")
	act.sh.Shutdown()
}

func Handler(opts []fx.Option, cmd *cobra.Command, args []string) {
	opts = append(opts, []fx.Option{
		// Prevents all the logging noise when building the service container
		fx.NopLogger,
		fx.Provide(newAction),
		fx.Invoke(func(*Action) {}),
	}...)

	fx.New(
		opts...,
	).Run()
}
