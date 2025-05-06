package vendors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	ContentTypeJSON       = "application/json"
	ContentTypeURLEncoded = "application/x-www-form-urlencoded"
)

// Config represents a generic config/credentials setup for third party integrations
type Config struct {
	BaseURL      string
	APIKey       string
	ClientID     string
	ClientSecret string
}

// Service represents a generic http service interface
type Service interface {
	Authenticate(req *http.Request) error
	Client() *http.Client
}

// Request is a request representation for integrations
type Request struct {
	SkipAuth    bool
	ContentType string
	BaseURL     string
	Resource    string
	Method      string
	Params      url.Values
	Payload     any
}

func (r Request) URL() string {
	return fmt.Sprintf("%v/%s?%s", r.BaseURL, r.Resource, r.Params.Encode())
}

func generateJSONBody(payload any) (io.Reader, error) {
	if payload == nil {
		return nil, nil
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "client - unable to marshal request struct")
	}

	return bytes.NewReader(requestBody), nil
}

func generateURLEncodedBody(payload any) (io.Reader, error) {
	if payload == nil {
		return nil, nil
	}

	requestBody, ok := payload.(url.Values)
	if !ok {
		return nil, fmt.Errorf("payload is not url.Values{} for content %s", ContentTypeURLEncoded)
	}

	return strings.NewReader(requestBody.Encode()), nil
}

// MakeHTTPRequest build, send and decode HTTP request/response made to an external service,
func MakeHTTPRequest(v Service, request Request, resp any) error {
	var (
		body io.Reader
		err  error
	)

	// if there's a request body, we add to the request payload, otherwise skip
	if request.ContentType == ContentTypeJSON {
		body, err = generateJSONBody(request.Payload)
	}

	if request.ContentType == ContentTypeURLEncoded {
		body, err = generateURLEncodedBody(request.Payload)
	}

	if err != nil {
		return errors.Wrapf(err, "client - unable to generate payload for content %s", request.ContentType)
	}

	log.Printf("making request to route %s with payload %s", request.URL(), request.Payload)

	// generate a new http request object
	req, err := http.NewRequest(request.Method, request.URL(), body)
	if err != nil {
		return errors.Wrap(err, "client - unable to create request body")
	}

	req.Header.Add("Content-Type", request.ContentType)

	// sets authentication headers to request
	if !request.SkipAuth {
		if err := v.Authenticate(req); err != nil {
			return err
		}
	}

	// make an http call to the outlaying service
	res, err := v.Client().Do(req)
	if err != nil {
		return errors.Wrap(err, "client - failed to execute request")
	}

	// read our response
	b, err := io.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		return errors.Wrap(err, "client - unable to read response body")
	}

	// validate response
	validResponses := map[int]bool{
		http.StatusOK:        true,
		http.StatusNoContent: true,
		http.StatusCreated:   true,
		http.StatusAccepted:  true,
	}

	if !validResponses[res.StatusCode] {
		return fmt.Errorf("invalid status code received, expected 200/204/201/202, got %v with body %s", res.StatusCode, b)
	}

	// do not unmarshal response on 204 or empty response
	if res.StatusCode == http.StatusNoContent || len(b) == 0 {
		fmt.Printf("got response %s code %d", string(b), res.StatusCode)
		return nil
	}

	// decode our response payload
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&resp); err != nil {
		fmt.Printf("got response %s code %d", string(b), res.StatusCode)
		return errors.Wrap(err, "unable to unmarshal response body")
	}

	return nil
}
