package user

import "net/http"

func Routes(mux *http.ServeMux, svc Service) {
	var (
		handler = Handler{Svc: svc}
	)

	mux.HandleFunc("/login", handler.Authenticate)
	mux.HandleFunc("/user", handler.AddUser)
}
