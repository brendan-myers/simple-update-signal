# Simple Update Signal

A simple project to demonstrate the difference between `update` and `signal` in Temporal workflows.

## Key Difference

- **Signal**: Fire-and-forget. The caller doesn't know if the workflow accepted or rejected the value.
- **Update**: Has validation and returns a result. The caller knows immediately if the value was rejected.

In this demo, the workflow only completes when it receives a number >= 50. With signals, invalid values are silently ignored. With updates, invalid values return an error to the caller.

## Usage

### Prerequisites

- A running Temporal server (e.g., `temporal server start-dev`)

### 1. Start the worker

```bash
go run worker/main.go
```

### 2. Start the workflow

```bash
go run starter/main.go -action start
```

### 3. Try sending values

**With Signal (no validation feedback):**
```bash
# This is accepted but workflow doesn't complete (number < 50)
go run starter/main.go -action signal -number 25

# This completes the workflow
go run starter/main.go -action signal -number 75
```

**With Update (validation feedback):**
```bash
# This returns an error: "Number must be >= 50"
go run starter/main.go -action update -number 25

# This succeeds and completes the workflow
go run starter/main.go -action update -number 75
```
