package home_automation

import (
	//"fmt"
	"github.com/metakeule/goh4/tag"
	"github.com/metakeule/goh4/tag/short"
	"net/http"

//	"os"
)

func List(rw http.ResponseWriter, req *http.Request) {
	tbody := tag.TBODY()

	for _, device := range Devices.Devices {
		switchUrl := "/on"
		status_ := short.AHref(switchUrl+"?name="+device.Name, "AN", tag.CLASS("btn"), tag.CLASS("btn-default"), tag.CLASS("disabled-btn"))
		switch_ := short.AHref(switchUrl+"?name="+device.Name, "AUS", tag.CLASS("btn"), tag.CLASS("btn-default"))

		if device.On {
			switchUrl = "/off"
			status_ = short.AHref(switchUrl+"?name="+device.Name, "AN", tag.CLASS("btn"), tag.CLASS("btn-warning"))
			switch_ = short.AHref(switchUrl+"?name="+device.Name, "AUS", tag.CLASS("btn"), tag.CLASS("disabled-btn"), tag.CLASS("btn-default"))

		}
		tbody.Add(
			tag.TR(
				tag.TD(device.Name),
				tag.TD(device.Plugger.Model()),
				tag.TD(
					status_, switch_),
				tag.TD(short.AHref("/edit?name="+device.Name, "bearbeiten", tag.CLASS("btn-primary"), tag.CLASS("btn"))),
			),
		)
	}
	table := tag.TABLE(tag.CLASS("table"),
		tag.CLASS("table-responsive"),
		tag.THEAD(
			tag.TR(
				tag.TH("Name"),
				tag.TH("Model"),
				tag.TH(""),
				tag.TH(""),
			),
		),
		tbody,
	)
	WriteLayout(table, rw)
}

func Details(rw http.ResponseWriter, req *http.Request) {
	urlValues := req.URL.Query()
	name := urlValues.Get("name")
	d, ok := Devices.Get(name)
	if !ok {
		WriteLayout("nicht gefunden", rw)
		return
	}
	d.ServeHTTP(rw, req)
}

func Edit(rw http.ResponseWriter, req *http.Request) {
	urlValues := req.URL.Query()
	name := urlValues.Get("name")
	var d *Device
	var ok bool
	if name == "" {
		pl := NewX10Plugger()
		d = &Device{
			Plugger:     pl,
			PluggerName: pl.Name(),
		}
	} else {
		d, ok = Devices.Get(name)
		if !ok {
			WriteLayout("nicht gefunden", rw)
			return
		}
	}
	pluggerForm := d.Plugger.Form(req)

	selectbox := tag.SELECT(tag.ATTR("name", "plugger"), tag.CLASS("form-control"))

	for k, _ := range PluggerMap {
		if d.Plugger.Name() == k {
			selectbox.Add(tag.OPTION(k, tag.ATTR("selected", "selected")))
		} else {
			selectbox.Add(tag.OPTION(k))
		}
	}

	target := "/change"
	if name != "" {
		target += "?name=" + name
	}

	form := short.FormPost(target,
		tag.CLASS("form-inline"),
		tag.ATTR("role", "form"),
		tag.DIV(tag.CLASS("form-group"),
			short.LabelFor("name", "Name", tag.BR(), tag.INPUT(
				tag.CLASS("form-control"),
				tag.ATTR("name", "name"),
				tag.ATTR("value", name),
			),
			),
		),
		tag.DIV(tag.CLASS("form-group"),
			short.LabelFor("plugger", "Plugger", tag.BR(), selectbox),
		),

		tag.HTML(pluggerForm),
		tag.DIV(tag.CLASS("form-group"),
			tag.LABEL(
				tag.BR(),
				tag.BUTTON(tag.CLASS("btn"),
					tag.CLASS("btn-primary"),
					tag.CLASS("form-control"),
					tag.ATTR("type", "submit", "style", "width:120px;"), "speichern"),
			),
		),
	)

	//layout.Replace(content__.Set(form)).WriteTo(rw)
	WriteLayout(form, rw)
}

func Change(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	name := req.URL.Query().Get("name")

	if name == "" {
		plugger := req.PostFormValue("plugger")
		d := NewDevice(req.PostFormValue("name"), PluggerMap[plugger]())
		d.Plugger.Post(rw, req)
	} else {
		d, ok := Devices.Get(name)
		if !ok {
			rw.WriteHeader(404)
			return
		}
		changeName := req.PostFormValue("name")
		d.Name = changeName
		if err := Devices.Update(name); err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		plugger := req.PostFormValue("plugger")

		//fmt.Println(plugger)
		d.Plugger = PluggerMap[plugger]()

		d.Plugger.Post(rw, req)
	}

	Devices.Save()
	//json.MarshalIndent(Device, prefix, indent)

	/*
		file, er := os.Create(repository)
		if er != nil {
			panic(er.Error())
		}

		file.WriteString("huho")
		file.Close()
	*/

	http.Redirect(rw, req, "/", 302)
}

func On(rw http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	d, ok := Devices.Get(name)
	if !ok {
		rw.WriteHeader(500)
		return
	}
	d.Plugger.On()
	d.On = true
	http.Redirect(rw, req, "/", 302)
}

func Off(rw http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	d, ok := Devices.Get(name)
	if !ok {
		rw.WriteHeader(500)
		return
	}
	d.Plugger.Off()
	d.On = false
	http.Redirect(rw, req, "/", 302)
}
