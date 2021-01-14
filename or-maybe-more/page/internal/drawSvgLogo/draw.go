// package drawSvgLogo creates and renders an svg logo
package drawSvgLogo

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	svg "github.com/ajstarks/svgo"
)

// Logo is an logo created using Standard Vector Graphics (via svgo)
type Logo struct {
	// true if the svg has been compiled
	Encoded bool
	Canvas  *svg.SVG
	LastErr error
}

func NewEmptyCanvas(w io.Writer) *Logo {
	var l = new(Logo)
	l.Canvas = svg.New(w)
	l.Encoded = false // not sure if needed b/c bool zero value is false
	return l

}

// Error implements the error interface; returns l.LastErr if present.
func (l *Logo) Error() string {
	return l.LastErr.Error()
}

func (l *Logo) Render() error {
	if !l.Encoded {
		return errors.New("Logo must be encoded before rendering")
	}
	// *BUG: This cant be right; too tired
	// *l.canvas.Writer = w

	// *TODO: find out what it rakes to build the svg outpout from an
	// *svg.SVG
	if errors.Unwrap((*l).LastErr) != nil {
		return fmt.Errorf("unable to render because of an unresolved error: %w\n", (*l).LastErr)
	}

	// for now
	return errors.New("rendering hasn't implemented")
}

// Renders a preview to stdout
func (l *Logo) Preview() {

}

func MakeOMMLogo(w int, h int) []byte {
	var b = new(bytes.Buffer)
	l := NewEmptyCanvas(b)
	c := l.Canvas
	var (
		x0         = 0
		y0         = 0
		topMargin  = 5
		leftMargin = 5
	)

	c.Start(w, h)
	c.Rect(x0, y0, w, h, "fill:skyblue;")
	c.Text(x0+leftMargin, y0+topMargin, "Or Maybe More", "font-size:30px;fill:black;text-anchor:middle;")
	// cloud
	c.Def()
	c.Ellipse(30, 70, 30, 15)

	// TODO finish this

}
