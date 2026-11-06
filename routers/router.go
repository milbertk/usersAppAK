package routers

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/milbertk/security"
	haiagent "github.com/milbertk/usersAppAK/handlers/HaiAgent"
	hcompany "github.com/milbertk/usersAppAK/handlers/Hcompany"
	hsupportrequest "github.com/milbertk/usersAppAK/handlers/HsupportRequest"
	husercreate "github.com/milbertk/usersAppAK/handlers/HuserCreate"
	huserlogin "github.com/milbertk/usersAppAK/handlers/HuserLogin"
	"github.com/milbertk/usersAppAK/handlers/common"
	securitygen "github.com/milbertk/usersAppAK/handlers/securityGen"
	Husers "github.com/milbertk/usersAppAK/handlers/users"
)

// ----------------------------------------------------------------------
// Configuration
// ----------------------------------------------------------------------

// allowedOrigins is the exact-match list for production.
// Add LAN IPs here for local cross-device testing, or set ALLOW_ANY_ORIGIN=1
// to bypass the list entirely (dev only — see corsHandler).
var allowedOrigins = map[string]bool{
	"http://localhost:3000":       true,
	"http://localhost:3001":       true,
	"http://localhost:3002":       true,
	"http://127.0.0.1:3000":       true,
	"http://192.168.100.110:3000": true,
	"http://192.168.100.110:3001": true,
	"http://192.168.100.110:3002": true,
	"http://192.168.100.110":      true,
}

// allowedHeaders is the set of request headers the frontend may send.
// We expose it as a constant so the preflight echoes exactly what we accept.
const (
	allowedMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	allowedHeaders = "Authorization, Content-Type, X-Requested-With, X-Captcha-Token"
	maxAgeSeconds  = "3600" // cache preflight for 1 hour
)

// ----------------------------------------------------------------------
// CORS — applied at the OUTERMOST layer, before mux ever sees the request.
// This is the key reliability win: OPTIONS preflights are answered
// immediately and never enter the router, so there's no way for mux's
// method matching to drop them.
// ----------------------------------------------------------------------

// init loads the CORS rules for the current environment (APP_ENV, see
// cors_config.go, 2026-10-25 security S9): production = only
// CORS_ALLOWED_ORIGINS + the origins file; test = any origin; not set = the
// local list + CORS_ALLOWED_ORIGINS (as before).
func init() {
	allowedOrigins = loadCORSOrigins(allowedOrigins)
}

// healthHandler answers GET / so a load balancer can check the service.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok","service":"usersAppAK"}`))
}

// isOriginAllowed returns true if `origin` is on the allowlist, or when the
// environment accepts any origin (APP_ENV=test, or ALLOW_ANY_ORIGIN=1 when
// APP_ENV is not "production").
func isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	if corsAllowAnyOrigin() {
		return true
	}
	return allowedOrigins[origin] || originAllowedByDomain(origin)
}

// writeCORS writes the standard CORS response headers when the request's
// Origin is allowed. Safe to call for any request — it no-ops on
// non-allowed origins.
func writeCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if !isOriginAllowed(origin) {
		return
	}
	h := w.Header()
	h.Set("Access-Control-Allow-Origin", origin)
	h.Set("Vary", "Origin")
	h.Set("Access-Control-Allow-Credentials", "true")
	h.Set("Access-Control-Allow-Methods", allowedMethods)
	h.Set("Access-Control-Allow-Headers", allowedHeaders)
	h.Set("Access-Control-Max-Age", maxAgeSeconds)
}

