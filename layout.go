package home_automation

import (
	"github.com/metakeule/goh4/tag"
	//"github.com/metakeule/goh4/tag/short"
	"net/http"
)

var (
	content__ = tag.HTML("content").Placeholder()

	/*
		layout_ = tag.HTML5(
			tag.HEAD(
				short.CssHref("/css/style.css"),
				short.CssHref("/jquery-ui-1.10.3/css/custom-theme/jquery-ui-1.10.3.custom.css"),
			),
			tag.BODY(
				tag.UL(
					tag.ID("menu"),
					tag.LI(
						short.AHref("/", "alle Geräte"),
					),
					tag.LI(
						short.AHref("/details", "Details"),
					),
					tag.LI(
						short.AHref("/irgendwo", "Irgendwo"),
					),
				),
				content__,
			),
		).Compile("layout")
	*/

	/*
	   <!-- Static navbar -->
	       <div class="navbar navbar-default navbar-static-top">
	         <div class="container">
	           <div class="navbar-header">
	             <button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse">
	               <span class="icon-bar"></span>
	               <span class="icon-bar"></span>
	               <span class="icon-bar"></span>
	             </button>
	             <a class="navbar-brand" href="/">Yannis Home Automation Dingsbums</a>
	           </div>
	           <div class="navbar-collapse collapse">
	             <ul class="nav navbar-nav">
	               <li class="active"><a href="/edit">neues Gerät</a></li>
	             </ul>
	           </div><!--/.nav-collapse -->
	         </div>
	       </div>

	*/

	layout = tag.Doc(tag.HTML(`<!DOCTYPE html>
<html>
  <head>
    <title>Yannis Home Automation Dingsbums</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <!-- Bootstrap -->
    <link href="/bootstrap-3.0.0/dist/css/bootstrap.min.css" rel="stylesheet" media="screen">
     <link href="/bootstrap-3.0.0/dist/css/bootstrap-theme.min.css" rel="stylesheet" media="screen">
     <link href="/jquery-ui-1.10.3/css/custom-theme/jquery-ui-1.10.3.custom.css" rel="stylesheet" media="screen">
     <link href="/css/style.css" rel="stylesheet" media="screen">
  </head>
  <body>
 <div class="container">

      <!-- Main component for a primary marketing message or call to action -->
      <div class="jumbotron">
`),
		content__,
		tag.HTML(` 


      </div>
    </div> <!-- /container --> <script src="/js/jquery-1.10.2.js"></script>
    <script src="/bootstrap-3.0.0/dist/js/bootstrap.min.js"></script>
    <script src="/jquery-ui-1.10.3/js/jquery-ui-1.10.3.custom.js"></script>
    <script src="/js/homeautomation.js"></script>
    `),
	).Compile("huho")
)

func WriteLayout(content interface{}, rw http.ResponseWriter) {
	layout.MustReplace(content__.Set(content)).WriteTo(rw)
}
