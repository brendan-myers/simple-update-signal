package app

import (
	"fmt"

	"go.temporal.io/sdk/workflow"
)

const (
	SignalName = "set-number-signal"
	UpdateName = "set-number-update"
)

type SetFlagInput struct {
	Number int
}

func SimpleWorkflow(ctx workflow.Context) (int, error) {
	logger := workflow.GetLogger(ctx)

	var input SetFlagInput

	// Register update handler - only sets flag if number > 50
	err := workflow.SetUpdateHandlerWithOptions(ctx, UpdateName,
		func(ctx workflow.Context, newInput SetFlagInput) (int, error) {
			input = newInput
			return input.Number, nil
		},
		workflow.UpdateHandlerOptions{
			Validator: func(ctx workflow.Context, input SetFlagInput) error {
				if input.Number < 50 {
					return fmt.Errorf("Number must be >= 50, got %d", input.Number)
				}
				return nil
			},
		},
	)
	if err != nil {
		return -1, err
	}

	// Handle signals in a goroutine
	workflow.Go(ctx, func(ctx workflow.Context) {
		signalChan := workflow.GetSignalChannel(ctx, SignalName)
		for {
			signalChan.Receive(ctx, &input)
		}
	})

	// Wait for flag to be set (either via signal or update)
	workflow.Await(ctx, func() bool { return input.Number >= 50 })

	logger.Info("Flag was set, workflow completing")
	return input.Number, nil
}
