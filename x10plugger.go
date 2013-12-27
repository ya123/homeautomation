package home_automation

import (
	"encoding/json"
	"fmt"
	"github.com/metakeule/goh4/tag"
	"github.com/metakeule/goh4/tag/short"
	"github.com/ya123/homeautomation/x10"
	"net/http"
	"strconv"
)

func init() {
	x10.InitCM11()
}

type X10Plugger struct {
	*x10.X10
}

func (this *X10Plugger) MarshalJSON() ([]byte, error) {
	s := &struct {
		HouseCode string
		UnitCode  int
		Model     string
	}{
		this.X10.HouseCode,
		this.X10.UnitCode,
		string(this.X10.Model),
	}
	return json.Marshal(s)
}

func (this *X10Plugger) UnmarshalJSON(data []byte) (err error) {
	s := &struct {
		HouseCode string
		UnitCode  int
		Model     string
	}{}
	err = json.Unmarshal(data, s)

	if err != nil {
		return
	}

	this.X10.HouseCode = s.HouseCode
	this.X10.UnitCode = s.UnitCode
	this.X10.Model = x10.Models[s.Model]
	this.X10.Name = "yannisdevice"
	return
}

func NewX10Plugger() Plugger {
	return &X10Plugger{&x10.X10{Name: "yannisdevice"}}
}
func (this *X10Plugger) Name() string {
	return "X10Plugger"
}

func (this *X10Plugger) On() {
	fmt.Printf("switching %s on\n", this.Name())
	this.Command(x10.On)
}

func (this *X10Plugger) Off() {
	fmt.Printf("switching %s off\n", this.Name())
	this.Command(x10.Off)
}

func (this *X10Plugger) Model() string {
	k := string(this.X10.Model)
	return models[k]
}

var models = map[string]string{
	"lm12": "Lamp dimmable",
	"lm15": "Lamp not dimmable",
	"am12": "Appliance",
	"tm12": "Tranceiver",
}

func (this *X10Plugger) Form(req *http.Request) string {

	modelselect := tag.SELECT(tag.ATTR("name", "model"), tag.CLASS("form-control"))

	for k, v := range models {
		if string(this.X10.Model) == k {
			modelselect.Add(
				tag.OPTION(tag.ATTR("value", k), v, tag.ATTR("selected", "selected")),
			)
		} else {
			modelselect.Add(
				tag.OPTION(tag.ATTR("value", k), v),
			)
		}

	}

	housecodeselect := tag.SELECT(
		tag.ATTR(
			"name", "housecode",
			"style", "width:70px;",
		),
		tag.CLASS("form-control"))

	housecodes := []string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
		"G",
		"H",
		"I",
		"J",
		"K",
		"L",
		"M",
		"N",
		"O",
		"P",
	}

	for _, v := range housecodes {
		if this.HouseCode == v {
			housecodeselect.Add(
				tag.OPTION(v, tag.ATTR("selected", "selected")),
			)
		} else {
			housecodeselect.Add(
				tag.OPTION(v))
		}

	}

	unitcodeselect := tag.SELECT(
		tag.ATTR(
			"name", "unitcode",
			"style", "width:70px;",
		),
		tag.CLASS("form-control"))

	for i := 1; i < 17; i++ {
		if this.UnitCode == i {
			unitcodeselect.Add(
				tag.OPTION(fmt.Sprintf("%v", i), tag.ATTR("selected", "selected")),
			)
		} else {
			unitcodeselect.Add(
				tag.OPTION(fmt.Sprintf("%v", i)),
			)
		}

	}

	return tag.Doc(
		tag.DIV(tag.CLASS("form-group"),
			short.LabelFor("model", "Model", tag.BR(), modelselect),
		),
		tag.DIV(tag.CLASS("form-group"),
			short.LabelFor("housecode", "Hauscode", tag.BR(), housecodeselect),
		),
		tag.DIV(tag.CLASS("form-group"),
			short.LabelFor("unitcode", "Unitcode", tag.BR(), unitcodeselect),
		),
	).String()
}

func (this *X10Plugger) Post(rw http.ResponseWriter, req *http.Request) {
	m := req.PostFormValue("model")
	this.X10.Model = x10.Models[m]
	h := req.PostFormValue("housecode")
	this.HouseCode = h
	u := req.PostFormValue("unitcode")
	i, _ := strconv.ParseInt(u, 10, 8)
	this.UnitCode = int(i)
}
