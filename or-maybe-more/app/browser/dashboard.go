package ui

import (
	"github.com/kdevb0x/or-maybe-more/app"
	"github.com/maxence-charriere/app/pkg/app"
)

type user = app.Client

type dashboard struct {
	app.Compo
	panes []pane
}

type pane struct {
	x, y          int64
	height, width int64
	content       *app.Component
}

type Dashboard struct {
	User *user
	*dashboard
}
