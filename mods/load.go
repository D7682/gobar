package mods

import (
	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/modules/sysinfo"
	"barista.run/outputs"
	"time"
)

// GetLoad is used to get the average cpu load.
func GetLoad() (bar.Module, error) {
	startTaskManager := click.RunLeft("urxvt", "-e", "htop")
	loadAvg := sysinfo.New().Output(func(s sysinfo.Info) bar.Output {
		out := outputs.Textf("%0.2f %0.2f", s.Loads[0], s.Loads[2])
		// Load averages are unusually high for a few minutes after boot.
		if s.Uptime < 10*time.Minute {
			// so don't add colours until 10 minutes after system start.
			return out
		}
		switch {
		case s.Loads[0] > 128, s.Loads[2] > 64:
			out.Urgent(true)
		case s.Loads[0] > 64, s.Loads[2] > 32:
			out.Color(colors.Scheme("bad"))
		case s.Loads[0] > 32, s.Loads[2] > 16:
			out.Color(colors.Scheme("degraded"))
		}
		out.OnClick(startTaskManager)
		return out
	})
	return loadAvg, nil
}