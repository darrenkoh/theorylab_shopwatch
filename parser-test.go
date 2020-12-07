package main

import (
	"fmt"

	parser "./src/tools"
)

func main() {

	fmt.Println("Going...")

	//url := "https://www.bhphotovideo.com/c/product/1595098-REG/sony_3005726_hd_camera_for_playstation.html"
	//url := "https://www.bhphotovideo.com/c/product/1600080-REG/microsoft_rrt_00001_xbox_series_x_1tb.html"
	//url := "https://www.bhphotovideo.com/c/product/1543054-REG/elmo_1430_mx_p2_visual_presenter_and.html"
	//url := "https://www.bhphotovideo.com/c/product/1555040-REG/lg_50un7300puf_un7300_50_class_hdr.html"
	url := "https://www.bhphotovideo.com/c/product/1604841-REG/apple_z125_mgn7_05_bh_13_3_macbook_air_with.html"
	p, err := parser.GetProductInfo(url)

	if p != nil && err != nil {
		println("Success")
	}
}
