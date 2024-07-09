package render

import (
	"github.com/leetcode-golang-classroom/golang-termdicator/internal/types"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type TermRender struct{}

func NewTermRender() *TermRender {
	return &TermRender{}
}

func (termRender *TermRender) RenderText(x, y int, msg string, color termbox.Attribute) {
	for _, ch := range msg {
		termbox.SetCell(x, y, ch, color, termbox.ColorDefault)
		w := runewidth.RuneWidth(ch)
		x += w
	}
}

func (termRender *TermRender) Start(ob types.OrderBook) error {

	for {
		err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		if err != nil {
			return err
		}
		ob.Render(0, 0)
		termbox.Flush()
		// TODO: this is blocking event, need to use non-block event for
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlC:
				return nil
			case termbox.KeyEsc:
				return nil
			}
		}
	}
}
