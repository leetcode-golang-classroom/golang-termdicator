package types

import "github.com/nsf/termbox-go"

type Renderer interface {
	RenderText(x int, y int, msg string, color termbox.Attribute)
}
