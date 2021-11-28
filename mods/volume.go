package mods

import (
	"barista.run/bar"
	"barista.run/colors"
	"barista.run/modules/volume"
	"barista.run/modules/volume/alsa"
	"barista.run/outputs"
	"barista.run/pango"
)

// Volume just takes care of reading volume data.
func Volume() (bar.Module, error) {
	volumeModule := volume.New(alsa.DefaultMixer()).Output(func(v volume.Volume) bar.Output {
		var iconName string
		if v.Mute {
			return outputs.Pango(pango.Icon("fa-volume-mute"), spacer, "MUT").Color(colors.Scheme("degraded"))
		}
		iconName = "off"
		pct := v.Pct()
		if pct > 66 {
			iconName = "up"
		} else if pct > 33 {
			iconName = "down"
		}

		return outputs.Pango(
			pango.Icon("fa-volume-"+iconName),
			spacer,
			pango.Textf("%2d%%", pct),
		)
	})

	return volumeModule, nil
}