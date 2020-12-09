package parser

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

type dict map[string]interface{}

func init() {
	merchantJSON, err := ioutil.ReadFile("../../../src/app/merchants.json")
	if err != nil {
		panic(fmt.Errorf("BhPhoto: Test Failed %v", err))
	}
	LoadConfig(merchantJSON)
}

func TestBhPhoto(t *testing.T) {
	urls := []dict{
		{
			"url":         "https://www.bhphotovideo.com/c/product/1595098-REG/sony_3005726_hd_camera_for_playstation.html",
			"expectprice": true,
			"instock":     true,
		},
		{
			"url":         "https://www.bhphotovideo.com/c/product/1600080-REG/microsoft_rrt_00001_xbox_series_x_1tb.html",
			"expectprice": false,
			"instock":     true,
		},
		{
			"url":         "https://www.bhphotovideo.com/c/product/1543054-REG/elmo_1430_mx_p2_visual_presenter_and.html",
			"expectprice": true,
			"instock":     true,
		},
		{
			"url":         "https://www.bhphotovideo.com/c/product/1555040-REG/lg_50un7300puf_un7300_50_class_hdr.html",
			"expectprice": true,
			"instock":     true,
		},
		{
			"url":         "https://www.bhphotovideo.com/c/product/1604841-REG/apple_z125_mgn7_05_bh_13_3_macbook_air_with.html",
			"expectprice": true,
			"instock":     true,
		},
	}

	print("Testing Bhphoto...")
	for _, url := range urls {
		p, err := GetProductInfo(url["url"].(string))
		if err != nil {
			t.Errorf("BhPhoto: We got error for Url: %s\n, Error is %v", url, err)
		} else if p == nil {
			t.Errorf("BhPhoto: We got nil object for Url: %s", url)
		}

		if len(p.Price) == 0 && url["expectprice"].(bool) {
			t.Errorf("BhPhoto: we are not getting any price for Url: %s", url)
		}

		if len(p.Available) == 0 && url["instock"].(bool) {
			t.Errorf("BhPhoto: we are not getting any available for Url: %s", url)
		}

		if len(p.Img) == 0 {
			t.Errorf("BhPhoto: we are not getting any image for Url: %s", url)
		}

		if len(p.Name) == 0 {
			t.Errorf("BhPhoto: we are not getting any name for Url: %s", url)
		}
		print(".")
		time.Sleep(2000 * time.Millisecond)
	}
	println("Done")
}

func TestBestbuy(t *testing.T) {
	merchanName := "Bestbuy"
	urls := []dict{
		{
			"url":         "https://www.bestbuy.com/site/apple-airpods-with-charging-case-latest-model-white/6084400.p?skuId=6084400",
			"expectprice": true,
			"instock":     true,
		},
		{
			"url":         "https://www.bestbuy.com/site/sony-playstation-5-console/6426149.p?skuId=6426149",
			"expectprice": true,
			"instock":     false,
		},
		{
			"url":         "https://www.bestbuy.com/site/apple-iphone-12-5g-64gb-black-verizon/6009897.p?skuId=6009897",
			"expectprice": true,
			"instock":     true,
		},
	}

	fmt.Printf("Testing %s...", merchanName)

	for _, url := range urls {
		p, err := GetProductInfo(url["url"].(string))
		if err != nil {
			t.Errorf("%s: We got error for Url: %s\n, Error is %v", merchanName, url, err)
		} else if p == nil {
			t.Errorf("%s: We got nil object for Url: %s", merchanName, url)
		}

		if len(p.Price) == 0 && url["expectprice"].(bool) {
			t.Errorf("%s: we are not getting any price for Url: %s", merchanName, url)
		}

		if len(p.Available) == 0 && url["instock"].(bool) {
			t.Errorf("%s: we are not getting any available for Url: %s", merchanName, url)
		}

		if len(p.Img) == 0 {
			t.Errorf("%s: we are not getting any image for Url: %s", merchanName, url)
		}

		if len(p.Name) == 0 {
			t.Errorf("%s: we are not getting any name for Url: %s", merchanName, url)
		}
		print(".")
		time.Sleep(2000 * time.Millisecond)
	}
	println("Done")
}
