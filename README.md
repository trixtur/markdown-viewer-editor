# Markdown Viewer/Editor

A native desktop markdown viewer and editor application written in Go.

## Features

- Native desktop application (no browser required)
- View markdown files with live rendering
- Edit markdown files with syntax highlighting
- Browse markdown files in a directory
- Split-pane view (editor + preview)
- Auto-save functionality
- Cross-platform (macOS, Windows, Linux)

## Tech Stack

**Built with**:
- **Fyne** - Native Go GUI toolkit
- **goldmark** - Markdown parser
- Pure Go (no web technologies)

## Getting Started

### Prerequisites

- Go 1.21+

### Installation

```bash
cd ~/code/markdown-viewer-editor
go mod init markdown-viewer-editor
go get fyne.io/fyne/v2@latest
go get github.com/yuin/goldmark@latest
go run main.go
```

### Building

```bash
# Build for current platform
go build -o markdown-viewer-editor main.go

# Or use Fyne's packaging tool for better integration
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -os darwin -icon icon.png
```

## Usage

1. Run the application
2. Click "Open Directory" to select a folder containing markdown files
3. Select a file from the left sidebar
4. Edit in the left pane, see live preview in the right pane
5. Changes are auto-saved

## Development

The application uses Fyne's widget system to create a split-pane interface with:
- File browser (left sidebar)
- Markdown editor (center-left)
- Live preview (center-right)

## Screenshot
<img width="1208" height="833" alt="markdown-editor" src="https://github.com/user-attachments/assets/e4e33e90-e380-4a2a-b30c-4164e30dac5f" />


## License

MIT
