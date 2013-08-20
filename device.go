package home_automation

import (
	"github.com/metakeule/goh4/tag"
	"net/http"
)

type Device struct {
	On       bool
	Disabled bool
	Name     string
	Plugger  Plugger
}

var Devices = map[string]*Device{}

func NewDevice(name string, plugger Plugger) *Device {
	d := &Device{
		Name:    name,
		Plugger: plugger,
	}
	Devices[name] = d
	return d

}

func (this *Device) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	status := "aus"
	if this.On {
		status = "an"
	}
	disabled := "aktiviert"
	if this.Disabled {
		disabled = "deaktiviert"
	}

	table := tag.TABLE(
		tag.TR(
			tag.TD(
				"Name",
			),
			tag.TD(
				this.Name,
			),
		),
		tag.TR(
			tag.TD(
				"Status",
			),
			tag.TD(
				status+"/"+disabled,
			),
		),
	)
	WriteLayout(table, rw)
}
