package main

import (
	"fmt"
	"regexp"

	ui "github.com/gizak/termui"
)

type modifierType string

const (
	filter      modifierType = "filter"
	highlighter modifierType = "highlighter"
)

type modifiers []modifier
type modifier struct {
	buffer
	active bool
	kind   modifierType
	color  string
}

func getModifierSpan(termWidth int) int {
	minWidth := 17
	columns := 12
	return (minWidth / ((termWidth / columns) | 1)) + 1
}

func (state AppState) toggleModifier(selection int) AppState {
	if selection < len(state.modifiers) {
		state.modifiers[selection].active = !state.modifiers[selection].active
	}
	return state
}

func (state AppState) getModifiers() []string {
	var listItems []string
	listItems = append(listItems, "[Highlighters](fg-cyan,fg-underline)")
	addHeader := true
	for i, f := range state.modifiers {
		var attrs, title string
		if f.kind == filter && addHeader {
			listItems = append(listItems, "", "[Filters](fg-cyan,fg-underline)")
			addHeader = false
		}
		if f.active {
			attrs = fmt.Sprintf("fg-black,bg-%s", f.color)
		} else {
			attrs = fmt.Sprintf("fg-%s", f.color)
		}
		if i == state.selectedMod {
			title = fmt.Sprintf("[[%d]](fg-black) [%s](%s,)", i+1, f.text, attrs)
		} else {
			title = fmt.Sprintf("[[%d]](fg-black) [%s](%s,)", i+1, f.text, attrs)
		}
		if i == state.selectedMod && (state.CurrentMode == modifierMode || state.CurrentMode == editModifier) {
			title = fmt.Sprintf("[[%d]](bg-cyan,fg-black) [%s](%s)", i+1, f.text, attrs)
		} else {
			title = fmt.Sprintf("[[%d] %s](%s)", i+1, f.text, attrs)
		}
		if state.CurrentMode == editModifier && i == state.selectedMod {
			title += "_"
		}
		listItems = append(listItems, title)
	}
	return listItems
}

func (mods modifiers) display(state AppState) *ui.Row {
	var listItems []string
	listItems = append(listItems, state.getModifiers()...)

	if state.selectedMod == len(state.modifiers) {
		listItems = append(listItems, "[ + Add Filter](fg-black,bg-green)")
	} else {
		listItems = append(listItems, "[ + Add Filter](fg-yellow)")
	}

	modList := ui.NewList()
	modList.Overflow = "wrap"
	modList.Items = listItems
	modList.Height = logViewHeight(state.termHeight)
	if state.CurrentMode == modifierMode || state.CurrentMode == editModifier {
		modList.BorderFg = ui.ColorGreen
	}
	modSpan := getModifierSpan(state.termWidth)
	return ui.NewCol(modSpan, 0, modList)
}

func (f File) filter(filters modifiers, height int, offSet int) File {
	var filteredLines File
	for i := range f {
		// Go through file in reverse
		line := f[len(f)-i-1]
		matchFilter := false
		for _, filter := range filters {
			if !filter.active {
				continue
			}
			matched, err := regexp.Match(filter.text, []byte(line))
			if err != nil {
				continue
			}
			if matched {
				matchFilter = true
				break
			}
		}
		if matchFilter {
			// Build up file in reverse
			filteredLines = append([]string{line}, filteredLines...)
			if len(filteredLines) == height+offSet {
				break
			}
		}
	}
	return filteredLines
}
