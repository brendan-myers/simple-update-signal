package main

import (
	"context"
	"flag"
	"log"

	"go.temporal.io/sdk/client"

	"simple-update-signal/app"
)

func main() {
	// Parse command line flags
	action := flag.String("action", "start", "Action: start, signal, or update")
	workflowID := flag.String("workflow-id", "simple-workflow", "Workflow ID")
	number := flag.Int("number", 0, "Number to send with signal or update")
	flag.Parse()

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	ctx := context.Background()

	switch *action {
	case "start":
		startWorkflow(ctx, c, *workflowID)
	case "signal":
		sendSignal(ctx, c, *workflowID, *number)
	case "update":
		sendUpdate(ctx, c, *workflowID, *number)
	default:
		log.Fatalf("Unknown action: %s. Use 'start', 'signal', or 'update'", *action)
	}
}

func startWorkflow(ctx context.Context, c client.Client, workflowID string) {
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "simple-signal-update-task-queue",
	}

	we, err := c.ExecuteWorkflow(ctx, options, app.SimpleWorkflow)
	if err != nil {
		log.Fatalln("Unable to start workflow:", err)
	}

	log.Printf("Started workflow: WorkflowID=%s, RunID=%s\n", we.GetID(), we.GetRunID())
	log.Println("Workflow is waiting for a signal or update to set the flag...")
}

func sendSignal(ctx context.Context, c client.Client, workflowID string, number int) {
	err := c.SignalWorkflow(ctx, workflowID, "", app.SignalName, app.SetFlagInput{Number: number})
	if err != nil {
		log.Fatalln("Unable to signal workflow:", err)
	}
	log.Printf("Sent signal with number %d\n", number)
}

func sendUpdate(ctx context.Context, c client.Client, workflowID string, number int) {
	handle, err := c.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
		WorkflowID:   workflowID,
		UpdateName:   app.UpdateName,
		Args:         []interface{}{app.SetFlagInput{Number: number}},
		WaitForStage: client.WorkflowUpdateStageCompleted,
	})
	if err != nil {
		log.Fatalln("Unable to send update:", err)
	}

	var result string
	err = handle.Get(ctx, &result)
	if err != nil {
		log.Fatalln("Update failed:", err)
	}
	log.Printf("Update succeeded: %s\n", result)
}
