package healthcheck

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// StartHealthcheckServer ...
func StartHealthcheckServer() error {
	port := os.Getenv("HEALTHCHECK_PORT")
	if port == "" {
		port = "80"
	}

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	fmt.Println("starting on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	return nil
}
