package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const ticMilli float64 = 100

var (
	xBounds     float64 = 800
	yBounds     float64 = 600
	opposingACC float64 = 5
	accCap      float64 = 100
	velCap      float64 = 100
)

type Plane struct {
	xvel float64
	yvel float64
	xacc float64
	yacc float64
	xpos float64
	ypos float64
	sync.Mutex
}

func (p *Plane) ManagePosition() {
	for {
		time.Sleep(time.Millisecond * time.Duration(ticMilli))
		p.Lock()
		p.updateVelocity()
		p.updatePosition()
		p.Unlock()
		log.Printf("%#v", p)
	}
}

func (p *Plane) updatePosition() {
	xdelta := p.xvel * (ticMilli / 1000)
	ydelta := p.yvel * (ticMilli / 1000)
	p.xpos, p.ypos = p.inbounds(xdelta, ydelta)
}

func (p *Plane) inbounds(xdelta, ydelta float64) (float64, float64) {
	xpos := xdelta + p.xpos
	ypos := ydelta + p.ypos
	if xpos > xBounds {
		xpos -= xBounds
		log.Println("CAPPED")
	}
	if ypos > yBounds {
		ypos -= yBounds
		log.Println("CAPPED")
	}
	return xpos, ypos
}

func (p *Plane) updateVelocity() {
	p.xacc -= opposingACC
	p.yacc -= opposingACC
	if p.xacc < 0 {
		p.xacc = 0
	}
	if p.yacc < 0 {
		p.yacc = 0
	}
	p.xvel += (p.xacc * (ticMilli / 1000))
	p.yvel += (p.yacc * (ticMilli / 1000))
	if p.xvel > velCap {
		p.xvel = velCap
	}
	if p.yvel > velCap {
		p.yvel = velCap
	}
	if p.xvel-opposingACC < 0 {
		p.xvel = 0
	} else {
		p.xvel -= opposingACC
	}
	if p.yvel-opposingACC < 0 {
		p.yvel = 0
	} else {
		p.yvel -= opposingACC
	}
}

func (p *Plane) AccellerateX(dir bool) {
	p.Lock()
	p.xacc += 10
	if p.xacc > accCap {
		p.xacc = accCap
	}
	p.Unlock()
}

func (p *Plane) AccellerateY(dir bool) {
	p.Lock()
	p.yacc += 10
	if p.yacc > accCap {
		p.yacc = accCap
	}
	p.Unlock()
}

func (p *Plane) Draw(s *sdl.Surface) {
	// Draw on the renderer\
	rect := sdl.Rect{int32(p.xpos), int32(p.ypos), 10, 10}
	s.FillRect(&rect, 0xffff0000)
}

func (p Plane) String() string {
	return fmt.Sprintf(`
	Xvel: %f
	Yvel: %f
	XPos: %f
	YPos: %f
	XAcc: %f
	YAcc: %f
	`, p.xvel, p.yvel, p.xpos, p.ypos, p.xacc, p.yacc)
}
