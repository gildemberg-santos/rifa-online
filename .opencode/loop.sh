#!/usr/bin/env bash
set -euo pipefail

# Loop Engineering Automation — Rifa Online
# Usage: bash .opencode/loop.sh [--dry-run]
#
# Lê docs/TASKS.md, encontra a próxima task pendente,
# e executa via cy-execute-task.

TASKS_FILE="docs/TASKS.md"
TRACKING_DIR=".opencode/tracking"
mkdir -p "$TRACKING_DIR"

echo "=== Loop Engineering: Rifa Online ==="
echo ""

# List tasks from TASKS.md with their status markers
echo "Tasks status:"
echo "-------------"
for task in $(grep -E '^### Task' "$TASKS_FILE" | sed 's/.*Task //;s/ —.*//'); do
  marker_file="$TRACKING_DIR/task-$task.done"
  if [ -f "$marker_file" ]; then
    echo "  [DONE]  Task $task"
  else
    echo "  [PEND]  Task $task"
  fi
done

echo ""
echo "Next: run 'opencode' and use /goal to execute the next pending task"
echo ""
echo "Loop primitives available:"
echo "  /goal \"Task X.Y: <name> is complete and verified\""
echo "  cy-execute-task (via Compozy skills)"
