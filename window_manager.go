package main

import (
	"errors"
	"github.com/nsf/termbox-go"
)

type Windower interface {
	PageUp()
	PageDown()
	LineUp()
	LineDown()
	Draw()
	DrawBorder(color termbox.Attribute)
	Inside(pos Position) bool
}

//Gestion entrada de varias ventadas
//permite seleccionar con los numeros del teclada
//y desplazar la ventana seleccionada
type WindowManager struct {
	windows map[string]Windower
	selectedWindow string
	
}

func (wm *WindowManager) Set(key string, win Windower) {
	wm.windows[key] = win
}

func (wm *WindowManager) Select(key string) error {
	if _, ok := wm.windows[key]; ok {
		wm.selectedWindow = key
		return nil
	}
	return errors.New("Not find window")
}

func (wm *WindowManager) Do(ev *termbox.Event) {
	switch ev.Type {
	case termbox.EventMouse:
		if ev.Key == termbox.MouseLeft {
			for key, win := range wm.windows {
				if win.Inside(Position{ev.MouseX, ev.MouseY}) {
					wm.selectedWindow = key
					wm.Draw()
				}
			}
		}
	}
}

func (wm *WindowManager) PageDown() {
	if _, ok := wm.windows[wm.selectedWindow]; ok {
		wm.windows[wm.selectedWindow].PageDown()
	}
}

func (wm *WindowManager) PageUp() {
	if _, ok := wm.windows[wm.selectedWindow]; ok {
		wm.windows[wm.selectedWindow].PageUp()
	}
}

func (wm *WindowManager) LineUp() {
	if _, ok := wm.windows[wm.selectedWindow]; ok {
		wm.windows[wm.selectedWindow].LineUp()
	}
}

func (wm *WindowManager) LineDown() {
	if _, ok := wm.windows[wm.selectedWindow]; ok {
		wm.windows[wm.selectedWindow].LineDown()
	}
}

func (wm *WindowManager) Draw() {
	for _, win := range wm.windows {
		win.Draw()
	}
	
	if win, ok := wm.windows[wm.selectedWindow]; ok {
		win.DrawBorder(termbox.ColorGreen)
	}
}

func (wm *WindowManager) DrawBorder(color termbox.Attribute) {
}

func (w *WindowManager) Inside(pos Position) bool{
	return false
}
