package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wolv89/troster/internal/models"
)

func Start(data *models.ScrapedData, port int) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/fixtures", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	mux.Handle("GET /", http.FileServer(http.Dir("web/public")))

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server running at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, mux)
}
