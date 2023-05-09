package presentation

import "net/http"

func (p *Presenter) Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
