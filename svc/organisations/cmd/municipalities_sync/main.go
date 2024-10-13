package municipalitiessync

import (
	"context"
	"fmt"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/municipalities"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type ActionInput struct {
	cmd  *cobra.Command
	args []string
}

type Action struct {
	handler *municipalities.SyncHandler
	sh      fx.Shutdowner
	cmd     *cobra.Command
	args    []string
}

func newAction(
	lc fx.Lifecycle,
	sh fx.Shutdowner,
	handler *municipalities.SyncHandler,
	input *ActionInput,
) *Action {
	act := &Action{
		sh:      sh,
		handler: handler,
		cmd:     input.cmd,
		args:    input.args,
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

	fmt.Println("Syncing municipalities...")
	res, err := act.handler.Handle(municipalities.SyncCommand{
		SourceCsvPath: act.args[0],
	})

	if err != nil {
		fmt.Printf("Error: %s\n\n", err.Error())

		if errGroup, ok := err.(common.ErrGroup); ok {
			for _, err := range errGroup.All() {
				fmt.Println(err.Error())
			}
		}

		act.sh.Shutdown(fx.ExitCode(1))
		return
	}

	fmt.Printf("Successfully synced %d municipalities from %s\n", res.Processed, res.Path)

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
