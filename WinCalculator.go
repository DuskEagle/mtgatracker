package calculator

import (
	"fmt"
	"ioutil"
	"json"
	"log"
	"math/big"
	"strings"
)

func WinCalculate(gameResults []GameResults) map[String]big.Rat {
	cardResults := map[String]big.Rat
	for _, gameResult := range gameResults {
		for _, card := range gameResult.Hand {
			if gameResult.Result != "won" && gameResult.Result != "lost" {
				log.Fatal(fmt.Sprintf("Unrecognized game result %s" gameResult.Result))
			} 
			mapLookup, ok := cardResults[card]
			if !ok {
				cardResults[card] = big.Rat(0,1)
			} else {
				*(cardResults[card].Denum) = *(cardResults[card].Denum) + 1
			}
			if gameResult.Result == "win" {
				*(cardResults[card].Num) = *(cardResults[card].Num) + 1
			}
				
		}
	}
}

func main() {
	resultsLog, err := ioutil.ReadFile("output")
	if err != nil {
		log.Fatal(err)
	}
	results = strings.Split(string(resultsLog), "\n")
	gameResults := []GameResult{}
	for _, result := range resultsLog {
		var gr GameResult
		err = json.Unmarshal([]byte(result), &gr)
		if err != nil {
			log.Fatal(err)
		}
		gameResults = append(gameResults, gr)
	}
	fmt.Println(WinCalculate(gameResults))
} 
