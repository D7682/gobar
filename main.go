package main

// Used the simple.go file

import (
	barista "barista.run"
	"barista.run/bar"
	"github.com/D7682/gobar/mods"
	"github.com/D7682/gobar/setDefault"
	"log"
)



func main() {
	setDefault.Config()
	modules := []func() (bar.Module, error){
		mods.Time,
		mods.GetWeather,
		mods.Volume,
		mods.GetMem,
		mods.DiskSpace,
		mods.CpuTemp,
		mods.UpdateCheck,
	}

	for _, module := range modules {
			barModule, err := module()
			if err != nil {
				log.Fatal(err)
			}
			barista.Add(barModule)
	}

	if err := barista.Run(); err != nil {
		log.Fatal(err)
	}
}
