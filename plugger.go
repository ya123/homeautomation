package home_automation

import (
	"net/http"
)

var PluggerMap = map[string]func() Plugger{
	//	"LedPlugger":    NewLedPlugger,
	//	"SimplePlugger": NewSimplePlugger,
	"X10Plugger": NewX10Plugger,
}

type Plugger interface {
	Name() string
	On()
	Off()
	Model() string
	Dimmable() bool
	Form(req *http.Request) string
	Post(rw http.ResponseWriter, req *http.Request)
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

/*
<div style="background-color:red;opacity:0.5">

Farbe: <input name="color" value="red" />
Helligkeit:<input name="brightness" value="0.5" />

</div>

*/
