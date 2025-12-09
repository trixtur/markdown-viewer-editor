# Usage Guide

## Command Line Usage

### Launch with a specific file

```bash
cd ~/code/markdown-viewer-editor
./markdown-viewer-editor ~/code/KSL_COOKIE_INVENTORY.md
```

This will:
1. Open the application
2. Load all markdown files in `~/code/` directory
3. Automatically select and display `KSL_COOKIE_INVENTORY.md`

### Launch with a directory

```bash
./markdown-viewer-editor ~/code/
```

This will open the application and load all markdown files from the specified directory.

### Launch without arguments

```bash
./markdown-viewer-editor
```

Opens the application with an empty state. Use **File → Open Directory** to browse files.

## Menu Shortcuts (Mac)

### File Menu
- **⌘O** - Open File
- **⌘N** - New File
- **⌘S** - Save

### Edit Menu
- **⌘X** - Cut
- **⌘C** - Copy
- **⌘V** - Paste
- **⌘A** - Select All

## Interface Layout

```
┌─────────────────────────────────────────────────────────┐
│  File  Edit                                             │
├──────────┬─────────────────────┬─────────────────────────┤
│          │                     │                         │
│  Files   │   Editor Pane       │   Preview Pane          │
│  ------  │   (Edit here)       │   (Live preview)        │
│          │                     │                         │
│  • file1 │                     │                         │
│  • file2 │                     │                         │
│  • file3 │                     │                         │
│          │                     │                         │
└──────────┴─────────────────────┴─────────────────────────┘
   20%              40%                   40%
```

## Tips

1. **Auto-save prompt**: If you have unsaved changes and try to open another file, you'll be prompted to save first
2. **File list**: Shows all `.md` files in the opened directory
3. **Live preview**: Updates automatically as you type in the editor
4. **Keyboard shortcuts**: Standard Mac shortcuts work throughout the app
