package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Route defines a header to be checked for actions, and an array of actions to handle.
type Route struct {
	HeaderName string
	Actions    []Action
}

// Action contains information about how a specified action should be handled.
type Action struct {
	ActionName          string
	Callback            func(e *Envelope)
	NeedsAuthentication bool
	ServiceType         string
}

// NewRoute produces a new route struct with appropriate header defaults.
func NewRoute() Route {
	return Route{
		HeaderName: "SOAPAction",
	}
}

// RoutingGroup defines a group of actions for a given service type.
type RoutingGroup struct {
	Route       *Route
	ServiceType string
}

// HandleGroup returns a routing group type for the given service type.
func (r *Route) HandleGroup(serviceType string) RoutingGroup {
	return RoutingGroup{
		Route:       r,
		ServiceType: serviceType,
	}
}

// Unauthenticated associates an action to a function to be handled without authentication.
func (r *RoutingGroup) Unauthenticated(action string, function func(e *Envelope)) {
	r.Route.Actions = append(r.Route.Actions, Action{
		ActionName:          action,
		Callback:            function,
		NeedsAuthentication: false,
		ServiceType:         r.ServiceType,
	})
}

// Authenticated associates an action to a function to be handled with authentication.
func (r *RoutingGroup) Authenticated(action string, function func(e *Envelope)) {
	r.Route.Actions = append(r.Route.Actions, Action{
		ActionName:          action,
		Callback:            function,
		NeedsAuthentication: true,
		ServiceType:         r.ServiceType,
	})
}

func (route *Route) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s via %s", r.Method, r.URL, r.Host)

		// Check if there's a header of the type we need.
		service, actionName := parseAction(r.Header.Get("SOAPAction"))
		if service == "" || actionName == "" || r.Method != "POST" {
			printError(w, "WiiSOAP can't handle this. Try again later.")
			return
		}

		// Verify this is a service type we know.
		switch service {
		case "ecs":
		case "ias":
			break
		default:
			printError(w, "Unsupported service type...")
			return
		}

		fmt.Println("[!] Incoming " + strings.ToUpper(service) + " request - handling for " + actionName)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			printError(w, "Error reading request body...")
			return
		}

		// Ensure we can route to this action before processing.
		// Search all registered actions and find a matching action.
		var action Action
		for _, routeAction := range route.Actions {
			if routeAction.ActionName == actionName && routeAction.ServiceType == service {
				action = routeAction
			}
		}

		// Action is only properly populated if we found it previously.
		if action.ActionName == "" && action.ServiceType == "" {
			printError(w, "WiiSOAP can't handle this. Try again later.")
			return
		}

		fmt.Println("Received:", string(body))

		// Insert the current action being performed.
		e, err := NewEnvelope(service, actionName, body)
		if err != nil {
			printError(w, "Error interpreting request body: "+err.Error())
			return
		}

		// Check for authentication.
		if action.NeedsAuthentication {
			success, err := checkAuthentication(e)
			// Catch-all in case of invalid formatting or true invalidity.
			if !success || (err != nil) {
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}
		}

		// Call this action.
		action.Callback(e)

		// The action has now finished its task, and we can serialize.
		// Output may or may not truly be XML depending on where things failed.
		// We'll expect the best, however.
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		success, contents := e.becomeXML()
		if !success {
			// This is not what we wanted, and we need to reflect that.
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(contents))
		fmt.Println("Response:", contents)
	})
}

var routeVerify *sql.Stmt

func routeInitialize() {
	var err error
	// What we select does not matter.
	routeVerify, err = db.Prepare(`SELECT DeviceId FROM userbase WHERE DeviceToken=? AND AccountId=? AND DeviceId=?`)
	if err != nil {
		log.Fatalf("route initialize: error preparing statement: %v\n", err)
	}
}

// checkAuthentication validates various factors from a given request requiring authentication.
func checkAuthentication(e *Envelope) (bool, error) {
	// Get necessary authentication identifiers.
	deviceToken, err := getKey(e.doc, "DeviceToken")
	if err != nil {
		return false, err
	}
	accountId, err := getKey(e.doc, "AccountId")
	if err != nil {
		return false, err
	}

	hash := validateTokenFormat(deviceToken)
	if hash == "" {
		return false, nil
	}

	// The Wii sends us a MD5 of its given token. We store a SHA-256 hash of that.
	hashedDeviceToken := fmt.Sprintf("%x", sha256.Sum256([]byte(hash)))

	// Check using various input given.
	row := routeVerify.QueryRow(hashedDeviceToken, accountId, e.DeviceId())

	var throwaway string
	err = row.Scan(&throwaway)

	if err == sql.ErrNoRows {
		return false, err
	} else if err != nil {
		// We shouldn't encounter other errors.
		log.Printf("error occurred while checking authentication: %v\n", err)
		return false, err
	} else {
		return true, nil
	}
}

// validateTokenFormat confirms the prefix and size of tokens,
// in a format such as WT-5d41402abc4b2a76b9719d911017c592.
// It returns an empty string on failure.
func validateTokenFormat(token string) string {
	success := len(token) == 35 && token[:3] == "WT-"

	if success {
		// Strips the WT- prefix off.
		return token[3:35]
	} else {
		return ""
	}
}

func printError(w http.ResponseWriter, reason string) {
	http.Error(w, reason, http.StatusInternalServerError)
	fmt.Println("Failed to handle request: " + reason)
}
