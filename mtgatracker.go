package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"strconv"
)

func GetTransactionStrings(mtgaLog string) []string {
	s := strings.Split(mtgaLog, "\n")
	all_transaction_strings := []string{}
	current_str_arr := []string{}
	processing_json := false
	for _, line := range s {
		if strings.HasPrefix(line, "{") {
			processing_json = true
		} 
		if processing_json {
			current_str_arr = append(current_str_arr, line)
		}
		if strings.HasPrefix(line, "}") {
			processing_json = false
			all_transaction_strings = append(all_transaction_strings, strings.Join(current_str_arr, "\n"))
			current_str_arr = []string{}
		}
	}
	return all_transaction_strings
}

func GetTransactions(mtgaLog string) []Transaction {
	transaction_strings := GetTransactionStrings(mtgaLog)
	transactions := make([]Transaction, 0, len(transaction_strings))
	for line, transaction_string := range transaction_strings {
		var t Transaction
		err := json.Unmarshal([]byte(transaction_string), &t)
		if err != nil {
			log.Print(fmt.Sprintf("Transaction number: %d", line))
			log.Print(fmt.Sprintf("Transaction contents: %s", transaction_string))
			log.Fatal(err)
		}
		transactions = append(transactions, t)
	}
	return transactions
}

/*
 * Returns a map of player id to map of turns to array of card ids.
 */
func GetHands(transactions []Transaction) map[int]map[int][]int {
	hand := map[int]map[int][]int{}
	for _, transaction := range transactions {
		greMessages := transaction.GreToClientEvent.GreToClientMessages
		for _, greMessage := range greMessages {
			turnInfo := greMessage.GameStateMessage.TurnInfo
			turn := -1
			if turnInfo != nil {
				turn = turnInfo.TurnNumber
			}
			zones := greMessage.GameStateMessage.Zones
			if zones != nil {
				for _, zone := range *zones {
					if zone.Type == "ZoneType_Hand" {
						if hand[zone.OwnerSeatId] == nil {
							hand[zone.OwnerSeatId] = map[int][]int{}
						}
						hand[zone.OwnerSeatId][turn] = zone.ObjectInstanceIds
					}
				}
			}
		}
	}
	return hand
}

func GetCardData() map[int]Card {
	cardDataString, err := ioutil.ReadFile("CardData.json")
	if err != nil {
		log.Fatal(err)
	}
	tmpCardData := map[string]Card{}
	err = json.Unmarshal(cardDataString, &tmpCardData)
	if err != nil {
		log.Fatal(err)
	}
	cardData := map[int]Card{}
	for id, _ := range tmpCardData {
		int_id, err := strconv.Atoi(id)
		if err != nil {
			log.Fatal(err)
		}
		cardData[int_id] = tmpCardData[id]
	}
	return cardData
} 

/*
 * Returns a mapping of the game object id of a card to the id of a card in the master
 * MTGA database.
 */
func GetGameObjectCards(transactions []Transaction) map[int]int {
	cardMap := map[int]int{}
	for _, transaction := range transactions {
		greMessages := transaction.GreToClientEvent.GreToClientMessages
		for _, greMessage := range greMessages {
			gameObjects := greMessage.GameStateMessage.GameObjects
			if gameObjects != nil {
				for _, gameObject := range *gameObjects {
					if gameObject.Type == "GameObjectType_Card" {
						cardMap[gameObject.InstanceId] = gameObject.GrpId
					}
				} 
			}
		}
	}
	return cardMap
}

/*
 * Get the winner's seat Id
 */
func GetWinner(transactions []Transaction) int {
	for _, transaction := range transactions {
		greMessages := transaction.GreToClientEvent.GreToClientMessages
		for _, greMessage := range greMessages {
			gameInfo := greMessage.GameStateMessage.GameInfo
			if gameInfo != nil {
				results := gameInfo.Results
				for _, result := range results {
					if result.Result == "ResultType_WinLoss" {
						return result.WinningTeamId
					}
				}
			}
		}
	}
	log.Fatal("No winner found")
	return 0
}

/*
 * Returns a map of seatId to player names
 */
func GetPlayerSeats(transactions []Transaction) map[int]string {
	seatMap := map[int]string{}
	for _, transaction := range transactions {
		gameRoom := transaction.MatchGameRoomStateChangedEvent
		if gameRoom != nil {
			players := gameRoom.GameRoomInfo.GameRoomConfig.ReservedPlayers
			for _, player := range players {
				seatMap[player.SystemSeatId] = player.PlayerName
			} 
		}
	}
	return seatMap
}

func handToCardNames(hands map[int][]int, cardData map[int]Card, gameObjectCards map[int]int) []string {
	handCardNames := make([]string, 0, len(hands))
	for _, playerHands := range hands {
		for _, cardId := range playerHands {
			gameObjectCardId := gameObjectCards[cardId]
			card, ok := cardData[gameObjectCardId]
			if !ok {
				log.Fatal("Card with game ID " + strconv.Itoa(cardId) + " and MTGA ID " + strconv.Itoa(gameObjectCardId) + " not found in card data.")
			}
			handCardNames = append(handCardNames, card.Name)
		}
	}
	return handCardNames
}

func printWinner(userName string, playerSeats map[int]string, winner int, handCardNames []string) {
	wonOrLost := ""
	if playerSeats[winner] == userName {
		wonOrLost = "won"
	} else {
		wonOrLost = "lost"
	}
	fmt.Println("You " + wonOrLost + " with cards \"" + strings.Join(handCardNames, "\", \"") + "\".")
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	userName := "Jolteon#55824"

	//mtgaLog, err := ioutil.ReadFile("C:\\Users\\joel\\AppData\\LocalLow\\Wizards Of The Coast\\MTGA\\output_log.txt")
	mtgaLog, err := ioutil.ReadFile("C:\\Users\\joel\\mtgaoutput\\output_log.txt")
	if err != nil {
		log.Fatal(err)
	}
	cardData := GetCardData()
	transactions := GetTransactions(string(mtgaLog))
	hands := GetHands(transactions)
	gameObjectCards := GetGameObjectCards(transactions)
	playerSeats := GetPlayerSeats(transactions)
	winner := GetWinner(transactions)
	
	userNamePlayerSeat := 0
	for i, playerName := range playerSeats {
		if playerName == userName {
			userNamePlayerSeat = i
		}
	}
	if userNamePlayerSeat == 0 {
		log.Fatal("Could not find userName's player ID")
	}

	handCardNames := handToCardNames(hands[userNamePlayerSeat], cardData, gameObjectCards)
	printWinner(userName, playerSeats, winner, handCardNames)

}
