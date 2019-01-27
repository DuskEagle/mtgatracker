package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"mtgatracker/parser"
	"strings"
)

func WinCalculate(gameResults []parser.GameResult) map[string]big.Rat {
	cardResults := map[string]big.Rat{}
	for _, gameResult := range gameResults {
		for _, card := range gameResult.Hand {
			if gameResult.Result != "won" && gameResult.Result != "lost" {
				log.Fatal(fmt.Sprintf("Unrecognized game result %s", gameResult.Result))
			} 
			mapLookup, ok := cardResults[card]
			if !ok {
				cardResults[card] = *big.NewRat(0,1)
			} else {
				denom := mapLookup.Denom()
				*denom = *denom.Add(denom, big.NewInt(1))
			}
			if gameResult.Result == "win" {
				num := mapLookup.Num()
				*num = *num.Add(num, big.NewInt(1))
			}
		}
	}
	return cardResults
}

func main() {
	resultsLog, err := ioutil.ReadFile("../parser/output")
	if err != nil {
		log.Fatal(err)
	}
	results := strings.Split(string(resultsLog), "\n")
	gameResults := []parser.GameResult{}
	for _, result := range results {
		var gr parser.GameResult
		err = json.Unmarshal([]byte(result), &gr)
		if err != nil {
			log.Fatal(err)
		}
		gameResults = append(gameResults, gr)
	}
	fmt.Println(WinCalculate(gameResults))
} 
