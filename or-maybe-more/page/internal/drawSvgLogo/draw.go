package drawSvgLogo

import (
	"errors"
	"fmt"
	"io"

	"github.com/ajstarks/svgo/svg"
)

// Logo is an logo created using Standard Vector Graphics (via svgo)
type Logo struct {
	// true if the svg has been compiled
	Encoded bool
	canvas  *svg.SVG
	LastErr error
}

func (l *Logo) Render(w io.WriteCloser) error {
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
