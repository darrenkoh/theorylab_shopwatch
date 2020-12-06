package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// Merchant Structure
type Config struct {
	Merchants []struct {
		Name  string `json:"name"`
		Xpath struct {
			Img   string `json:"img"`
			Price string `json:"price"`
		} `json:"xpath"`
	} `json:"merchants"`
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

func main() {
	fmt.Println("Going...")
	var merchants Config
	url := "https://www.bhphotovideo.com/c/product/1595098-REG/sony_3005726_hd_camera_for_playstation.html"
	doc, err := htmlquery.LoadURL(url)
	merchantJson, err := ioutil.ReadFile("merchants.json")
	err = json.Unmarshal(merchantJson, &merchants)
	if err != nil {
		fmt.Println("error:", err)
	}

	node, err := htmlquery.Query(doc, merchants.Merchants[0].Xpath.Img)
	if err != nil {
		panic(`not a valid xpath`)
	}

	// Download the Image
	imgURL := getAttr("src", node.Attr)
	img, err := getImg(imgURL)
	ioutil.WriteFile("a.jpg", img, 0644)

	// Get the price
	price, err := htmlquery.Query(doc, merchants.Merchants[0].Xpath.Price)
	println(price.FirstChild.Data)
}
