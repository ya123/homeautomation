package home_automation

import (
	"github.com/metakeule/goh4/tag"
	"github.com/metakeule/goh4/tag/short"
	"net/http"
)

var (
	content__ = tag.HTML("content").Placeholder()

	layout = tag.HTML5(
		tag.HTML("\n<html>\n"),
		tag.HEAD(),
		tag.BODY(
			tag.UL(
				tag.LI(
					short.AHref("/", "alle Geräte"),
				),
				tag.LI(
					short.AHref("/details", "Details"),
				),
			),
			content__,
		),
		tag.HTML("\n</html>"),
	).Compile()
)

func WriteLayout(content interface{}, rw http.ResponseWriter) {
	layout.Replace(content__.Set(content)).WriteTo(rw)
}
