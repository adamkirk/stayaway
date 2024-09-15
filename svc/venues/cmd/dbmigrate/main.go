package dbmigrate

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/adamkirk-stayaway/venues/internal/db"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type Action struct {
	migrator db.Migrator
	sh fx.Shutdowner
}

func newAction(
	lc fx.Lifecycle,
	sh fx.Shutdowner,
	migrator db.Migrator,
) *Action {
	act := &Action{
		sh: sh,
		migrator: migrator,
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
	err := act.migrator.Migrate()
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
		fx.Provide(newAction),
		fx.Invoke(func(*Action) {}),
	}...)

	fx.New(
		opts...,
	  ).Run()
}
