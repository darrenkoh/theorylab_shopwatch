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
	"github.com/antonmedv/expr"
	"golang.org/x/net/html"
)

// Config Structure
type Config struct {
	Merchants []Merchant `json:"merchant"`
}

// Img Structure
type Img struct {
	Path          string `json:"path"`
	Transformeval string `json:"transformeval"`
}

// Xpath Structure
type Xpath struct {
	ProductName string   `json:"productname"`
	Img         Img      `json:"img"`
	Available   string   `json:"available"`
	Price       []string `json:"price"`
}

// Merchant Structure
type Merchant struct {
	Name             string `json:"name"`
	Xpath            Xpath  `json:"xpath"`
	AvailableKeyword string `json:"availablekeyword,omitempty"`
}

// Product Structure
type Product struct {
	Name      string
	Img       []byte
	Available string
	Price     []float64
	Currency  string
	HTMLBody  string
}

var config Config

// LoadConfig Load Merchant Configuration
func LoadConfig(data []byte) {
	err := json.Unmarshal(data, &config)
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

func getImg(url string, transformeval string) (img []byte, err error) {

	// We might need to transform the source URL
	if len(transformeval) > 0 {
		env := map[string]interface{}{
			"s":     url,
			"trim":  strings.Trim,
			"split": strings.Split,
		}

		program, err := expr.Compile(transformeval, expr.Env(env))
		if err != nil {
			return nil, fmt.Errorf("GET Img Error: %v", err)
		}

		output, err := expr.Run(program, env)
		if err != nil {
			return nil, fmt.Errorf("GET Img Error: %v", err)
		}

		url = output.(string)
	}

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
		idx := 1
		if len(prc) == 1 {
			idx = 0
		}
		f, _ := strconv.ParseFloat(strings.Replace(prc[idx], ",", "", -1), 64)
		return f
	}
	return -1
}

func getAvailability(available string, config Merchant) string {
	if strings.Contains(strings.ToLower(available), strings.ToLower(config.AvailableKeyword)) {
		return "Yes"
	}
	return "No"
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

	client := &http.Client{
		Timeout: 30 * time.Second,
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
	//ioutil.WriteFile("test.html", data, 0644)
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
	product.Name = prodName.Data

	// Download the Image
	node, err := htmlquery.Query(doc, merchantConfig.Xpath.Img.Path)
	if err != nil {
		fmt.Printf("Invalid Xpath: %v", err)
		return nil, err
	}
	imgURL := getAttr("src", node.Attr)
	img, err := getImg(imgURL, merchantConfig.Xpath.Img.Transformeval)
	product.Img = img

	// Stock Available
	product.Available = "Yes"
	stock, err := htmlquery.Query(doc, merchantConfig.Xpath.Available)
	if err != nil {
		fmt.Printf("Invalid Xpath: %v", err)
		return nil, err
	}
	if stock != nil {
		product.Available = getAvailability(stock.Data, merchantConfig)
	}

	// Get the price
	// 1.0 collect all prices
	product.Currency = "USD"
	for _, xpath := range merchantConfig.Xpath.Price {
		price, err := htmlquery.Query(doc, xpath)
		if err != nil {
			fmt.Printf("Invalid Xpath: %v", err)
		} else if price != nil {
			product.Price = append(product.Price, getPrice(price.Data))
		}
	}

	return &product, nil
}
