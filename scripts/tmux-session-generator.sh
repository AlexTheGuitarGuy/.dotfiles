#!/bin/bash

# Execute zoxide query -l | fzf to select a directory
selected_directory=$(zoxide query -l | fzf)

# Extract the last two directories from the selected directory
session_name=$(basename "$(dirname "$selected_directory")")_$(basename "$selected_directory")
echo $session_name

# Check if the session name contains any characters that might cause issues
if [[ "$session_name" =~ [^[:alnum:]_-]+ ]]; then
    echo "Invalid session name. Exiting."
    exit 1
fi

# Check if the tmux session already exists
if tmux has-session -t "$session_name" 2>/dev/null; then
    # Select first window
    tmux select-window -t "$session_name:1"

    # Attach to the session
    tmux attach-session -t "$session_name"
else
    # Change to the selected directory
    cd "$selected_directory"

    # Create a tmux session with the name of the selected directory
    tmux new-session -d -s "$session_name"

    # Create 3 new windows
    tmux new-window -t "$session_name:$i"
    tmux new-window -t "$session_name:$i"
    tmux new-window -t "$session_name:$i"

    # Select first window
    tmux select-window -t "$session_name:1"

    # Attach to the newly created session
    tmux attach-session -t "$session_name"
fi
