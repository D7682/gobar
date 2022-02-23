package mods

import (
	barista "barista.run"
	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/modules/meminfo"
	"barista.run/outputs"
	"barista.run/pango"
	"fmt"
)

// GetMem gets the amount of free ram left on the hard drive.
// Try to get the amount of used ram correctly onto here before recompiling.
func GetMem() {
	startTaskManager := click.RunLeft("urxvt", "-e", "htop")
	usedMemModule := meminfo.New().Output(func(m meminfo.Info) bar.Output {
		memoryLeft := m.Available().Gigabytes()

		out := outputs.Pango(pango.Icon("material-memory"), fmt.Sprintf("%0.2f", memoryLeft))
		switch {
		case memoryLeft < 2.0:
			out.Color(colors.Scheme("bad"))
		default:
			out.Color(colors.Scheme("good"))
		}
		out.OnClick(startTaskManager)
		return out
	})

	barista.Add(usedMemModule)
}
