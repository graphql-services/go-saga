package healthcheck

import (
	"fmt"
	"log"
	"net/http"
)

func StartHealthcheckServerOnDefaultPort() error {
	return StartHealthcheckServer("80")
}

// StartHealthcheckServer ...
func StartHealthcheckServer(port string) error {
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	fmt.Println("starting on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	return nil
}
