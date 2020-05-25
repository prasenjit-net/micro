package homepage

import (
	"log"
	"micro/util"
	"net/http"
	"path/filepath"
)

type HomePage struct {
}

func New() *HomePage {
	return &HomePage{}
}

func (h *HomePage) Home(w http.ResponseWriter, r *http.Request) {
	log.Print("Homepage request processed")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello World!\n"))
	util.LogAndIgnoreError(err)
}

func (h *HomePage) Register(prefix string, mux *http.ServeMux) {
	homePath := filepath.Join(prefix, "/")
	mux.HandleFunc(homePath, util.LogReqTime(h.Home))
}