// corsHandler wraps the entire router. It handles OPTIONS preflights
// itself (always returning 204 immediately) and writes CORS headers
// on every other response so downstream handlers don't need to.
func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeCORS(w, r)

		// Short-circuit ALL preflights here. mux never sees them.
		if r.Method == http.MethodOptions {
			// If the origin wasn't allowed, the headers above weren't
			// written — the browser will reject anyway. Return 204 either
			// way; the meaningful signal is the presence/absence of the
			// Allow-Origin header.
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logging + panic recovery: see request_log.go (requestLogMiddleware,
// 2026-11-01). It replaced loggingMiddleware and recoveryMiddleware.

// ----------------------------------------------------------------------
// Router
// ----------------------------------------------------------------------

// NewRouter returns the fully-wired http.Handler.
//
// Layering (outermost first):
//
//	corsHandler                  →  preflights answered, CORS headers set
//	  mux.Router                 →  route matching
//	    requestLogMiddleware     →  one log line per call, request id, panic -> 500
//	      your handlers
//
// Note we return http.Handler (not *mux.Router) because corsHandler is
// the outermost wrapper. Update main.go to:
//
//	http.ListenAndServe(":8080", routers.NewRouter())
func NewRouter() http.Handler {
	router := mux.NewRouter()

	// Mux-level middleware: applied AFTER CORS has already handled preflights.
	router.Use(requestLogMiddleware)

	// Routes
	router.HandleFunc("/", healthHandler).Methods("GET")
	router.HandleFunc("/createuser", husercreate.HuserCreate).Methods("GET")
	router.Handle(
		"/companiescreate",
		security.RequireActiveUser(http.HandlerFunc(hcompany.HcompanyCreate)),
	).Methods("POST")
	router.Handle(
		"/companiesget",
		security.RequireActiveUser(http.HandlerFunc(hcompany.HcompanyGetAll)),
	).Methods("GET")
	router.HandleFunc("/login", huserlogin.HuserLogin).Methods("POST")
	// Only the OWNER / MANAGER of the company (companyId query param or
	// X-Company-ID header) can search, by exact email (2026-10-25, security S10).
	router.Handle(
		"/searchUser",
		security.RequireActiveCompanyUser(http.HandlerFunc(Husers.HsearchUsersByEmailNameOrId)),
	).Methods("GET")
	router.Handle(
		"/userRequestCreate",
		security.RequireActiveUser(http.HandlerFunc(hcompany.HCreateCompanyUserRequest)),
	).Methods("POST")
	router.Handle(
		"/searchRequest",
		security.RequireActiveUser(http.HandlerFunc(hcompany.HSearchCompanyUserRequests)),
	).Methods("GET")
	router.Handle(
		"/cancelrequest",
		security.RequireActiveUser(http.HandlerFunc(hcompany.HUpdateCompanyUserRequestStatus)),
	).Methods("POST", "PUT")
	router.Handle(
		"/notifications",
		security.RequireActiveUser(http.HandlerFunc(common.HGetNotifications)),
	).Methods("GET")
	router.Handle(
		"/invitations",
		security.RequireActiveUser(http.HandlerFunc(hcompany.HGetPendingInvitationsByUser)),
	).Methods("GET")
	router.Handle(
		"/processrequestbyinvited",
		security.RequireActiveUser(http.HandlerFunc(hcompany.HProcessCompanyUserInvitation)),
	).Methods("POST", "PUT")
	router.Handle(
		"/searchuserbycompany",
		security.RequireActiveUser(http.HandlerFunc(hcompany.HSearchCompanyUsers)),
	).Methods("GET")
	router.Handle(
		"/updateuserofcompany",
		security.RequireActiveUser(http.HandlerFunc(hcompany.HUpdateCompanyUser)),
	).Methods("POST", "PUT")
	router.Handle(
		"/getuserinfo",
		security.AuthMiddleware(http.HandlerFunc(securitygen.HuserMe)),
	).Methods("GET")
	router.Handle(
		"/logout",
		http.HandlerFunc(securitygen.HuserLogout),
	).Methods("POST")
	router.Handle(
		"/sendinvitation",
		security.RequireActiveUser(http.HandlerFunc(huserlogin.HSendInvitation)),
	).Methods("POST")
	router.Handle(
		"/changepassword",
		security.RequireActiveUser(http.HandlerFunc(securitygen.HChangePassword)),
	).Methods("POST")
	router.Handle(
		"/transfercompanyowner",
		security.RequireActiveUser(http.HandlerFunc(hcompany.HTransferCompanyOwner)),
	).Methods("POST", "PUT")
	router.Handle(
		"/companiesmy",
		security.RequireActiveUser(http.HandlerFunc(hcompany.HcompanyGetByUser)),
	).Methods("GET")

	// Support requests (2026-10-24) — logged users only. Saved in
	// tbl_support_request and emailed to SUPPORT_EMAIL_TO.
	router.Handle(
		"/supportrequests",
		security.RequireActiveUser(http.HandlerFunc(hsupportrequest.HCreateSupportRequest)),
	).Methods("POST")

	// AI Agent (Anthropic / Claude Haiku) — logged users only
	router.Handle(
		"/aichat",
		security.RequireActiveUser(http.HandlerFunc(haiagent.HaiSendMessage)),
	).Methods("POST")
	router.Handle(
		"/aiconversations",
		security.RequireActiveUser(http.HandlerFunc(haiagent.HaiGetConversations)),
	).Methods("GET")
	router.Handle(
		"/aiconversationmessages",
		security.RequireActiveUser(http.HandlerFunc(haiagent.HaiGetMessages)),
	).Methods("GET")

	// Fallbacks. Because CORS already short-circuits OPTIONS at the outer
	// layer, these will only fire for genuine mismatches (e.g. typo'd path,
	// wrong HTTP method on a real route). Both write CORS headers so the
	// browser can actually display the error in DevTools.
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeCORS(w, r)
		http.Error(w, `{"status":"NotFound","message":"unknown route"}`, http.StatusNotFound)
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeCORS(w, r)
		http.Error(w, `{"status":"MethodNotAllowed","message":"wrong method for this route"}`,
			http.StatusMethodNotAllowed)
	})

	router.Handle(
		"/canchangefirebasepassword",
		security.RequireActiveUser(http.HandlerFunc(securitygen.HCanChangeFirebasePassword)),
	).Methods("GET")

	router.Handle(
		"/changeemail",
		security.RequireActiveUser(http.HandlerFunc(securitygen.HChangeEmail)),
	).Methods("POST", "PUT")

	// Startup log
	log.Println("========================================")
	if os.Getenv("ALLOW_ANY_ORIGIN") == "1" {
		log.Println("⚠️  ALLOW_ANY_ORIGIN=1 — any Origin will be accepted (DEV ONLY)")
	} else {
		log.Printf("Allowed origins: %s", strings.Join(keys(allowedOrigins), ", "))
	}
	log.Println("Rutas registradas:")
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, _ := route.GetMethods()
			log.Printf("  %-8v %s", methods, pathTemplate)
		}
		return nil
	})
	log.Println("========================================")

	// Wrap the whole router with the CORS handler.
	return corsHandler(router)
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
