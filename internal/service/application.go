package service

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/pflow-dev/pflow-eth/internal/config"
	"net/http"
)

func (s *Service) appSource() string {
	out := `<!doctype html>
	<html lang="en">
	<head>
		<title>pflow | StateMachine </title>
        <meta charset="utf-8"/>
        <meta name="viewport" content="width=device-width,initial-scale=1"/>
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

// add any front-end routes needed by react router here
var applicationRoutes = []string{
	"/",
}

// all return the index page, assuming the front-end router will handle the route
func (s *Service) applicationRoutes(box *rice.Box) {
	for _, route := range applicationRoutes {
		s.WrapHandler(route, func(vars map[string]string, w http.ResponseWriter, r *http.Request) {
			if err := s.applicationPage.Execute(w, nil); err != nil {
				s.Logger.Printf("Error rendering index page: %v\n", err)
			}
		})
	}

	s.Router.HandleFunc("/{file}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		http.StripPrefix("/", http.FileServer(box.HTTPBox())).ServeHTTP(w, r)
	})

	s.Router.HandleFunc("/static/js/{jsBuild}",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			f, boxErr := box.Open("static/js/" + vars["jsBuild"])
			if boxErr != nil {
				http.Error(w, boxErr.Error(), http.StatusNotFound)
				return
			}

			fileInfo, fileErr := f.Stat()
			if fileErr != nil {
				http.Error(w, fileErr.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set("Cache-Control", "public, max-age=31536000")
			http.ServeContent(w, r, "main."+vars["jsBuild"]+".js", fileInfo.ModTime(), f)
		})

	s.Router.HandleFunc("/static/css/{cssBuild}",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			f, boxErr := box.Open("static/css/" + vars["cssBuild"])
			if boxErr != nil {
				http.Error(w, boxErr.Error(), http.StatusNotFound)
				return
			}
			fileInfo, fileErr := f.Stat()
			if fileErr != nil {
				http.Error(w, fileErr.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set("Cache-Control", "public, max-age=31536000")
			http.ServeContent(w, r, "main."+vars["cssBuild"]+".css", fileInfo.ModTime(), f)
		})

}
