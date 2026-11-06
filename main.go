package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/milbertk/usersAppAK/helpers"
	"github.com/milbertk/usersAppAK/routers"
)

func main() {
	// .env is optional: real environment variables win (2026-10-24).
	_ = godotenv.Load()

	// Logs with levels; JSON in production (LOG_FORMAT, 2026-11-01). The old
	// log.Println lines go through the same logger.
	helpers.SetupLogger()

	// PORT (default 8080) and CORS_ALLOWED_ORIGINS come from the environment.
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	handler := routers.NewRouter()

	log.Println("✅ Server running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

/*package main

import (
	"log"
	"net/http"

	"github.com/milbertk/usersAppAK/routers"
	"github.com/rs/cors"
)

func main() {
	router := routers.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		// Debug: true, // opcional para ver logs del cors
	})

	handler := c.Handler(router)

	log.Println("✅ Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
*/
