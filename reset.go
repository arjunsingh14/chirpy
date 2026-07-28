package main

import(
	"net/http"
)


func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
		if cfg.platform != "dev" {
			respondWithError(w, 403, "FORBIDDEN")
			return
		}
		cfg.db.ResetUsers(r.Context())
		w.WriteHeader(http.StatusOK)
}
