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

type Config struct {
	Merchants []Merchant `json:"merchant"`
}
type Img struct {
	SourceAttribute string `json:"sourceattribute"`
	Path            string `json:"path"`
	Transformeval   string `json:"transformeval"`
}
type Available struct {
	Path     string `json:"path"`
	Keyword  string `json:"keyword"`
	Operator string `json:"operator"`
}
type Xpath struct {
	Productname string      `json:"productname"`
	Img         Img         `json:"img"`
	Available   []Available `json:"available"`
	Price       []string    `json:"price"`
}
type RequestHeaderOverwrite struct {
	Cookie string `json:"cookie"`
}
type Merchant struct {
	Name                   string                 `json:"name"`
	Xpath                  Xpath                  `json:"xpath"`
	RequestHeaderOverwrite RequestHeaderOverwrite `json:"request_header_overwrite"`
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
	price = strings.TrimSpace(price)
	prc := regexp.MustCompile("[$]").Split(price, 2)
	if prc != nil {
		idx := 1
		if len(prc) == 1 {
			idx = 0
		}
		t := prc[idx]
		for _, x := range [3]string{",", "+"} {
			t = strings.Replace(t, x, "", -1)
		}
		f, _ := strconv.ParseFloat(t, 64)
		return f
	}
	return -1
}

func getAvailability(available string, config Available) string {
	source := strings.Trim(strings.TrimSpace(available), ".")
	if config.Operator == "equal" {
		if strings.ToLower(source) == strings.ToLower(config.Keyword) {
			return "Yes"
		}
		return "No"

	} else if config.Operator == "contain" {
		if strings.Contains(strings.ToLower(source), strings.ToLower(config.Keyword)) {
			return "Yes"
		}
		return "No"
	}
	return "No"
}

func getMerchantFromURL(url string) (merchant *Merchant) {
	// Detect Merchan
	merchantConfig := config.Merchants[0]
	for _, m := range config.Merchants {
		if strings.Contains(url, m.Name) {
			merchantConfig = m
			break
		}
	}
	return &merchantConfig
}

// GetProductInfoByURL load url and return Product
func GetProductInfoByURL(url string) (info *Product, err error) {

	merchantConfig := getMerchantFromURL(url)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Set("Referer", "https://www.google.com")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")

	if len(merchantConfig.RequestHeaderOverwrite.Cookie) > 0 {
		req.Header.Set("cookie", merchantConfig.RequestHeaderOverwrite.Cookie)
	}

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}
	html := string(data)

	/* Write to temp for testing */
	ioutil.WriteFile("test.html", data, 0644)

	return GetProductInfoByHTML(html, merchantConfig)
}

// GetProductInfoByHTML parse html content and return Product
func GetProductInfoByHTML(html string, merchantConfig *Merchant) (info *Product, err error) {

	product := Product{HTMLBody: html}

	docReader := bytes.NewReader([]byte(html))
	doc, err := htmlquery.Parse(docReader)
	if err != nil {
		return nil, err
	}

	// Product Name
	prodName, err := htmlquery.Query(doc, merchantConfig.Xpath.Productname)
	if err != nil {
		fmt.Printf("Invalid Xpath: %v", err)
		return nil, err
	}
	product.Name = strings.TrimSpace(prodName.Data)

	// Download the Image
	node, err := htmlquery.Query(doc, merchantConfig.Xpath.Img.Path)
	if err != nil {
		fmt.Printf("Invalid Xpath: %v", err)
		return nil, err
	}
	srcattribute := "src"
	if len(merchantConfig.Xpath.Img.SourceAttribute) > 0 {
		srcattribute = merchantConfig.Xpath.Img.SourceAttribute
	}

	imgURL := getAttr(srcattribute, node.Attr)
	img, err := getImg(imgURL, merchantConfig.Xpath.Img.Transformeval)
	product.Img = img

	// Stock Available
	product.Available = "No"
	for _, a := range merchantConfig.Xpath.Available {
		stock, err := htmlquery.Query(doc, a.Path)
		if err != nil {
			fmt.Printf("Invalid Xpath: %v", err)
			return nil, err
		}
		if stock != nil {
			product.Available = getAvailability(stock.Data, a)
		}
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
