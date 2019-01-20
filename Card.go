package main

type Card struct {
	Name   string `json:"name"`
	Set    string `json:"set"`
	Images struct {
		Small   string `json:"small"`
		Normal  string `json:"normal"`
		Large   string `json:"large"`
		ArtCrop string `json:"art_crop"`
	} `json:"images"`
	Type        string   `json:"type"`
	Cost        []string `json:"cost"`
	Cmc         int      `json:"cmc"`
	Rarity      string   `json:"rarity"`
	Cid         string   `json:"cid"`
	Frame       []int    `json:"frame"`
	Artist      string   `json:"artist"`
	Dfc         string   `json:"dfc"`
	Collectible bool     `json:"collectible"`
	Craftable   bool     `json:"craftable"`
	DfcID       int      `json:"dfcId"`
	Rank        int      `json:"rank"`
}
