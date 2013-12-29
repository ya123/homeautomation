package home_automation

import (
	"fmt"
	"github.com/metakeule/goh4/tag"
	"github.com/metakeule/goh4/tag/short"
	"net/http"
	"strconv"

//	"os"
)

func dimmerInput(value int) string {
	//return fmt.Sprintf(`<input type="text" value="%d" disabled="disabled" class="dimmer-input form-control">`, value)
	return fmt.Sprintf(`<input type="text" value="%d" disabled="disabled" class="dimmer-input">`, value)
}

// <span class="badge">42</span>

/*
func dimmerNumber(value int) string {
	return fmt.Sprintf(`<span class="dimmer-number badge">%d<`, value)
}
*/

func Dim(rw http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	val := req.URL.Query().Get("value")
	d, ok := Devices.Get(name)
	if !ok {
		fmt.Printf("can't find device %s \n", name)
		rw.WriteHeader(500)
		return
	}
	if !d.Plugger.Dimmable() {
		fmt.Printf("device %s not dimmable \n", name)
		rw.WriteHeader(500)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	i, err := strconv.Atoi(val)
	if err != nil {
		rw.Write([]byte(`{"result": 0}`))
		return
	}

	ok = d.Plugger.Dim(i)

	if !ok {
		rw.Write([]byte(`{"result": 0}`))
		return
	}
	rw.Write([]byte(fmt.Sprintf(`{"result": %d}`, i)))
}

func List(rw http.ResponseWriter, req *http.Request) {
	tbody := tag.TBODY()

	for _, device := range Devices.Devices {
		//switchUrl := "/on"
		//status_ := short.AHref(switchUrl+"?name="+device.Name, "AN", tag.CLASS("btn"), tag.CLASS("btn-default"), tag.CLASS("disabled-btn"))
		//switch_ := short.AHref(switchUrl+"?name="+device.Name, "AUS", tag.CLASS("btn"), tag.CLASS("btn-default"))

		/*
			if device.On {
				switchUrl = "/off"
				status_ = short.AHref(switchUrl+"?name="+device.Name, "AN", tag.CLASS("btn"), tag.CLASS("btn-warning"))
				switch_ = short.AHref(switchUrl+"?name="+device.Name, "AUS", tag.CLASS("btn"), tag.CLASS("disabled-btn"), tag.CLASS("btn-default"))

			}
		*/

		var dimmer = ""
		if device.Plugger.Dimmable() {
			dimmer = `<div class="dimmer"></div>` + dimmerInput(0)
		}

		switchDiv := tag.DIV(tag.CLASS("switch"), tag.CLASS("btn"), device.Name)

		if device.On {
			switchDiv.AddClass("on")
			switchDiv.AddClass("btn-danger")
		} else {
			switchDiv.AddClass("btn-primary")
		}

		tbody.Add(
			tag.TR(
				tag.ATTR("data-name", device.Name),
				tag.TD(tag.CLASS("first"), switchDiv),
				tag.TD(tag.CLASS("middle"), tag.HTML(dimmer)),
				//	tag.TD(device.Plugger.Model()),
				//tag.TD(	status_, switch_),
				tag.TD(tag.CLASS("last"), short.AHref("/edit?name="+device.Name, tag.CLASS("btn-sm"), tag.CLASS("btn"),
					tag.SPAN(tag.CLASS("glyphicon"), tag.CLASS("glyphicon-cog")),
				)),
			),
		)
	}
	table := tag.TABLE(tag.CLASS("table"),
		tag.CLASS("table-responsive"),
		/*
			tag.THEAD(
				tag.TR(
					tag.TH("Name"),
					tag.TH("Model"),
					tag.TH(""),
					tag.TH(""),
				),
			),
		*/
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
	rw.Write([]byte("ok"))
	//http.Redirect(rw, req, "/", 302)
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
	//http.Redirect(rw, req, "/", 302)
	rw.Write([]byte("ok"))
}
