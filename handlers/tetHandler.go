package handlers

import (
	"fmt"
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	// Get query param ?action=1
	action := r.URL.Query().Get("action")

	fmt.Println(action)
	if action == "1" {
		fmt.Fprintln(w, "🏠 Welcome to the Home Page!")
	} else {
		fmt.Fprintln(w, "🏠 Welcome to the Home Page 2!")
	}
}
