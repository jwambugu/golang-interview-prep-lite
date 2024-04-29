package user

import "net/http"

func Routes(svc *Service) *http.ServeMux {
	var (
		mux     = http.NewServeMux()
		handler = Handler{Svc: *svc}
	)

	mux.HandleFunc("/user", handler.AddUser)

	return mux
}
