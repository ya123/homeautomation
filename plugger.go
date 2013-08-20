package home_automation

import (
	"fmt"
	"github.com/metakeule/goh4/tag"
	"net/http"
	"strconv"
)

var PluggerMap = map[string]func() Plugger{
	"LedPlugger": NewLedPlugger,
}

type Plugger interface {
	Name() string
	Form(req *http.Request) string
	Post(rw http.ResponseWriter, req *http.Request)
}

type LedPlugger struct {
	Color      string
	Brightness float32
}

func NewLedPlugger() *LedPlugger {
	return &LedPlugger{}
}

func (this *LedPlugger) Name() string {
	return "LedPlugger"
}

func (this *LedPlugger) Form(req *http.Request) string {
	return tag.DIV(
		tag.ATTR("style", "background-color:"+this.Color+";opacity:"+fmt.Sprintf("%v", this.Brightness)),
		"Farbe: ", tag.INPUT(
			tag.ATTR("name", "color"),
			tag.ATTR("value", this.Color),
		),
		"Helligkeit:", tag.INPUT(
			tag.ATTR("name", "brightness"),
			tag.ATTR("value", fmt.Sprintf("%v", this.Brightness)),
		),
	).String()
}

/*
<div style="background-color:red;opacity:0.5">

Farbe: <input name="color" value="red" />
Helligkeit:<input name="brightness" value="0.5" />

</div>

*/

func (this *LedPlugger) Post(rw http.ResponseWriter, req *http.Request) {
	color := req.PostFormValue("color")
	brightness := req.PostFormValue("brightness")

	this.Color = color
	b, _ := strconv.ParseFloat(brightness, 32)
	this.Brightness = float32(b)
}
