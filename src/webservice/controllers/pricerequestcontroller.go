package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"theorylab.com/shopwatch/internal/pkg/utils"

	"theorylab.com/shopwatch/internal/pkg/parser"
	"theorylab.com/shopwatch/src/webservice/models"
)

// isInit: Flag to mark if the initialization was done
var isInit bool = false

// PriceRequestController structure
type PriceRequestController struct {
	priceRequestPattern *regexp.Regexp
}

func newPriceRequestController() *PriceRequestController {
	return &PriceRequestController{
		priceRequestPattern: regexp.MustCompile(`^/pricecrequests/(\d+)/?`),
	}
}

// init function
// This has to be called once at the beginning to initialize the shop watch parser package
func (pr PriceRequestController) init() (bool, error) {
	// Init the parser package
	var doesExist bool = false
	currDir, err := os.Getwd()
	fmt.Println("Current Working Dir", currDir)

	merchantFile := "../app/merchants.json"
	doesExist, err = utils.CheckExists(merchantFile)
	if err != nil {
		return false, err
	}

	if !doesExist {
		return false, nil
	}

	merchantJSON, err := ioutil.ReadFile(merchantFile)
	if err != nil {
		return false, err
	}

	parser.LoadConfig(merchantJSON)
	isInit = true
	return isInit, nil
}

func (pr PriceRequestController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/pricerequests" {
		switch r.Method {
		case http.MethodPost:
			pr.processRequest(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (pr *PriceRequestController) processRequest(w http.ResponseWriter, r *http.Request) {
	p, err := pr.parsePriceRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse Price Request object"))
		return
	}

	if len(p.URL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No Url to request price"))
		return
	}

	if !isInit {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Price Parser has not been initialized"))
		return
	}

	productInfoPtr, err := parser.GetProductInfoByURL(p.URL)

	// for testing
	//var productInfo parser.Product
	//productInfo.Name = "Ipad 11"
	//productInfo.Img = []byte{1, 2, 3, 4}
	//productInfo.Available = "Available"
	//productInfo.Price = []float64{100.12, 23.25}
	//productInfo.Currency = "$"
	//productInfo.HTMLBody = "abcd"
	//productInfoPtr := &productInfo
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

func (pr *PriceRequestController) parsePriceRequest(r *http.Request) (models.PriceRequest, error) {
	dec := json.NewDecoder(r.Body)
	var p models.PriceRequest
	err := dec.Decode(&p)
	if err != nil {
		return models.PriceRequest{}, err
	}
	return p, nil
}
