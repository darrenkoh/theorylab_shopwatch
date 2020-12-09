package parser

import (
	"io/ioutil"
	"testing"
	"time"
)

type dict map[string]interface{}

func TestBhPhoto(t *testing.T) {
	merchantJSON, err := ioutil.ReadFile("../../../src/app/merchants.json")
	if err != nil {
		t.Errorf("BhPhoto: Test Failed %v", err)
	}
	loadConfig(merchantJSON)
	urls := []dict{
		{
			"url":         "https://www.bhphotovideo.com/c/product/1595098-REG/sony_3005726_hd_camera_for_playstation.html",
			"expectprice": true,
		},
		{
			"url":         "https://www.bhphotovideo.com/c/product/1600080-REG/microsoft_rrt_00001_xbox_series_x_1tb.html",
			"expectprice": false,
		},
		{
			"url":         "https://www.bhphotovideo.com/c/product/1543054-REG/elmo_1430_mx_p2_visual_presenter_and.html",
			"expectprice": true,
		},
		{
			"url":         "https://www.bhphotovideo.com/c/product/1555040-REG/lg_50un7300puf_un7300_50_class_hdr.html",
			"expectprice": true,
		},
		{
			"url":         "https://www.bhphotovideo.com/c/product/1604841-REG/apple_z125_mgn7_05_bh_13_3_macbook_air_with.html",
			"expectprice": true,
		},
	}

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

		if len(p.Available) == 0 {
			t.Errorf("BhPhoto: we are not getting any available for Url: %s", url)
		}

		if len(p.Img) == 0 {
			t.Errorf("BhPhoto: we are not getting any image for Url: %s", url)
		}

		if len(p.Name) == 0 {
			t.Errorf("BhPhoto: we are not getting any name for Url: %s", url)
		}
		println("Sleeping for 2 seconds...")
		time.Sleep(2000 * time.Millisecond)
		println("Continue next test...")
	}
	println("Done")
}
