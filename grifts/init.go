package grifts

import (
	"github.com/L1ghtman2k/ScoreTrakWeb/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
