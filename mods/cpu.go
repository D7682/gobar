package mods

import (
	barista "barista.run"
	"barista.run/bar"
	"barista.run/colors"
	"barista.run/modules/cputemp"
	"barista.run/outputs"
	"barista.run/pango"
	"github.com/martinlindhe/unit"
	"time"
)

func CpuTemp() {
	cpuModule := cputemp.New().RefreshInterval(2 * time.Second).Output(func(temp unit.Temperature) bar.Output {
		out := outputs.Pango(
			pango.Icon("mdi-fan"), spacer,
			pango.Textf("%2d℃", int(temp.Celsius())),
		)
		switch {
		case temp.Celsius() > 90:
			out.Urgent(true)
		case temp.Celsius() > 70:
			out.Color(colors.Scheme("bad"))
		case temp.Celsius() > 60:
			out.Color(colors.Scheme("degraded"))
		}
		return out
	})

	barista.Add(cpuModule)
}