# Task CLI

Task CLI is a powerful command-line task management application written in Go. It helps you organize, track, and manage your tasks efficiently from the terminal.

## Features

- **Simple and intuitive** command-line interface
- **Add tasks** with title, description, project, priority, and due date
- **List tasks** with filtering and sorting options
- **Mark tasks as completed** and track time spent
- **Remove tasks** completely from the system
- **Edit tasks** to update their information
- **Local storage** of tasks in JSON format
- **Color-coded status indicators** for task status (pending, completed, overdue)

## Installation

### Prerequisites

- Go 1.16 or higher

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/kevin7254/task.git
   ```

2. Navigate to the project directory:
   ```bash
   cd task
   ```

3. Build and install the application:

   Option 1: Build manually:
   ```bash
   go build -o task
   ```
   Then, move the binary to your PATH (optional):
   ```bash
   sudo mv task /usr/local/bin/
   ```

   Option 2: Install directly:
   ```bash
   go install .
   ```
   (This installs the binary directly to your `$GOBIN` path, which is usually `~/go/bin`, but you can configure it differently.)

### Difference between `go build` and `go install`

- **`go build`** compiles the code and creates an executable binary in the current directory or the directory specified via the `-o` flag. It allows you to manually manage the binary, such as moving it to your system's PATH.

- **`go install`** compiles the code and installs the resulting binary directly into the `$GOBIN` directory. It is a simpler option when you want the binary installed for immediate use without additional steps to move it.

## Usage

### Adding Tasks

Add a new task with a title:
```bash
task add "Complete project report"
```

Add a task with additional details:
```bash
task add "Force push to prod" --project work --priority 2 --due 2023-12-31 --description "Push the latest changes to production"
```

Options:
- `--description, -d`: Add a detailed description
- `--project, -p`: Assign to a project (default: "work")
- `--priority, -P`: Set priority (1=Low, 2=Medium, 3=High)
- `--due`: Set due date (format: YYYY-MM-DD)

### Listing Tasks

List all incomplete tasks:
```bash
task list
```

List all tasks including completed ones:
```bash
task list --completed
```

Filter tasks by project:
```bash
task list --project personal
```

Sort tasks by priority:
```bash
task list --sort priority
```

View tasks in full detail:
```bash
task list --view full
```

Options:
- `--project, -p`: Filter by project
- `--completed, -c`: Include completed tasks
- `--sort, -s`: Sort by "id", "priority", or "due"
- `--view`: Set view format ("basic" or "full")

### Completing Tasks

Mark a task as completed:
```bash
task do 1
```

Mark multiple tasks as completed:
```bash
task do 1 2 3
```

Track time spent on a task:
```bash
task do 1 --time 30
```

Options:
- `--time, -t`: Time spent on the task in minutes

### Removing Tasks

Remove a task completely:
```bash
task remove 1
```

Remove multiple tasks:
```bash
task remove 1 2 3
```

### Editing Tasks

Edit a task's title:
```bash
task edit 1 --title "New task title"
```

## Task Status Indicators

- ⏳ Pending task
- ✅ Completed task
- ⚠️ Overdue task

## Storage

Tasks are stored in a JSON file located at `~/.task/tasks.json`.

## Roadmap

### Upcoming Features (maybe)

- Group tasks by project or priority
- Enhanced color support
- Undo/redo functionality
- Interactive add/edit mode
- Show specific task details
- Clear all tasks command
- User profiles
- Cloud synchronization
