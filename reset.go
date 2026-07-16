package main

import(
	"net/http"
)


func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits.Store(0)
		w.WriteHeader(http.StatusOK)
}
