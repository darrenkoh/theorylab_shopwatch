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

func validate(t *testing.T, urls []dict, merchanName string) {
	for _, url := range urls {
		p, err := GetProductInfoByURL(url["url"].(string))
		if err != nil {
			t.Errorf("%s: We got error for Url: %s\n, Error is %v", merchanName, url["url"], err)
		} else if p == nil {
			t.Errorf("%s: We got nil object for Url: %s", merchanName, url["url"])
		}

		hasPrice := len(p.Price) > 0
		expectPrice := url["expectprice"].(bool)
		if (hasPrice && !expectPrice) ||
			(!hasPrice && expectPrice) {
			t.Errorf("%s: wrong price expectation (hasPrice=%v vs expected=%v) for Url: %s", merchanName, hasPrice, expectPrice, url["url"])
		}

		isAvailable := p.Available == "Yes"
		expectAvailable := url["instock"].(bool)
		if (isAvailable && !expectAvailable) ||
			(!isAvailable && expectAvailable) {
			t.Errorf("%s: wrong in stock expection (isAvailable=%v vs expectAvailable=%v) for Url: %s", merchanName, isAvailable, expectAvailable, url["url"])
		}

		if len(p.Img) == 0 {
			t.Errorf("%s: we are not getting any image for Url: %s", merchanName, url["url"])
		}

		if len(p.Name) == 0 {
			t.Errorf("%s: we are not getting any name for Url: %s", merchanName, url["url"])
		}
		print(".")
		time.Sleep(2000 * time.Millisecond)
	}
	println("Done")
}

func TestBhPhoto(t *testing.T) {
	urls := []dict{
		{
			"url":         "https://www.bhphotovideo.com/c/product/1595098-REG/sony_3005726_hd_camera_for_playstation.html",
			"expectprice": true,
			"instock":     true,
		},
		{
			"url":         "https://www.bhphotovideo.com/c/product/1344710-REG/asus_ph_gt1030_o2g_geforce_gt_1030_2gb.html",
			"expectprice": true,
			"instock":     true,
		},
		{
			"url":         "https://www.bhphotovideo.com/c/product/1543054-REG/elmo_1430_mx_p2_visual_presenter_and.html",
			"expectprice": true,
			"instock":     false,
		},
		{
			"url":         "https://www.bhphotovideo.com/c/product/1555040-REG/lg_50un7300puf_un7300_50_class_hdr.html",
			"expectprice": true,
			"instock":     true,
		},
		{
			"url":         "https://www.bhphotovideo.com/c/product/1604841-REG/apple_z125_mgn7_05_bh_13_3_macbook_air_with.html",
			"expectprice": true,
			"instock":     false,
		},
	}

	print("Testing Bhphoto")
	validate(t, urls, "Bhphoto")
}

func TestBestbuy(t *testing.T) {
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

	print("Testing BestBuy")
	validate(t, urls, "BestBuy")
}

func TestWalmart(t *testing.T) {
	urls := []dict{
		{
			"url":         "https://www.walmart.com/ip/Echelon-Connect-Sport-Indoor-Cycling-Exercise-Bike-with-6-Month-Free-Membership-120-value/533034706",
			"expectprice": true,
			"instock":     true,
		},
		{
			"url":         "https://www.walmart.com/ip/LEGO-Star-Wars-A-wing-Starfighter-75275-Building-Toy-Cool-Gift-Idea-for-Creative-Adults-1-673-Pieces/554462454",
			"expectprice": true,
			"instock":     false,
		},
		{
			"url":         "https://www.walmart.com/ip/seort/360463987",
			"expectprice": true,
			"instock":     false,
		},
		{
			"url":         "https://www.walmart.com/ip/Jackery-SolarSaga-60W-Solar-Panel-Explorer-160-240-500-Portable-Generator-Foldable-Charger-Summer-Camping-Van-RV-Can-t-Charge-440-PowerPro/939780115?wmlspartner=wlpa&selectedSellerId=101019293",
			"expectprice": true,
			"instock":     true,
		},
	}

	print("Testing Walmart")
	validate(t, urls, "Walmart")
}

func TestAmazon(t *testing.T) {
	urls := []dict{
		{
			"url":         "https://www.amazon.com/dp/B07XJ8C8F7",
			"expectprice": true,
			"instock":     false,
		},
		{
			"url":         "https://www.amazon.com/dp/B076PRWVFG",
			"expectprice": true,
			"instock":     true,
		},
		{
			"url":         "https://www.amazon.com/PULSE-3D-Wireless-Headset-PlayStation-5/dp/B08FC6QLKN",
			"expectprice": false,
			"instock":     false,
		},
		{
			"url":         "https://www.amazon.com/Jackery-SolarSaga-Portable-Explorer-Foldable/dp/B07Q71LX84/ref=redir_mobile_desktop?ie=UTF8&aaxitk=TV8Ms-qJfKXkDQYwn-NthQ&hsa_cr_id=3606977080801&ref_=sbx_be_s_sparkle_mcd_asin_1",
			"expectprice": true,
			"instock":     true,
		},
	}

	print("Testing Amazon")
	validate(t, urls, "Amazon")
}

func TestNewegg(t *testing.T) {
	urls := []dict{
		{
			"url":         "https://www.newegg.com/grey-with-blue-diamond-cut-msi-prestige-15-a10sc-296-mainstream/p/N82E16834155448?Item=N82E16834155448",
			"expectprice": false,
			"instock":     false,
		},
		{
			"url":         "https://www.newegg.com/p/N82E16868110292",
			"expectprice": true,
			"instock":     false,
		},
		{
			"url":         "https://www.newegg.com/msi-geforce-rtx-2070-super-rtx-2070-super-gaming-x-trio/p/N82E16814137439",
			"expectprice": true,
			"instock":     true,
		},
	}

	print("Testing Newegg")
	validate(t, urls, "Newegg")
}

func TestAdorama(t *testing.T) {
	urls := []dict{
		{
			"url":         "https://www.adorama.com/nolpminimk2.html",
			"expectprice": true,
			"instock":     true,
		},
		{
			"url":         "https://www.adorama.com/icalpe6nh.html",
			"expectprice": true,
			"instock":     true,
		},
	}

	print("Testing Adorama")
	validate(t, urls, "Adorama")
}

func TestEtsy(t *testing.T) {
	urls := []dict{
		{
			"url":         "https://www.etsy.com/listing/505391661/cat-beginner-diy-felting-kit-wool-felt?ga_order=most_relevant&ga_search_type=all&ga_view_type=gallery&ga_search_query=&ref=sc_gallery-1-4&plkey=8a8b101dcb15ca341ee1b9e9f5fb528119e7b7d6%3A505391661&frs=1&bes=1",
			"expectprice": true,
			"instock":     true,
		},
		{
			"url":         "https://www.etsy.com/listing/504195172/funny-stickers-dont-be-a-prick-cactus?ga_order=most_relevant&ga_search_type=all&ga_view_type=gallery&ga_search_query=&ref=sr_gallery-1-5&frs=1&bes=1",
			"expectprice": true,
			"instock":     true,
		},
	}

	print("Testing Etsy")
	validate(t, urls, "Etsy")
}
