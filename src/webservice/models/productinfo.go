package models

import (
	b64 "encoding/base64"

	"theorylab.com/shopwatch/internal/pkg/parser"
)

type ProductInfo struct {
	Name      string
	Image     string
	Available string
	Price     []float64
	Currency  string
	HTMLBody  string
}

func CreateProductInfoObject(p parser.Product) (prodInfo *ProductInfo, err error) {
	var pInfo ProductInfo

	if len(p.Name) > 0 {
		pInfo.Name = p.Name
	} else {
		pInfo.Name = "Unknown Product"
	}

	if p.Img != nil {
		pInfo.Image = b64.URLEncoding.EncodeToString(p.Img)
	}

	pInfo.Available = p.Available
	pInfo.Price = p.Price
	pInfo.Currency = p.Currency
	pInfo.HTMLBody = p.HTMLBody

	return &pInfo, nil
}
