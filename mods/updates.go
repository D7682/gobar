package mods

import (
	barista "barista.run"
	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/colors"
	"barista.run/outputs"
	"github.com/gen2brain/beeep"
	"github.com/martinohmann/barista-contrib/modules/updates"
	"github.com/martinohmann/barista-contrib/modules/updates/pacman"
	"log"
	"time"
)

func UpdateCheck() {
	updatesModule := updates.New(pacman.Provider).Output(func(info updates.Info) bar.Output {
		text := outputs.Textf("%d updates", info.Updates).OnClick(click.Left(func() {
			if err := beeep.Notify("Available Pacman Updates", info.PackageDetails.String(), ""); err != nil {
				log.Fatal(err)
			}
		}))

		switch count := info.Updates; {
		case count == 0:
			return nil
		case count > 125:
			return text.Color(colors.Scheme("bad"))
		default:
			return text.Color(colors.Scheme("good"))
		}
	})

	updatesModule = updatesModule.Every(time.Second * 30)

	barista.Add(updatesModule)
}
