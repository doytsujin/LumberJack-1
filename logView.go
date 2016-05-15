package main

import ui "github.com/gizak/termui"

// File contains the lines of a given file
type File struct {
	Name   string
	Lines  []string
	Active bool
}

// Display returns a list object representing the file
func (f File) Display(height int) *ui.List {
	list := ui.NewList()
	list.Height = height
	if f.Active {
		list.BorderFg = ui.ColorWhite
	} else {
		list.BorderFg = ui.ColorYellow
	}
	list.BorderLabel = f.Name
	sliceStart := len(f.Lines) - (height - 2)
	if sliceStart < 0 {
		sliceStart = 0
	}
	list.Items = f.Lines[sliceStart:]
	return list
}

// LogViews is a list of Files
type LogViews struct {
	Files []*File
}

// Display returns a Row object representing all of the logViews
func (lv LogViews) Display(height int) *ui.Row {
	listBlocks := []*ui.List{}
	for _, file := range lv.Files {
		logView := file.Display(height)
		logView.BorderLeft = false
		listBlocks = append(listBlocks, logView)
	}
	if len(listBlocks) > 0 {
		listBlocks[0].BorderLeft = true
	}

	logViewColumns := []*ui.Row{}
	numColumnsEach := 6 //numColumns / 1 //len(state.logViews)
	for _, logViewBlock := range listBlocks {
		logViewColumns = append(logViewColumns, ui.NewCol(numColumnsEach, 0, logViewBlock))
	}
	return ui.NewRow(logViewColumns...)
}

// Select the File at index i
func (lv LogViews) Select(i int) {
	if len(lv.Files) <= i || i < 0 {
		return
	}
	for _, file := range lv.Files {
		file.Active = false
	}
	lv.Files[i].Active = true
}
