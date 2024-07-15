package service

import (
	"github.com/pflow-dev/pflow-eth/config"
	"html/template"
)

func (s *Server) IndexPage() *template.Template {
	return s.indexPage
}

func (s *Server) IndexTemplateSource() string {
	out := `<!doctype html>
	<html lang="en">
	<head>
		<title>pflow | StateMachine </title>
        <meta charset="utf-8"/>
        <meta name="viewport" content="width=device-width,initial-scale=1"/>
        <meta name="msapplication-TileColor" content="#da532c">
        <meta name="theme-color" content="#ffffff">
	<link href="/static/css/main.`

	out += config.CssBuild + `.css" rel="stylesheet">`
	out += SessionDataScript
	out += `<script defer="defer" src=/static/js/main.` + config.JsBuild + `.js></script>`
	out += `</head>
		<body><noscript>You need to enable JavaScript to run this app.</noscript>
        <div id="root"></div>
    </body>
    </html>`

	return out
}

const (
	SessionDataScript = `<script>
	const NOTE = "inject js code here";
</script>`
)
