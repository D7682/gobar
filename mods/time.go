package mods

import (
	barista "barista.run"
	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/modules/clock"
	"barista.run/outputs"
	"barista.run/pango"
	"time"
)

// Time will return a module with the time information based on your computer local time.
func Time() {
	localtime := clock.Local().Output(time.Second, func(now time.Time) bar.Output {
		return outputs.Pango(
			pango.Icon("material-today").Color(colors.Scheme("dim-icon")),
			now.Format(" Mon, Jan-02-2006 "),
			pango.Icon("material-access-time").Color(colors.Scheme("dim-icon")),
			now.Format(" 03:04:05PM"),
		).OnClick(click.RunLeft("gnome-calendar"))
	})

	barista.Add(localtime)
}
