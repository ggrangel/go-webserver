package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerHealthz)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerDisplayMetrics)

	mux.HandleFunc("/api/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/chirps", handlerCreateChirp(filepathRoot+"/database.json"))

	mux.HandleFunc("GET /api/chirps/", handlerGetChirps)

	mux.HandleFunc("POST /api/users", handlerCreateUser)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server started from %s on port %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
