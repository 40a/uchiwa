package uchiwa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sensu/uchiwa/uchiwa/auth"
	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// FilterAggregates is a function that filters aggregates
var FilterAggregates func(aggregates *[]interface{}, token *jwt.Token) []interface{}

// FilterChecks is a function that filters checks
var FilterChecks func(checks *[]interface{}, token *jwt.Token) []interface{}

// FilterClients is a function that filters clients
var FilterClients func(clients *[]interface{}, token *jwt.Token) []interface{}

// FilterAggregates is a function that filters datacenters
var FilterDatacenters func(datacenters []*structs.Datacenter, token *jwt.Token) []*structs.Datacenter

// FilterEvents is a function that filters events
var FilterEvents func(events *[]interface{}, token *jwt.Token) []interface{}

// FilterAggregates is a function that filters aggregates
var FilterStashes func(stashes *[]interface{}, token *jwt.Token) []interface{}

// FilterAggregates is a function that filters aggregates
var FilterSubscriptions func(subscriptions *[]string, token *jwt.Token) []string

// FilterGetRequest is a function that filters GET requests.
var FilterGetRequest func(string, *jwt.Token) bool

// FilterPostRequest is a function that filters POST requests.
var FilterPostRequest func(*jwt.Token, *interface{}) bool

// FilterSensuDataData is a function that filters Sensu Data.
var FilterSensuData func(*jwt.Token, *structs.Data) *structs.Data

