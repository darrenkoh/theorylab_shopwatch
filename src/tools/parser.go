package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// Config Structure
type Config struct {
	Merchants []struct {
		Name  string `json:"name"`
		Xpath struct {
			Img       string `json:"img"`
			Available string `json:"available"`
			Price     string `json:"price"`
		} `json:"xpath"`
	} `json:"merchants"`
}

// Product Structure
type Product struct {
	Img       []byte
	Available string
	Price     float64
	Currency  string
	HtmlBody  string
}

var config Config

func init() {
	merchantJSON, err := ioutil.ReadFile("merchants.json")
	err = json.Unmarshal(merchantJSON, &config)
	if err != nil {
		fmt.Printf("Unable to load merchants.json: %v", err)
	}
}

func getHTML(url string) (body string, err error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", fmt.Errorf("GET Error: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Read Body: %v", err)
	}

	return string(data), nil
}

func getAttr(key string, attrs []html.Attribute) string {
	for _, a := range attrs {
		if key == a.Key {
			return a.Val
		}
	}
	return ""
}

func getImg(url string) (img []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET Error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func getPrice(price string) float64 {
	prc := regexp.MustCompile("[$]").Split(price, 2)
	if prc != nil {
		f, _ := strconv.ParseFloat(strings.Replace(prc[1], ",", "", -1), 64)
		return f
	}
	return -1
}

// GetProductInfo load url and return Product
func GetProductInfo(url string) (info *Product, err error) {
	var product Product
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}
	/* Write to temp for testing */
	//ioutil.WriteFile("test2.html", data, 0644)
	docReader := bytes.NewReader(data)
	doc, err := htmlquery.Parse(docReader)
	if err != nil {
		return nil, err
	}

	// Download the Image
	node, err := htmlquery.Query(doc, config.Merchants[0].Xpath.Img)
	if err != nil {
		fmt.Printf("Invalid Xpath: %v", err)
		return nil, err
	}
	imgURL := getAttr("src", node.Attr)
	img, err := getImg(imgURL)
	product.Img = img

	/* Write to file for testing */
	//ioutil.WriteFile("a.jpg", img, 0644)

	// Stock Available
	stock, err := htmlquery.Query(doc, config.Merchants[0].Xpath.Available)
	if err != nil {
		fmt.Printf("Invalid Xpath: %v", err)
		return nil, err
	}
	product.Available = stock.FirstChild.Data
	print("Stock Status: ")
	println(stock.FirstChild.Data)

	// Get the price
	price, err := htmlquery.Query(doc, config.Merchants[0].Xpath.Price)
	if err != nil {
		fmt.Printf("Invalid Xpath: %v", err)
	}
	if price != nil {
		product.Price = getPrice(price.FirstChild.Data)
		fmt.Printf("Price: %v", product.Price)
	} else {
		product.Price = -1
		println("NA")
	}

	return &product, nil
}
