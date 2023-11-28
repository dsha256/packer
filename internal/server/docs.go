package server

import "net/http"

const (
	docsTheneoURL = "https://app.theneo.io/rad/gymshark/get-packets"
)

func (s *Server) docsHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, docsTheneoURL, http.StatusFound)
}
