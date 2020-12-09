package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

func RegisterControllers() {
	pr := newPriceRequestController()

	http.Handle("/pricerequests", *pr)
	http.Handle("/pricerequests/", *pr)
}

func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
