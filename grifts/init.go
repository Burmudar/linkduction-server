package grifts

import (
	"github.com/Burmudar/linkduction-server/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
