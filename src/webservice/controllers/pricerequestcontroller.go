package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"theorylab.com/shopwatch/internal/pkg/parser"

	"theorylab.com/shopwatch/src/webservice/models"
)

type PriceRequestController struct {
	priceRequestPattern *regexp.Regexp
}

func newPriceRequestController() *PriceRequestController {
	return &PriceRequestController{
		priceRequestPattern: regexp.MustCompile(`^/pricecrequests/(\d+)/?`),
	}
}

func (pr PriceRequestController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/pricerequests" {
		switch r.Method {
		case http.MethodPost:
			pr.ProcessRequest(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (pr *PriceRequestController) ProcessRequest(w http.ResponseWriter, r *http.Request) {
	p, err := pr.ParsePriceRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse Price Request object"))
		return
	}

	if len(p.Url) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No Url to request price"))
		return
	}

	//productInfoPtr, err := parser.GetProductInfo(p.Url)

	// for testing
	var productInfo parser.Product
	productInfo.Name = "Ipad 11"
	productInfo.Img = []byte{1, 2, 3, 4}
	productInfo.Available = "Available"
	productInfo.Price = []float64{100.12, 23.25}
	productInfo.Currency = "$"
	productInfo.HTMLBody = "abcd"
	productInfoPtr := &productInfo
	// end testing

	if err != nil || productInfoPtr == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to retrieve product info"))
		return
	}

	pInfo, err := models.CreateProductInfoObject(*productInfoPtr)
	if err != nil || productInfoPtr == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to create product info"))
		return
	}

	charCount, err := models.AddPriceRequest(p, pInfo)
	if err != nil || charCount == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not add Price Request"))
		return
	}

	if charCount == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error writing to file"))
		return
	}
	encodeResponseAsJSON(p, w)
}

func (pr *PriceRequestController) ParsePriceRequest(r *http.Request) (models.PriceRequest, error) {
	dec := json.NewDecoder(r.Body)
	var p models.PriceRequest
	err := dec.Decode(&p)
	if err != nil {
		return models.PriceRequest{}, err
	}
	return p, nil
}
