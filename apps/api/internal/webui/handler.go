package webui

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	files    fs.FS
	hasIndex bool
}

func NewHandler(files fs.FS, hasIndex bool) http.Handler {
	return &Handler{files: files, hasIndex: hasIndex}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.NotFound(w, r)
		return
	}
	if !h.hasIndex {
		writePlain(w, http.StatusServiceUnavailable, "Frontend belum di-build. Jalankan `make build` atau gunakan Vite dev server.")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" || path == "index.html" {
		h.serveIndex(w, r)
		return
	}
	if strings.Contains(path, ".") {
		http.FileServer(http.FS(h.files)).ServeHTTP(w, r)
		return
	}
	h.serveIndex(w, r)
}

func (h *Handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	data, err := fs.ReadFile(h.files, "index.html")
	if err != nil {
		writePlain(w, http.StatusInternalServerError, "index.html tidak ditemukan.")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(data))
}

func writePlain(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = io.WriteString(w, message)
}
