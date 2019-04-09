package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	evts := getEvents()
	ticker := time.Tick(time.Millisecond * 20)
	p := &Plane{
		xacc: 0,
		yacc: 0,
	}
	go p.ManagePosition()
	go func() {
		for {

			select {
			case event := <-evts:
				switch evt := event.(type) {
				case *sdl.KeyboardEvent:
					log.Println(evt)
					switch evt.Keysym.Scancode {
					case 82:
						log.Println("Up")
						p.AccellerateY()
					case 81:
						log.Println("Down")
						p.AccellerateY()
					case 80:
						p.AccellerateX()
						log.Println("Left")
					case 79:
						log.Println("Right")
						p.AccellerateX()
					default:
						log.Println(evt.Keysym.Scancode)
					}
				case *sdl.QuitEvent:
					println("Quit")
					os.Exit(0)
				}

			}

		}
	}()
	runtime.LockOSThread()
	for range ticker {
		surface.FillRect(nil, 0)
		window.UpdateSurface()
		p.Draw(surface)
		window.UpdateSurface()
	}
}

func getEvents() chan sdl.Event {
	evts := make(chan sdl.Event)
	go func() {
		for {
			select {
			case evts <- sdl.WaitEvent():
			}
		}
	}()
	return evts
}
