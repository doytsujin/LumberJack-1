package main

import tail "github.com/hpcloud/tail"

// Files list
type Files map[string]File

// File contains the lines of a given file
type File []string

func (state AppState) getSelectedFileName() string {
	return state.LogViews[state.selected].FileName
}

func (state AppState) getSelectedView() LogView {
	return state.LogViews[state.selected]
}

func (state AppState) getSelectedFile() File {
	return state.Files[state.getSelectedFileName()]
}

func (state AppState) getFile(fileName string) File {
	return state.Files[fileName]
}

// WatchFile Action
type WatchFile struct {
	FileName string
}

// Apply WatchFile
func watchFile(fileName string, actions chan<- Action) {
	addNewLine := func(fileName string, newLine string) {
		actions <- AppendLine{FileName: fileName, Line: newLine}
	}
	go tailFile(fileName, addNewLine)
}

func addWatchers(fileNames []string, actions chan<- Action) {
	for _, fileName := range fileNames {
		watchFile(fileName, actions)
	}
}

// AppendLine to file
type AppendLine struct {
	FileName string
	Line     string
}

// Apply the AppendLine
func (action AppendLine) Apply(state AppState) AppState {
	file := state.Files[action.FileName]
	file = append(file, action.Line)
	state.Files[action.FileName] = file
	return state
}

func tailFile(fileName string, callback func(string, string)) {
	t, err := tail.TailFile(fileName, tail.Config{
		Follow:    true,
		Logger:    tail.DiscardingLogger,
		MustExist: true,
	})
	if err != nil {
		panic(err)
	}
	for line := range t.Lines {
		callback(fileName, line.Text)
	}
}
