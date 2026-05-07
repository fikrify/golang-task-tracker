# Task Tracker CLI

Go implementation of the [Task Tracker](https://roadmap.sh/projects/task-tracker) project from roadmap.sh.

## Usage

```bash
# Add task
./task-cli add "Buy groceries"

# Update task
./task-cli update <id> "Buy groceries and cook dinner"

# Delete task
./task-cli delete <id>

# Mark status
./task-cli mark-todo <id>
./task-cli mark-in-progress <id>
./task-cli mark-done <id>

# List all tasks
./task-cli list

# List by status
./task-cli list todo
./task-cli list in-progress
./task-cli list done
```

## Build

```bash
go build -o task-cli .
```

## Test

```bash
go test ./...
```