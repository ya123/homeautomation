package home_automation

import (
	"github.com/metakeule/goh4/tag"
	"github.com/metakeule/goh4/tag/short"
	"net/http"
)

func List(rw http.ResponseWriter, req *http.Request) {
	table := tag.TABLE()

	for _, device := range Devices {
		table.Add(
			tag.TR(
				tag.TD(device.Name),
				tag.TD(short.AHref("/details?name="+device.Name, "Details")),
				tag.TD(short.AHref("/edit?name="+device.Name, "bearbeiten")),
			),
		)
	}
	WriteLayout(table, rw)
}

func Details(rw http.ResponseWriter, req *http.Request) {
	urlValues := req.URL.Query()
	name := urlValues.Get("name")
	d, ok := Devices[name]
	if !ok {
		WriteLayout("nicht gefunden", rw)
		return
	}
	d.ServeHTTP(rw, req)
}

func Edit(rw http.ResponseWriter, req *http.Request) {
	urlValues := req.URL.Query()
	name := urlValues.Get("name")
	d, ok := Devices[name]
	if !ok {
		WriteLayout("nicht gefunden", rw)
		return
	}
	pluggerForm := d.Plugger.Form(req)
	form := short.FormPost("/change?name="+name,
		tag.SELECT(
			tag.OPTION("Plugger1")),
		tag.HTML(pluggerForm),
		short.InputSubmit("ändern"))

	//layout.Replace(content__.Set(form)).WriteTo(rw)
	WriteLayout(form, rw)
}

func Change(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	name := req.URL.Query().Get("name")
	d, ok := Devices[name]
	if !ok {
		rw.WriteHeader(404)
		return
	}
	d.Plugger.Post(rw, req)
	http.Redirect(rw, req, "/", 302)
}
