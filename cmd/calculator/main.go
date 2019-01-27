package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"mtgatracker/parser"
	"sort"
	"strings"
)

type CardResults map[string]*big.Rat

func WinCalculate(gameResults []parser.GameResult) CardResults {
	cardResults := CardResults{}
	for _, gameResult := range gameResults {
		for _, card := range gameResult.Hand {
			if gameResult.Result != "won" && gameResult.Result != "lost" {
				log.Fatal(fmt.Sprintf("Unrecognized game result %s", gameResult.Result))
			} 
			mapLookup, ok := cardResults[card]
			if !ok {
				mapLookup = big.NewRat(0,1)
				cardResults[card] = mapLookup
			} else {
				denom := mapLookup.Denom()
				*denom = *denom.Add(denom, big.NewInt(1))
			}
			if gameResult.Result == "won" {
				num := mapLookup.Num()
				*num = *num.Add(num, big.NewInt(1))
			}
		}
	}
	return cardResults
}

func PrettyPrint(cardResults CardResults) string {
	type kv struct {
		Key   string
		Value *big.Rat
	}
	crs := []kv{}
	for k, v := range cardResults {
		crs = append(crs, kv{k, v})
	}
	sort.Slice(crs, func(i, j int) bool {
		return crs[i].Value.Cmp(crs[j].Value) > 0
	})
	resultArr := []string{}
	for _, kv := range crs {
		fl, _ := kv.Value.Float32()
		resultArr = append(resultArr, fmt.Sprintf("%-30s %.0f%%", fmt.Sprintf("\"%s\":", kv.Key), fl*100))
	}
	return strings.Join(resultArr, "\n")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	resultsLog, err := ioutil.ReadFile("output")
	if err != nil {
		log.Fatal(err)
	}
	results := strings.Split(string(resultsLog), "\n")
	gameResults := []parser.GameResult{}
	for _, result := range results {
		if len(result) > 0 {
			var gr parser.GameResult
			err = json.Unmarshal([]byte(result), &gr)
			if err != nil {
				log.Fatal(err)
			}
			gameResults = append(gameResults, gr)
		}
	}
	cardResults := WinCalculate(gameResults)
	fmt.Println(PrettyPrint(cardResults))
} 
