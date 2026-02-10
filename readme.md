# snip

`snip` is a terminal UI for saving and browsing command snippets.
Built using [BubbleTea](https://github.com/charmbracelet/bubbletea)

<img width="1470" height="855" alt="snipSnapshot" src="https://github.com/user-attachments/assets/980f88d0-17e2-42a0-9e36-41694ae7bc26" />

## Features

- View saved snippets in a keyboard-driven list UI
- Add new commands with name, description, and usage
- Edit existing commands
- Delete commands with confirmation
- Persist snippets to disk automatically


## Keybindings

### List Screen

- `↑` / `↓`: Navigate snippets
    - `j` / `k` also supported
- `a`: Add a new snippet
- `e`: Edit the selected snippet
- `d`: Delete the selected snippet
- `ctrl+c`: Quit

### Add/Edit Screen

- `tab` / `shift+tab`: Move between fields
- `enter`: Advance field, then save when focused on usage
- `esc`: Return to the list screen

### Delete Confirmation

- `y` or `enter`: Confirm delete
- `n` or `esc`: Cancel

## Data Storage

Snippets are stored in:

```text
~/.config/snip/commands.json
```

The file is created automatically if it does not exist.

## Project Structure

- `main.go`: Root Bubble Tea app and screen routing
- `models/`: TUI models for list/add/edit flows
- `internal/handlers/`: Load/save/update/delete helpers
- `ui/`: Styling
