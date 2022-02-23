package main

// Used the simple.go file

import (
	barista "barista.run"
	"github.com/D7682/gobar/mods"
	"github.com/D7682/gobar/setDefault"
	"log"
)

func main() {
	setDefault.Config()

	mods.Time()
	mods.GetWeather()
	mods.Volume()
	mods.GetMem()
	mods.DiskSpace()
	mods.CpuTemp()
	mods.UpdateCheck()

	if err := barista.Run(); err != nil {
		log.Fatal(err)
	}
}
