package funcs

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-chi/chi"
)

func HttpApi() {
	router := chi.NewRouter()
	WithRouter(router)
	if err := http.ListenAndServe(":8888", router); err != nil {
		panic(err)
	}
}

var staticRoutePath = "/static/"
var staticPath = "static"
var staticRoutePathLen = len(staticRoutePath)

func WithRouter(mux *chi.Mux) {

	dir, err := os.Stat(staticPath)
	if os.IsNotExist(err) || (dir != nil && !dir.IsDir()) {
		log.Fatal(err)
	}

	m := sync.Map{}

	mux.Route(staticRoutePath, func(root chi.Router) {
		root.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hi"))
		})

		assetsDir, err := filepath.Abs(staticPath)
		if err != nil {
			log.Fatal(err)
		}

		serverFile(root, "/static", "/", http.Dir(assetsDir), &m)
	})
}

func serverFile(r chi.Router, basePath string, path string, root http.FileSystem, m *sync.Map) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(basePath+path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path
		if len(reqPath) > staticRoutePathLen {
			key := reqPath[staticRoutePathLen:]
			etag, ok := m.Load(key)
			if ok {
				if v, ok := etag.(string); ok {
					w.Header().Set("Etag", fmt.Sprintf("\"%s\"", v))
				}
			}
		}

		fs.ServeHTTP(w, r)
	}))
}
