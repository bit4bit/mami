package main

import (
	"github.com/nsf/termbox-go"
	"github.com/bit4bit/GAMI"
	"flag"
	"os"
	"fmt"
)

var (
	address = flag.String("server","","Servidor Asterisk ip:puerto")
	amiuser = flag.String("amiuser","admin","AMI user")
	amipass = flag.String("amipass","root","AMI pass")
)


func monitorEvents(g *gami.AMIClient, win *WindowHistory, items []Drawable) {
	
	for {
		select {

		case event := <-g.Events:
			win.AddEvent(*event)

			go func(){draw(items)}()
				
		}
	}
}

func main() {
	flag.Parse()
	fmt.Print("MAMI v0.0 APP para estudiar eventos asterisk\n")
	if *address == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	gami, err := gami.Dial(*address)
	if err != nil {
		fmt.Fprint(os.Stderr, err,"\n")
		os.Exit(1)
	}
	defer func(){ gami.Close()}()

	if err := gami.Login(*amiuser, *amipass); err != nil {
		fmt.Fprint(os.Stderr, err,"\n")
		os.Exit(1)
	}
	gami.Action("Events", map[string]string{"EventMask":"on"})

	termbox.Init()
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	


	winMG := &WindowManager{make(map[string]Windower),""}

	winHistory := NewWindowHistory(0,0, 40, 20)
	
	//muestra historico de eventos en arbol
	winMG.Set("history", winHistory)

	winDesc := NewWindow(40,0, 45, 20)
	winDesc.Title = "Desc Event"
	//al seleccionar evento, la descripcion
	winMG.Set("event", winDesc)

	items := make([]Drawable,1)
	items[0] = winMG

	
	go monitorEvents(gami, winHistory, items)

	
	draw(items)
mainLoop:
	for {
		ev := termbox.PollEvent()
		
		switch  ev.Type {
		case termbox.EventKey:
				
			switch ev.Key {
			case termbox.KeyEsc:
				break mainLoop;
			case termbox.KeyArrowUp:
				winMG.LineUp()
				draw(items)
			case termbox.KeyArrowDown:
				winMG.LineDown()
				draw(items)
			case termbox.KeyPgup:
				winMG.PageUp()
				draw(items)
			case termbox.KeyPgdn:
				winMG.PageDown()
				draw(items)
			}
		case termbox.EventResize:
			draw(items)
		case termbox.EventMouse:
			
			if winHistory.Inside(Position{ev.MouseX, ev.MouseY}) {
				winDesc.Clear()
			}

			if amievent := winHistory.Get(&ev); amievent != nil {
				winDesc.Write([]byte(PrettyEvent(amievent)))
			}
			
			draw(items)

		}
		winMG.Do(&ev)
		
	}
	
}
