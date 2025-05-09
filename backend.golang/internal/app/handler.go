package app

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/amadeus"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/flightsky"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/googleflights"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/redis"
	"github.com/rubengp99/golang-flights-challenge/internal/workflow"
	"github.com/rubengp99/golang-flights-challenge/pkg"
)

func defaultOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Max-Age", "600")
	serveResponse(nil, http.StatusNoContent, w)
}

func getQueryParams(value interface{}, r *http.Request) error {
	decoder := schema.NewDecoder()
	// decoder lookup for values on the json tag, instead of the default schema tag
	decoder.SetAliasTag("json")
	decoder.IgnoreUnknownKeys(true)

	// for demo purposes, we need to simplify our date parsing to ignore time
	decoder.RegisterConverter(time.Time{}, func(value string) reflect.Value {
		if v, err := time.Parse("2006-01-02", value); err == nil {
			return reflect.ValueOf(v)
		}
		return reflect.Value{} // this is the same as the private const invalidType
	})

	if err := decoder.Decode(value, r.URL.Query()); err != nil {
		return errors.Wrap(err, "handler - failed to decode query params")
	}

	return nil
}

func createToken(secretKey, clientID string) (string, int64, error) {
	expiration := time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": clientID,
			"exp":      expiration,
		})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", expiration, err
	}

	return tokenString, expiration, nil
}

// RetrieveBestFlightsHandler handles best flights lookup
func RetrieveBestFlightsHandler(redisClient redis.Service,
	googleflightService googleflights.Service,
	amadeusService amadeus.Service,
	flightskyService flightsky.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params pkg.QueryParams
		if err := getQueryParams(&params, r); err != nil {
			serveResponse(newError(err.Error()), http.StatusInternalServerError, w)
			return
		}

		if err := validateBestFlightsParamsRequest(params); err != nil {
			serveResponse(newError(err.Error()), http.StatusBadRequest, w)
			return
		}

		wf := workflow.RetrieveBestFlights(redisClient, googleflightService, amadeusService, flightskyService)
		res, err := wf(params)
		if err != nil {
			serveResponse(newError(err.Error()), http.StatusInternalServerError, w)
			return
		}
		serveResponse(res, http.StatusOK, w)
	})
}

// LoginHandler represents login handler functionality
func LoginHandler(appCreds pkg.CrendetialsRequest, secretKey string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var u pkg.CrendetialsRequest
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			serveResponse(newError(err.Error()), http.StatusInternalServerError, w)
			return
		}

		if err := validateCrendetialsRequest(u); err != nil {
			serveResponse(newError(err.Error()), http.StatusBadRequest, w)
			return
		}

		if u.ClientID == appCreds.ClientID && u.ClientSecret == appCreds.ClientSecret {
			token, exp, err := createToken(secretKey, u.ClientID)
			if err != nil {
				serveResponse(newError(err.Error()), http.StatusInternalServerError, w)
				return
			}

			response := pkg.CredentialsResponse{
				AccessToken: token,
				ExpIn:       exp,
			}

			serveResponse(response, http.StatusOK, w)
			return
		}

		serveResponse(newError("Invalid Credentials"), http.StatusUnauthorized, w)
	})
}

// SubcribeToFlightOfferUpdatesHandler handles periodic updates to a flight search criteria using websockets
func SubcribeToFlightOfferUpdatesHandler(redisClient redis.Service,
	googleflightService googleflights.Service,
	amadeusService amadeus.Service,
	flightskyService flightsky.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we need mandatory params in order to subscribe to updates for an specific request
		var params pkg.QueryParams
		if err := getQueryParams(&params, r); err != nil {
			serveResponse(newError(err.Error()), http.StatusInternalServerError, w)
			return
		}

		if err := validateBestFlightsParamsRequest(params); err != nil {
			serveResponse(newError(err.Error()), http.StatusBadRequest, w)
			return
		}

		// we need mandatory params in order to subscribe to updates for an specific request
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				wf := workflow.RetrieveBestFlights(redisClient, googleflightService, amadeusService, flightskyService)
				res, err := wf(params)
				if err != nil {
					serveResponse(newError(err.Error()), http.StatusInternalServerError, w)
					return
				}

				if err := conn.WriteJSON(res); err != nil {
					log.Println("Write error:", err)
					return
				}
			}
		}
	})
}
