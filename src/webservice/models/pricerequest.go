package models

import (
	"encoding/json"
	"errors"
	"os"
)

type PriceRequest struct {
	URL      string `json:"Url"`
	MinPrice int    `json:"Minprice"`
	MaxPrice int    `json:"Maxprice"`
	User     string `json:"User"`
}

func AddPriceRequest(p PriceRequest, prod *ProductInfo) (int, error) {
	var prodName string
	numBytesWritten := 0
	pstring, err := json.Marshal(p)
	if err != nil {
		return numBytesWritten, errors.New("Unable to convert price request to string to save")
	}

	pInfoString, err := json.Marshal(*prod)
	if err != nil {
		return numBytesWritten, errors.New("Unable to convert product info to string to save")
	}
	prodName = (*prod).Name
	if len(prodName) > 0 {
		// Saving the Request
		fullPathReqFileName := "./Req_" + (*prod).Name
		f, err := os.Create(fullPathReqFileName)
		if err != nil {
			return numBytesWritten, errors.New("Unable to create file " + fullPathReqFileName)
		}
		numBytesWritten, err = f.Write(pstring)
		if err != nil {
			return numBytesWritten, errors.New("Error saving file " + fullPathReqFileName)
		}

		// Saving the Product
		fullPathProdFileName := "./Prod_" + prod.Name
		f, err = os.Create(fullPathProdFileName)
		if err != nil {
			return numBytesWritten, errors.New("Unable to create file " + fullPathProdFileName)
		}

		numBytesWritten, err = f.Write(pInfoString)
		if err != nil {
			return numBytesWritten, errors.New("Error saving file " + fullPathReqFileName)
		}
	}
	return numBytesWritten, err
}