// Aggregates
func (u *Uchiwa) aggregatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// verify that the authenticated user is authorized to access this resource
	token := auth.GetTokenFromContext(r)

	resources := strings.Split(r.URL.Path, "/")

	if len(resources) == 2 {
		token := auth.GetTokenFromContext(r)
		aggregates := FilterAggregates(&u.Data.Aggregates, token)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(aggregates); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	} else if len(resources) == 5 {
		check := resources[3]
		dc := resources[2]
		issued := resources[4]

		if check == "" || dc == "" || issued == "" {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		unauthorized := FilterGetRequest(dc, token)
		if unauthorized {
			http.Error(w, fmt.Sprint(""), http.StatusNotFound)
			return
		}

		aggregate, err := u.GetAggregateByIssued(check, issued, dc)
		if err != nil {
			http.Error(w, fmt.Sprint(err), 500)
			return
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(aggregate); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

// Checks
func (u *Uchiwa) checksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := auth.GetTokenFromContext(r)
	checks := FilterChecks(&u.Data.Checks, token)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(checks); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

// Clients
func (u *Uchiwa) clientsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" && r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := auth.GetTokenFromContext(r)

	resources := strings.Split(r.URL.Path, "/")
	if len(resources) == 2 && r.Method == "GET" {
		clients := FilterClients(&u.Data.Clients, token)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(clients); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	} else if len(resources) != 4 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	client := resources[3]
	dc := resources[2]

	if client == "" || dc == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if r.Method == "DELETE" {
		// verify that the authenticated user is authorized to access this resource
		unauthorized := FilterGetRequest(dc, token)
		if unauthorized {
			http.Error(w, fmt.Sprint(""), http.StatusNotFound)
			return
		}

		err := u.DeleteClient(client, dc)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "GET" {
		unauthorized := FilterGetRequest(dc, token)
		if unauthorized {
			http.Error(w, fmt.Sprint(""), http.StatusNotFound)
			return
		}

		data, err := u.GetClient(client, dc)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusNotFound)
			return
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(data); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// Config
func (u *Uchiwa) configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resources := strings.Split(r.URL.Path, "/")

	if len(resources) == 2 {
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(u.PublicConfig); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		if resources[2] == "auth" {
			fmt.Fprintf(w, "%s", u.PublicConfig.Uchiwa.Auth)
		} else {
			http.Error(w, "", http.StatusNotFound)
			return
		}
	}
}

// Datacenters
func (u *Uchiwa) datacentersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := auth.GetTokenFromContext(r)
	datacenters := FilterDatacenters(u.Data.Dc, token)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(datacenters); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

// Events
func (u *Uchiwa) eventsHandler(w http.ResponseWriter, r *http.Request) {
	token := auth.GetTokenFromContext(r)

	if r.Method == "DELETE" {
		resources := strings.Split(r.URL.Path, "/")
		if len(resources) != 5 {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		check := resources[4]
		client := resources[3]
		dc := resources[2]

		unauthorized := FilterGetRequest(dc, token)
		if unauthorized {
			http.Error(w, fmt.Sprint(""), http.StatusNotFound)
			return
		}

		err := u.ResolveEvent(check, client, dc)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

	} else if r.Method == "GET" {
		events := FilterEvents(&u.Data.Events, token)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(events); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func (u *Uchiwa) getSensuHandler(w http.ResponseWriter, r *http.Request) {
	token := auth.GetTokenFromContext(r)
	data := FilterSensuData(token, u.Data)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

func (u *Uchiwa) healthHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	var err error
	if r.URL.Path[1:] == "health/sensu" {
		err = encoder.Encode(u.Data.Health.Sensu)
	} else if r.URL.Path[1:] == "health/uchiwa" {
		err = encoder.Encode(u.Data.Health.Uchiwa)
	} else {
		err = encoder.Encode(u.Data.Health)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
	}
}

// Results
func (u *Uchiwa) resultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := auth.GetTokenFromContext(r)
	results := FilterChecks(&u.Data.Results, token)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(results); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

// Stashes
func (u *Uchiwa) stashesHandler(w http.ResponseWriter, r *http.Request) {
	token := auth.GetTokenFromContext(r)
	resources := strings.Split(r.URL.Path, "/")

	if r.Method == "DELETE" && len(resources) >= 3 {
		dc := resources[2]
		path := strings.Join(resources[3:], "/")

		unauthorized := FilterGetRequest(dc, token)
		if unauthorized {
			http.Error(w, fmt.Sprint(""), http.StatusNotFound)
			return
		}

		err := u.DeleteStash(dc, path)
		if err != nil {
			http.Error(w, "Could not create the stash", http.StatusNotFound)
		}
	} else if r.Method == "GET" {
		stashes := FilterStashes(&u.Data.Stashes, token)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(stashes); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var data stash
		err := decoder.Decode(&data)
		if err != nil {
			http.Error(w, "Could not decode body", http.StatusInternalServerError)
			return
		}

		// verify that the authenticated user is authorized to access this resource
		unauthorized := FilterGetRequest(data.Dc, token)
		if unauthorized {
			http.Error(w, fmt.Sprint(""), http.StatusNotFound)
			return
		}

		if token != nil && token.Claims["Username"] != nil {
			data.Content["username"] = token.Claims["Username"]
		}

		err = u.PostStash(data)
		if err != nil {
			http.Error(w, "Could not create the stash", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

// Subscriptions
func (u *Uchiwa) subscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token := auth.GetTokenFromContext(r)
	subscriptions := FilterSubscriptions(&u.Data.Subscriptions, token)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(subscriptions); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
		return
	}
}

// WebServer starts the web server and serves GET & POST requests
func (u *Uchiwa) WebServer(publicPath *string, auth auth.Config) {

	// Private endpoints
	http.Handle("/aggregates", auth.Authenticate(http.HandlerFunc(u.aggregatesHandler)))
	http.Handle("/aggregates/", auth.Authenticate(http.HandlerFunc(u.aggregatesHandler)))
	http.Handle("/checks", auth.Authenticate(http.HandlerFunc(u.checksHandler)))
	http.Handle("/clients", auth.Authenticate(http.HandlerFunc(u.clientsHandler)))
	http.Handle("/clients/", auth.Authenticate(http.HandlerFunc(u.clientsHandler)))
	http.Handle("/config", auth.Authenticate(http.HandlerFunc(u.configHandler)))
	http.Handle("/datacenters", auth.Authenticate(http.HandlerFunc(u.datacentersHandler)))
	http.Handle("/events", auth.Authenticate(http.HandlerFunc(u.eventsHandler)))
	http.Handle("/events/", auth.Authenticate(http.HandlerFunc(u.eventsHandler)))
	http.Handle("/results", auth.Authenticate(http.HandlerFunc(u.resultsHandler)))
	http.Handle("/stashes", auth.Authenticate(http.HandlerFunc(u.stashesHandler)))
	http.Handle("/stashes/", auth.Authenticate(http.HandlerFunc(u.stashesHandler)))
	http.Handle("/subscriptions", auth.Authenticate(http.HandlerFunc(u.subscriptionsHandler)))

	// Legacy private endpoint
	http.Handle("/get_sensu", auth.Authenticate(http.HandlerFunc(u.getSensuHandler)))

	// Static files
	http.Handle("/", http.FileServer(http.Dir(*publicPath)))

	// Public endpoints
	http.Handle("/config/", http.HandlerFunc(u.configHandler))
	http.Handle("/health", http.HandlerFunc(u.healthHandler))
	http.Handle("/health/", http.HandlerFunc(u.healthHandler))
	http.Handle("/login", auth.GetIdentification())

	listen := fmt.Sprintf("%s:%d", u.Config.Uchiwa.Host, u.Config.Uchiwa.Port)
	logger.Infof("Uchiwa is now listening on %s", listen)
	logger.Fatal(http.ListenAndServe(listen, nil))
}
