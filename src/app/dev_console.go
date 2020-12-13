package main

import (
	"fmt"
	"io/ioutil"

	"theorylab.com/shopwatch/internal/pkg/parser"
)

func main() {

	fmt.Println("Going...")

	merchantJSON, err := ioutil.ReadFile("merchants.json")
	if err != nil {
		panic(fmt.Errorf("Loading configuration failed %v", err))
	}
	parser.LoadConfig(merchantJSON)

	//url := "https://www.bhphotovideo.com/c/product/1595098-REG/sony_3005726_hd_camera_for_playstation.html"
	//url := "https://www.bhphotovideo.com/c/product/1600080-REG/microsoft_rrt_00001_xbox_series_x_1tb.html"
	//url := "https://www.bhphotovideo.com/c/product/1543054-REG/elmo_1430_mx_p2_visual_presenter_and.html"
	//url := "https://www.bhphotovideo.com/c/product/1555040-REG/lg_50un7300puf_un7300_50_class_hdr.html"
	//url := "https://www.bhphotovideo.com/c/product/1604841-REG/apple_z125_mgn7_05_bh_13_3_macbook_air_with.html"
	//url := "https://www.bestbuy.com/site/apple-airpods-with-charging-case-latest-model-white/6084400.p?skuId=6084400"
	//url := "https://www.bestbuy.com/site/sony-playstation-5-console/6426149.p?skuId=6426149"
	//url := "https://www.bestbuy.com/site/apple-iphone-12-5g-64gb-black-verizon/6009897.p?skuId=6009897"
	//url := "https://www.walmart.com/ip/Echelon-Connect-Sport-Indoor-Cycling-Exercise-Bike-with-6-Month-Free-Membership-120-value/533034706"
	//url := "https://www.walmart.com/ip/LEGO-Star-Wars-A-wing-Starfighter-75275-Building-Toy-Cool-Gift-Idea-for-Creative-Adults-1-673-Pieces/554462454"
	//url := "https://www.amazon.com/dp/B07XJ8C8F7"
	//url := "https://www.amazon.com/dp/B076PRWVFG"
	//url := "https://www.amazon.com/PULSE-3D-Wireless-Headset-PlayStation-5/dp/B08FC6QLKN"
	//url := "https://www.newegg.com/grey-with-blue-diamond-cut-msi-prestige-15-a10sc-296-mainstream/p/N82E16834155448?Item=N82E16834155448"
	//url := "https://www.adorama.com/nolpminimk2.html"
	//url := "https://www.adorama.com/icalpe6nh.html"
	url := "https://www.etsy.com/listing/505391661/cat-beginner-diy-felting-kit-wool-felt?ga_order=most_relevant&ga_search_type=all&ga_view_type=gallery&ga_search_query=&ref=sc_gallery-1-4&plkey=8a8b101dcb15ca341ee1b9e9f5fb528119e7b7d6%3A505391661&frs=1&bes=1"

	p, err := parser.GetProductInfoByURL(url)
	if p != nil && err == nil {
		fmt.Printf("Name: %s\n", p.Name)
		fmt.Printf("Price: %v\n", p.Price)
		fmt.Printf("Currency: %s\n", p.Currency)
		fmt.Printf("Available: %s\n", p.Available)
		fmt.Printf("Img Length: %d\n", len(p.Img))
		fmt.Printf("Body Length: %d\n", len(p.HTMLBody))
		ioutil.WriteFile("test.html", []byte(p.HTMLBody), 0644)
		ioutil.WriteFile("test.jpg", p.Img, 0644)
	}
}
