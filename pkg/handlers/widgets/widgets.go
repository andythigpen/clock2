package widgets

import (
	"context"
	"io"
)

type Widget interface {
	ShouldDisplay() bool
	Render(context.Context, io.Writer)
}
