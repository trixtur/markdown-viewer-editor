package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type MarkdownEditor struct {
	app           fyne.App
	window        fyne.Window
	currentFile   string
	currentDir    string
	editor        *widget.Entry
	preview       *widget.RichText
	previewScroll *container.Scroll
	fileList      *widget.List
	files         []string
	isDirty       bool
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Markdown Viewer/Editor")
	myWindow.Resize(fyne.NewSize(1200, 800))

	editor := &MarkdownEditor{
		app:    myApp,
		window: myWindow,
		files:  []string{},
	}

	editor.setupUI()

	// Check if a file was provided via command line
	if len(os.Args) > 1 {
		filePath := os.Args[1]

		// If it's a directory, load it
		if info, err := os.Stat(filePath); err == nil && info.IsDir() {
			editor.currentDir = filePath
			editor.loadDirectory(filePath)
		} else if err == nil {
			// It's a file, load it directly
			editor.currentDir = filepath.Dir(filePath)
			editor.loadDirectory(editor.currentDir)
			// Find and select the file in the list
			for i, f := range editor.files {
				if f == filePath {
					editor.fileList.Select(i)
					break
				}
			}
		}
	}

	myWindow.ShowAndRun()
}

func (e *MarkdownEditor) setupUI() {
	// Create editor widget with better sizing
	e.editor = widget.NewMultiLineEntry()
	e.editor.SetPlaceHolder("Select a markdown file or create a new one...")
	e.editor.OnChanged = func(content string) {
		e.isDirty = true
		e.updatePreview(content)
	}
	// Make editor text more readable
	e.editor.TextStyle = fyne.TextStyle{Monospace: true}

	// Create preview widget with better spacing
	e.preview = widget.NewRichTextFromMarkdown("")
	e.preview.Wrapping = fyne.TextWrapWord
	e.preview.Truncation = fyne.TextTruncateOff

	// Create file list
	e.fileList = widget.NewList(
		func() int {
			return len(e.files)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			label.SetText(filepath.Base(e.files[i]))
			// Add some padding for better readability
			label.Truncation = fyne.TextTruncateClip
		},
	)
	e.fileList.OnSelected = func(id widget.ListItemID) {
		e.loadFile(e.files[id])
	}

	// Setup Mac-style menu
	e.setupMenu()

	// Create split container with generous padding
	// Double-wrap in padding for more space
	editorPadded := container.NewPadded(container.NewPadded(e.editor))
	editorScroll := container.NewScroll(editorPadded)

	// Add extra padding to preview for better readability
	previewPadded := container.NewPadded(
		container.NewPadded(
			container.NewPadded(e.preview),
		),
	)
	e.previewScroll = container.NewScroll(previewPadded)

	// Main split: Editor on left, Preview on right
	mainSplit := container.NewHSplit(
		editorScroll,
		e.previewScroll,
	)
	mainSplit.SetOffset(0.5)

	// Create layout with file list on left
	fileListHeader := widget.NewLabel("Files")
	fileListHeader.TextStyle = fyne.TextStyle{Bold: true}

	fileListWithPadding := container.NewPadded(e.fileList)
	fileListContainer := container.NewBorder(
		container.NewPadded(fileListHeader),
		nil, nil, nil,
		container.NewScroll(fileListWithPadding),
	)

	appSplit := container.NewHSplit(
		fileListContainer,
		mainSplit,
	)
	appSplit.SetOffset(0.2)

	// Main container with overall padding
	mainContent := container.NewPadded(appSplit)
	e.window.SetContent(mainContent)
}

