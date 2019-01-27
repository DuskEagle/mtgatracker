package main

import (
	"fmt"
	"sort"
	"strings"
)

type GameResult struct {
	Result string
	Hand []string
}

func (gr GameResult) String() string {
	var resultString string
	if gr.Result == "won" {
		resultString = "\"Won\",  "
	} else if gr.Result == "lost" {
		resultString = "\"Lost\", "
	} else {
		resultString = "\"???\",  "
	}
	handArrStr := []string{}
	sortedHand := make([]string, len(gr.Hand))
	copy(sortedHand, gr.Hand)
	sort.Strings(sortedHand)
	for i, card := range sortedHand {
		cardString := fmt.Sprintf("\"%s\"", card)
		if i == 0 {
			handArrStr = append(handArrStr, cardString)
		} else {
			handArrStr = append(handArrStr, fmt.Sprintf("        %s", cardString))
		}
	}
	return fmt.Sprintf("%s%s", resultString, strings.Join(handArrStr, "\n"))
}
