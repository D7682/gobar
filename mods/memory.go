package mods

import (
	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/format"
	"barista.run/modules/meminfo"
	"barista.run/outputs"
	"barista.run/pango"
)

// GetMem gets the amount of free ram left on the hard drive.
func GetMem() (bar.Module, error) {
	startTaskManager := click.RunLeft("urxvt", "-e", "htop")
	freeMemModule := meminfo.New().Output(func(m meminfo.Info) bar.Output {
		out := outputs.Pango(pango.Icon("material-memory"), format.IBytesize(m.Available()))
		freeGigs := m.Available().Gigabytes()
		switch {
		case freeGigs < 0.5:
			out.Urgent(true)
		case freeGigs < 1:
			out.Color(colors.Scheme("bad"))
		case freeGigs < 2:
			out.Color(colors.Scheme("degraded"))
		case freeGigs > 12:
			out.Color(colors.Scheme("good"))
		}
		out.OnClick(startTaskManager)
		return out
	})
	return freeMemModule, nil
}
