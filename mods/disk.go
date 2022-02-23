package mods

import (
	barista "barista.run"
	"barista.run/bar"
	"barista.run/format"
	"barista.run/modules/diskspace"
	"barista.run/outputs"
)

// DiskSpace Will return the amount of space left on the "/" root file partition.
func DiskSpace() {
	diskspacemod := diskspace.New("/").Output(func(i diskspace.Info) bar.Output {
		return outputs.Textf("%s/%s avail", format.IBytesize(i.Used()), format.IBytesize(i.Total))
	})

	barista.Add(diskspacemod)
}