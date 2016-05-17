package main

import "testing"

func TestEnterAddsActionFromNormalMode(t *testing.T) {
	state := NewAppState([]string{}, 10)
	state.CurrentMode = normalMode
	store := NewStore()
	KeyPress{Key: "<enter>"}.Apply(state, store.Actions)
	action := <-store.Actions
	changeMode, ok := action.(ChangeMode)
	if !ok || changeMode.Mode != selectCategoryMode {
		t.Fail()
	}
}

func TestSwitchingFocus(t *testing.T) {
	state := NewAppState([]string{}, 10)
	state.CurrentMode = normalMode
	store := NewStore()

	KeyPress{Key: "<backspace>"}.Apply(state, store.Actions)
	action := <-store.Actions
	changeSelection, ok := action.(ChangeSelection)
	if !ok || changeSelection.Direction != left {
		t.Fail()
	}

	KeyPress{Key: "C-l"}.Apply(state, store.Actions)
	action = <-store.Actions
	changeSelection, ok = action.(ChangeSelection)
	if !ok || changeSelection.Direction != right {
		t.Fail()
	}
}

func TestEnterAddsActionFromSelectCategoryMode(t *testing.T) {
	state := NewAppState([]string{}, 10)
	state.CurrentMode = selectCategoryMode
	store := NewStore()
	KeyPress{Key: "<enter>"}.Apply(state, store.Actions)
	action := <-store.Actions
	changeMode, ok := action.(ChangeMode)
	if !ok || changeMode.Mode != normalMode {
		t.Fail()
	}
}

func TestKeyPressAddsTypeKeyInSelectCategoryMode(t *testing.T) {
	state := NewAppState([]string{}, 10)
	state.CurrentMode = selectCategoryMode
	store := NewStore()
	KeyPress{Key: "a"}.Apply(state, store.Actions)
	action := <-store.Actions
	typeKey, ok := action.(TypeKey)
	if !ok || typeKey.Key != "a" {
		t.Fail()
	}
}

func TestKeyPressAddsBackspaceInSelectCategoryMode(t *testing.T) {
	state := NewAppState([]string{}, 10)
	state.CurrentMode = selectCategoryMode
	store := NewStore()
	KeyPress{Key: "C-8"}.Apply(state, store.Actions)
	action := <-store.Actions
	_, ok := action.(Backspace)
	if !ok {
		t.Fail()
	}
}