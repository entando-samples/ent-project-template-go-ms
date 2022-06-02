package main

import (
	kca "github.com/GermanoGiudici/keycloak-go-adapter"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const healthCheckResponse = "{\"status\":\"UP\"}"
const apiResponse = "{\"payload\":\"test Data\"}"

func Protect(f http.HandlerFunc, roles []string, any bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorized, httpStatus, err := kca.Protect(r, roles, any)
		log.Printf("URL %v authorized: %t err: %v\n", r.URL.Path, authorized, err)
		if !authorized {
			w.WriteHeader(httpStatus)
		} else {
			f(w, r)
		}
	}
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	var keycloakClientId = os.Getenv("KEYCLOAK_CLIENT_ID")                  //"internal"
	var keycloakRealm = os.Getenv("KEYCLOAK_REALM")                         //"entando"
	var serverServletContextPath = os.Getenv("SERVER_SERVLET_CONTEXT_PATH") //""
	var keycloakAuthUrl = os.Getenv("KEYCLOAK_AUTH_URL")                    //"http://localhost:9080/auth"
	serverServletContextPath = serverServletContextPath + "/"

	type OpenApiConfig struct {
		ServersURL       string
		AuthorizationURL string
		KeycloakRealm    string
	}
	serverUrl := strings.Replace(keycloakAuthUrl, "/auth", "", 1)
	kca.Init(keycloakClientId, serverUrl, keycloakRealm)

	tplOpenApi := template.Must(template.ParseFiles("open-api.tpl"))
	tplSwaggerUIInit := template.Must(template.ParseFiles("swagger-ui/swagger-initializer.tpl"))

	tplConfig := OpenApiConfig{
		ServersURL:       serverServletContextPath,
		AuthorizationURL: keycloakAuthUrl,
		KeycloakRealm:    keycloakRealm,
	}

	r := mux.NewRouter()
	subrouter := r.PathPrefix(serverServletContextPath).Subrouter()

	sh := http.StripPrefix(serverServletContextPath+"swagger-ui", http.FileServer(http.Dir("./swagger-ui")))
	r.PathPrefix(serverServletContextPath + "swagger-ui/").Handler(sh).Name("static")

	subrouter.HandleFunc("/open-api.json", func(w http.ResponseWriter, r *http.Request) {
		err := tplOpenApi.Execute(w, tplConfig)
		if err != nil {
			log.Println(err)
		}
	})
	subrouter.HandleFunc("/swagger-ui/swagger-initializer.js", func(w http.ResponseWriter, r *http.Request) {
		err := tplSwaggerUIInit.Execute(w, tplConfig)
		if err != nil {
			log.Println(err)
		}
	})

	subrouter.HandleFunc("/actuator/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(healthCheckResponse))
	}).Methods("GET")

	subrouter.HandleFunc("/api/example", Protect(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(apiResponse))
	}, []string{"first-role"}, true)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", stripTrailingSlashesMiddleware(handlers.LoggingHandler(os.Stdout, r))))
}

func stripTrailingSlashesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Index(r.URL.Path, "/swagger-ui/") != -1 {
			next.ServeHTTP(w, r)
		}
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}
