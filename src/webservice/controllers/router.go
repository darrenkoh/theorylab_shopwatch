package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

// RegisterControllers method registers all controllers will be used for the webservice
func RegisterControllers() {
	pr := newPriceRequestController()

	http.Handle("/pricerequests", *pr)
	http.Handle("/pricerequests/", *pr)
}

func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
