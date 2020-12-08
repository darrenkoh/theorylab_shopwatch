package main

import (
	"fmt"
	"io/ioutil"

	parser "./src/tools"
)

func main() {

	fmt.Println("Going...")

	//url := "https://www.bhphotovideo.com/c/product/1595098-REG/sony_3005726_hd_camera_for_playstation.html"
	//url := "https://www.bhphotovideo.com/c/product/1600080-REG/microsoft_rrt_00001_xbox_series_x_1tb.html"
	//url := "https://www.bhphotovideo.com/c/product/1543054-REG/elmo_1430_mx_p2_visual_presenter_and.html"
	//url := "https://www.bhphotovideo.com/c/product/1555040-REG/lg_50un7300puf_un7300_50_class_hdr.html"
	//url := "https://www.bhphotovideo.com/c/product/1604841-REG/apple_z125_mgn7_05_bh_13_3_macbook_air_with.html"
	bestbuy := "https://www.bestbuy.com/site/apple-airpods-with-charging-case-latest-model-white/6084400.p?skuId=6084400"
	//bestbuy := "https://www.bestbuy.com/site/sony-playstation-5-console/6426149.p?skuId=6426149"
	//etsy := "https://www.etsy.com/listing/892352400/christmas-gift-wooden-name-puzzle?ga_order=most_relevant&ga_search_type=all&ga_view_type=gallery&ga_search_query=name+puzzle&ref=sc_gallery-1-1&plkey=1242ea853ed2f5cbc112f343f44ce2c469b91461%3A892352400&pro=1&frs=1&col=1"
	p, err := parser.GetProductInfo(bestbuy)

	if p != nil && err != nil {
		println(p)
		ioutil.WriteFile("bestbuy.html", []byte(p.HTMLBody), 0644)
	}
}
