package main

import (
//	"fmt"
//	"bytes"
	"github.com/nsf/termbox-go"
//	"os"
//	"io"
//	"bufio"
)

type Position struct {
	X, Y int
}


type Window struct {
	X, Y int
	Width, Height int

	TextColor termbox.Attribute
	BorderColor termbox.Attribute

	Title string

	//la linea de donde se empieza a mostrar el buffer
	OffsetView int

	//Posicion del cursor en la linea donde escribe
	//no interesa la columna ya que la ventana
	//es de visualizacion 
	LineCursor int

	//almacena contenido de la ventana
	Body []byte
}

type Drawable interface {
	Draw()
}

func NewWindow(x, y, width, height int) *Window {
	return &Window{
		X: x,
		Y: y,
		Width: width, 
		Height: height, 
		TextColor: termbox.ColorWhite,
		BorderColor: termbox.ColorWhite,
		Title: "",
		LineCursor: 1,
		OffsetView: 0,
		Body: make([]byte, 0),
	}
}

//Limite de escritura
func (w *Window) LimitWrite() int {
	return	w.Width - 2
}

func (w *Window) LimitLines() int {
	return	w.Height - 1 //-2 -1 borde izquierd -1 borde derecho
}


func (w *Window) PageUp() {
	w.OffsetView -= w.LimitLines() / 3
	if w.OffsetView < 0 {
		w.OffsetView = 0
	}
}

func (w *Window) PageDown() {
	w.OffsetView += w.LimitLines() / 3
	if w.OffsetView >= len(w.Body) {
		w.OffsetView = len(w.Body)
	}

}

func (w *Window) LineDown() {
	w.OffsetView++
	if w.OffsetView >= len(w.Body) {
		w.OffsetView = len(w.Body)
	}
}

func (w *Window) LineUp() {
	w.OffsetView--
	if w.OffsetView < 0 {
		w.OffsetView = 0
	}

	
}

//Ubica cursor al principio de la ventana
func (w *Window) SetCursor() {
	termbox.SetCursor(w.X + 1,w.Y + 1)
}

func (w *Window) Write(data []byte) (int, error) {

	w.Body = append(w.Body, data...)
	return len(data), nil
}

func (w *Window) Clear() {
	w.Body = make([]byte, 0)
	//w.OffsetView = 0
	//w.LineCursor = 0
}

func (w *Window) ClearScreen() {
	for ix := 0; ix < w.Width; ix++ {
		for iy := 0; iy < w.Height; iy++ {
			termbox.SetCell(ix + w.X, iy + w.Y, rune(' '), termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	draw_box(w.X, w.Y, w.Width, w.Height, w.BorderColor)
}

func (w *Window) Draw() {
	draw_box(w.X, w.Y, w.Width, w.Height, w.BorderColor)

	cursor := Position{0,0}

	//winCountChars := w.LimitLines() * w.LimitWrite()


	offsetView := w.OffsetView

	for cx := 0; cx < len(w.Body); cx++ {

		if cursor.X % w.LimitWrite() == 0 {
			cursor.X = 1
			cursor.Y++
		}

		//cortamos al superar lineas
		//y mostramos caracter de que hay mas lineas 
		if cursor.Y % w.LimitLines() == 0 {
			termbox.SetCell(w.X + w.Width ,w.Y + w.Height, rune('/'), termbox.ColorRed, w.BorderColor)
			break
		}
		cursor.X++
		


		c := w.Body[cx]

		if c == '\n' {

			if offsetView > 0 {
				cursor.X=0
				cursor.Y=0
				offsetView--
			}else{
				cursor.Y++
				cursor.X=1
			}
		}
		if offsetView == 0 {
			termbox.SetCell(w.X + cursor.X, w.Y + cursor.Y, rune(c), termbox.ColorWhite, termbox.ColorDefault)
		}

	}

	for cx := range w.Title {
		termbox.SetCell(w.X + w.Width/2 - len(w.Title) + cx ,w.Y, rune(byte(w.Title[cx])), termbox.ColorRed, w.BorderColor)
	}
}

func (w *Window) DrawBorder(color termbox.Attribute) {
	draw_box(w.X, w.Y, w.Width, w.Height, color)
}

func (w *Window) Inside(pos Position) bool{
	return pos.X >= w.X && pos.X <= w.X + w.Width && pos.Y >= w.Y && pos.Y <= w.Y + w.Height 
}

func draw_box(x, y, width, height int, color termbox.Attribute) {
	for ix := 0; ix <= width; ix++ {
		termbox.SetCell(x + ix, y, rune(' '), termbox.ColorDefault, color)
		termbox.SetCell(x + ix, y + height, rune(' '), termbox.ColorDefault, color)
	}

	for iy := 0; iy < height; iy++ {
		termbox.SetCell(x, y+iy, rune(' '), termbox.ColorDefault, color)
		termbox.SetCell(x+width, y+iy, rune(' '), termbox.ColorDefault, color)
	}

}

func draw(draws []Drawable) {


	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for c := range draws {
		if draws[c] != nil {
			draws[c].Draw()
		}
	}

	termbox.Flush()
}
