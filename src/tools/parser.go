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
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// Config Structure
type Config struct {
	Merchants []struct {
		Name  string `json:"name"`
		Xpath struct {
			ProductName string `json:"productname"`
			Img         string `json:"img"`
			Available   string `json:"available"`
			Price       string `json:"price"`
		} `json:"xpath"`
	} `json:"merchants"`
}

// Product Structure
type Product struct {
	Name      string
	Img       []byte
	Available string
	Price     float64
	Currency  string
	HTMLBody  string
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

	// Detect Merchan
	merchantConfig := config.Merchants[0]
	for _, m := range config.Merchants {
		if strings.Contains(url, m.Name) {
			merchantConfig = m
			break
		}
	}

	var product Product
	//response, err := http.Get(url)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "PostmanRuntime/7.26.5")
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)

	product.HTMLBody = string(data)

	if err != nil {
		return nil, err
	}
	/* Write to temp for testing */
	ioutil.WriteFile("test.html", data, 0644)
	docReader := bytes.NewReader(data)
	doc, err := htmlquery.Parse(docReader)
	if err != nil {
		return nil, err
	}

	// Product Name
	prodName, err := htmlquery.Query(doc, merchantConfig.Xpath.ProductName)
	if err != nil {
		fmt.Printf("Invalid Xpath: %v", err)
		return nil, err
	}
	product.Name = prodName.FirstChild.Data
	print("Product Name: ")
	println(product.Name)

	// Download the Image
	node, err := htmlquery.Query(doc, merchantConfig.Xpath.Img)
	if err != nil {
		fmt.Printf("Invalid Xpath: %v", err)
		return nil, err
	}
	imgURL := getAttr("src", node.Attr)
	img, err := getImg(imgURL)
	product.Img = img

	/* Write to file for testing */
	ioutil.WriteFile("a.jpg", img, 0644)

	// Stock Available
	stock, err := htmlquery.Query(doc, merchantConfig.Xpath.Available)
	if err != nil {
		fmt.Printf("Invalid Xpath: %v", err)
		return nil, err
	}
	product.Available = stock.FirstChild.Data
	print("Stock Status: ")
	println(stock.FirstChild.Data)

	// Get the price
	price, err := htmlquery.Query(doc, merchantConfig.Xpath.Price)
	if err != nil {
		fmt.Printf("Invalid Xpath: %v", err)
	}
	if price != nil {
		product.Price = getPrice(price.FirstChild.Data)
		product.Currency = "USD"
		fmt.Printf("Price: %v", product.Price)
	} else {
		product.Price = -1
		println("NA")
	}

	return &product, nil
}
