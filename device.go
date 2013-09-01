package home_automation

import (
	"encoding/json"
	"fmt"
	"github.com/metakeule/goh4/tag"
	"net/http"
	"os"
	"path"
	"sync"
)

type deviceMap struct {
	*sync.Mutex
	Devices map[string]*Device
}

//var repository = `D:/GO/GOPATH/src/github.com/ya123/homeautomation/devices.json`

var GOPATH = os.Getenv("GOPATH")

var repository = path.Join(GOPATH, "src", "github.com", "ya123", "homeautomation", "devices.json")

func Load() {
	file, err := os.Open(repository)
	if err != nil {
		panic(err.Error())
	}

	defer file.Close()
	info, _ := file.Stat()
	b := make([]byte, info.Size())
	_, err = file.Read(b)

	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(b, &Devices.Devices)
	if err != nil {
		panic(err.Error())
	}
}

func (this *deviceMap) Save() {
	this.Lock()
	defer this.Unlock()

	b, err := json.MarshalIndent(this.Devices, "", " ")
	if err != nil {
		panic(err.Error())
	}

	_ = os.Remove(repository)

	file, e := os.Create(repository)
	if e != nil {
		panic(e.Error())
	}
	_, er := file.Write(b)
	defer file.Close()

	if er != nil {
		panic(er.Error())
	}

}

func (this *deviceMap) Set(name string, d *Device) {
	this.Lock()
	defer this.Unlock()
	this.Devices[name] = d
}

func (this *deviceMap) Get(name string) (d *Device, ok bool) {
	this.Lock()
	defer this.Unlock()
	d, ok = this.Devices[name]
	return
}

func (this *deviceMap) Update(oldName string) (err error) {
	this.Lock()
	defer this.Unlock()
	d := this.Devices[oldName]
	if oldName == d.Name {
		return nil
	}
	if _, ok := this.Devices[d.Name]; ok {
		return fmt.Errorf("device with name %s exists already, choose another name", d.Name)
	}
	this.Devices[d.Name] = d
	delete(this.Devices, oldName)
	return nil
}

//var Devices = map[string]*Device{}
var Devices = &deviceMap{&sync.Mutex{}, map[string]*Device{}}

func NewDevice(name string, plugger Plugger) *Device {
	d := &Device{
		Name:        name,
		Plugger:     plugger,
		PluggerName: plugger.Name(),
	}
	Devices.Set(name, d)
	return d

}

type Device struct {
	On          bool
	Disabled    bool
	Name        string
	PluggerName string
	Plugger     Plugger
}

func (this *Device) UnmarshalJSON(data []byte) error {
	d := &struct {
		Disabled    bool
		Name        string
		PluggerName string
	}{}

	err := json.Unmarshal(data, d)

	if err != nil {
		return err
	}

	this.Disabled = d.Disabled
	this.Name = d.Name
	this.PluggerName = d.PluggerName
	this.On = false

	pl, ok := PluggerMap[this.PluggerName]
	if !ok {
		return fmt.Errorf("PLugger unknown: %s", this.PluggerName)
	}

	this.Plugger = pl()

	p := &struct {
		Plugger Plugger
	}{
		Plugger: this.Plugger,
	}
	err = json.Unmarshal(data, p)
	return err
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
