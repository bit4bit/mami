package main

import (
	"github.com/nsf/termbox-go"
	"github.com/bit4bit/GAMI"
	"strconv"
	"bytes"
	"strings"
)

type HistoryEvents struct {
	//eventos derrotero
	Events map[string][]gami.AMIEvent

	CSeq int
}

func (he *HistoryEvents) Add(event gami.AMIEvent) string {
	//key := strconv.Itoa(he.CSeq) + "." + event.Id
	key := "noseq"
	if _, ok := event.Params["Uniqueid"]; ok {
		key = event.Params["Uniqueid"]		
	}
	
	nkey := key + "|" + strconv.Itoa(len(he.Events[key]))
	event.Params["mami-key"] = nkey
	if _, ok := he.Events[key]; ok {
		he.Events[key] = append(he.Events[key], event)
	} else {
		he.Events[key] = make([]gami.AMIEvent,0)
		he.Events[key] = append(he.Events[key], event)
	}


	he.CSeq++
	return nkey
}

//Toma un evento partiendo de un unico identificador
func (he *HistoryEvents) Get(key string) *gami.AMIEvent {
	ids := strings.Split(key,"|")
	if event, ok := he.Events[ids[0]]; ok {

		skey, _ := strconv.Atoi(ids[1])
		return &event[skey]
	}
	return nil
}

type WindowHistory struct {
	win *Window
	

	//Mapea posicion de linea
	//key de  HistoryEvents
	Events map[int]string
	EventLineCount int
	History HistoryEvents
}

func NewWindowHistory(x, y, width, height int) *WindowHistory {
	w := &WindowHistory{
		win: NewWindow(x, y, width, height),
		Events: make(map[int]string),
		EventLineCount: 0,
		History: HistoryEvents{make(map[string][]gami.AMIEvent), 0 },
	}
	w.win.Title = "History"

	return w
}

func (wh *WindowHistory) LineUp() {
	wh.win.LineUp()
}

func (wh *WindowHistory) PageUp() {
	wh.win.PageUp()
}

func (wh *WindowHistory) PageDown() {
	wh.win.PageDown()
}

func (wh *WindowHistory) LineDown() {
	wh.win.LineDown()
}

func (wh *WindowHistory) AddEvent(ev gami.AMIEvent) {

	wh.Events[wh.EventLineCount] = wh.History.Add(ev)

	var show bytes.Buffer
	show.Reset()
	wh.win.Clear()
	bx := 0
	for _, events := range wh.History.Events {

		
		for ix, event := range events {
			wh.Events[bx] = event.Params["mami-key"]
			if _, ok := event.Params["Uniqueid"]; ok {
				if ix == 0 {
					show.WriteString(">")
				}else{
					show.WriteString("\\--")
				}

				show.WriteString("(")
				show.WriteString(event.Params["Uniqueid"])
				show.WriteString(")")
			}
			show.WriteString(event.Id)
			show.WriteString("\n")
			bx++
		}
	}

	wh.win.Write(show.Bytes())

	wh.EventLineCount++
}

func (wh *WindowHistory) Get(ev *termbox.Event) *gami.AMIEvent {
	switch ev.Key {
	case termbox.MouseLeft:
		if wh.win.Inside(Position{ev.MouseX, ev.MouseY}) {
			lineClick := ev.MouseY - 1 + wh.win.OffsetView 
			if _, ok := wh.Events[lineClick]; ok {
				return wh.History.Get(wh.Events[lineClick])
			}
		}

	}
	return nil
}

func (wh *WindowHistory) Draw() {

	wh.win.Draw()
}

func (wh *WindowHistory) DrawBorder(color termbox.Attribute) {
}

func PrettyEvent(ev *gami.AMIEvent) string {
	var pretty bytes.Buffer

	pretty.WriteString(ev.Id)
	pretty.WriteString("\n")
	for k,v := range ev.Params {
		pretty.WriteString("\t")
		pretty.WriteString(k)
		pretty.WriteString("=")
		pretty.WriteString(v)
		pretty.WriteString("\n")
	}

	return pretty.String()
}

func (w *WindowHistory) Inside(pos Position) bool{
	return w.win.Inside(pos)
}
