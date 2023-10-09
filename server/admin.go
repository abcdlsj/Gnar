package server

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/abcdlsj/pipe/logger"
)

var (
	//go:embed assets
	assetsFs embed.FS
)

func (s *Server) startAdmin() {
	tmpl := template.Must(template.New("").ParseFS(assetsFs, "assets/*.html"))

	fe, _ := fs.Sub(assetsFs, "assets/static")
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.FS(fe))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{
			"forwards": s.forwards,
		}); err != nil {
			logger.ErrorF("execute index.html error: %v", err)
		}
	})

	logger.InfoF("admin server started on port %d", s.cfg.AdminPort)
	if err := http.ListenAndServe(":"+strconv.Itoa(s.cfg.AdminPort), nil); err != nil {
		logger.FatalF("admin server error: %v", err)
	}
}