func (e *MarkdownEditor) setupMenu() {
	// File menu
	openFileItem := fyne.NewMenuItem("Open File...", func() {
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil || file == nil {
				return
			}
			filePath := file.URI().Path()
			e.currentDir = filepath.Dir(filePath)
			e.loadDirectory(e.currentDir)
			for i, f := range e.files {
				if f == filePath {
					e.fileList.Select(i)
					break
				}
			}
		}, e.window)
	})

	openDirItem := fyne.NewMenuItem("Open Directory...", func() {
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil || dir == nil {
				return
			}
			e.currentDir = dir.Path()
			e.loadDirectory(dir.Path())
		}, e.window)
	})

	newFileItem := fyne.NewMenuItem("New File...", func() {
		e.createNewFile()
	})

	saveItem := fyne.NewMenuItem("Save", func() {
		e.saveCurrentFile()
	})

	fileMenu := fyne.NewMenu("File",
		openFileItem,
		openDirItem,
		fyne.NewMenuItemSeparator(),
		newFileItem,
		saveItem,
	)

	// Edit menu with standard shortcuts
	editMenu := fyne.NewMenu("Edit")

	// Main menu
	mainMenu := fyne.NewMainMenu(
		fileMenu,
		editMenu,
	)

	e.window.SetMainMenu(mainMenu)
	e.window.SetMaster()
}

func (e *MarkdownEditor) loadDirectory(dirPath string) {
	files, err := FindMarkdownFiles(dirPath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("error loading directory: %v", err), e.window)
		return
	}

	e.files = files
	e.fileList.Refresh()

	if len(e.files) > 0 {
		e.fileList.Select(0)
	}
}

func (e *MarkdownEditor) loadFile(filePath string) {
	if e.isDirty && e.currentFile != "" {
		dialog.ShowConfirm("Unsaved Changes",
			"You have unsaved changes. Do you want to save before opening another file?",
			func(save bool) {
				if save {
					e.saveCurrentFile()
				}
				e.doLoadFile(filePath)
			}, e.window)
		return
	}
	e.doLoadFile(filePath)
}

func (e *MarkdownEditor) doLoadFile(filePath string) {
	content, err := LoadFileContent(filePath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("error reading file: %v", err), e.window)
		return
	}

	e.currentFile = filePath
	e.editor.SetText(content)
	e.updatePreview(content)
	e.isDirty = false

	e.window.SetTitle(fmt.Sprintf("Markdown Viewer/Editor - %s", filepath.Base(filePath)))
}

func (e *MarkdownEditor) saveCurrentFile() {
	if e.currentFile == "" {
		dialog.ShowInformation("No File", "Please open or create a file first", e.window)
		return
	}

	err := SaveFileContent(e.currentFile, e.editor.Text)
	if err != nil {
		dialog.ShowError(fmt.Errorf("error saving file: %v", err), e.window)
		return
	}

	e.isDirty = false
	dialog.ShowInformation("Saved", fmt.Sprintf("File saved: %s", filepath.Base(e.currentFile)), e.window)
}

func (e *MarkdownEditor) createNewFile() {
	if e.currentDir == "" {
		dialog.ShowError(fmt.Errorf("please open a directory first"), e.window)
		return
	}

	// Create filename entry
	filenameEntry := widget.NewEntry()
	filenameEntry.SetPlaceHolder("filename (without .md extension)")

	// Create form dialog
	formItems := []*widget.FormItem{
		widget.NewFormItem("Filename", filenameEntry),
	}

	dialog.ShowForm("New File", "Create", "Cancel", formItems, func(confirmed bool) {
		if !confirmed {
			return
		}

		filename := filenameEntry.Text
		if filename == "" {
			return
		}

		if !strings.HasSuffix(filename, ".md") {
			filename = filename + ".md"
		}

		newFilePath := filepath.Join(e.currentDir, filename)

		// Create empty file
		err := CreateMarkdownFile(newFilePath)
		if err != nil {
			dialog.ShowError(fmt.Errorf("error creating file: %v", err), e.window)
			return
		}

		// Reload directory and select new file
		e.loadDirectory(e.currentDir)

		// Find and select the new file
		for i, f := range e.files {
			if f == newFilePath {
				e.fileList.Select(i)
				break
			}
		}
	}, e.window)
}

func (e *MarkdownEditor) updatePreview(markdown string) {
	// Parse markdown and create rich text
	// Note: Fyne's RichText widget doesn't support scrolling to anchor links
	// This is a known limitation of the Fyne framework
	// External links (http/https) will open in the browser when clicked
	e.preview.ParseMarkdown(markdown)
	e.preview.Refresh()
}
