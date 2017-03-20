package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

type ShopstyleV2ApiResponse struct {
	Products []struct {
		Brand struct {
			Name string `json:"name"`
		} `json:"brand"`
		ClickURL    string `json:"clickUrl"`
		Currency    string `json:"currency"`
		Description string `json:"description"`
		ExtractDate string `json:"extractDate"`
		Image       struct {
			Sizes struct {
				Best struct {
					ActualHeight int    `json:"actualHeight"`
					ActualWidth  int    `json:"actualWidth"`
					Height       int    `json:"height"`
					SizeName     string `json:"sizeName"`
					URL          string `json:"url"`
					Width        int    `json:"width"`
				} `json:"Best"`
				Original struct {
					ActualHeight int    `json:"actualHeight"`
					ActualWidth  int    `json:"actualWidth"`
					SizeName     string `json:"sizeName"`
					URL          string `json:"url"`
				} `json:"Original"`
			} `json:"sizes"`
		} `json:"image"`
		InStock    bool    `json:"inStock"`
		Locale     string  `json:"locale"`
		Name       string  `json:"name"`
		Price      float32 `json:"price"`
		PriceLabel string  `json:"priceLabel"`
		Retailer   struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Score int    `json:"score"`
		} `json:"retailer"`
		SeeMoreLabel string `json:"seeMoreLabel"`
		Sizes        []struct {
			Name string `json:"name"`
		} `json:"sizes"`
		UnbrandedName string `json:"unbrandedName"`
	} `json:"products"`
}

func callShopStyle() (results *ShopstyleV2ApiResponse, err error) {
	url := fmt.Sprintf("http://api.shopstyle.com/api/v2/products?pid=%s&fts=sunglasses&offset=%d&limit=5", os.Getenv("SHOPSTYLE_API_TOKEN"), 0)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")

	log.Info(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Info("Shopstyle response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	log.Info(results)

	return results, nil
}
