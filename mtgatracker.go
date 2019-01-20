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
 * Returns a map of (player, turn) to array of card ids.
 */
func GetHand(transaction Transaction, player int, turn int) []int {
	greMessages := transaction.GreToClientEvent.GreToClientMessages

	for _, greMessage := range greMessages {
		turnInfo := greMessage.GameStateMessage.TurnInfo
		if turnInfo != nil {
			if turnInfo.TurnNumber == turn {
				zones := greMessage.GameStateMessage.Zones
				if zones != nil {
					for _, zone := range *zones {
						if zone.Type == "ZoneType_Hand" && zone.OwnerSeatId == player {
							return zone.ObjectInstanceIds
						}
					}
				}
			}
		}
	}
	log.Print("No hand found")
	return nil
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
	hand := GetHand(transactions[1], 1, 18)
	gameObjectCards := GetGameObjectCards(transactions)
	playerSeats := GetPlayerSeats(transactions)
	winner := GetWinner(transactions)
	
	handCardNames := make([]string, 0, len(hand))
	for _, cardId := range hand {
		gameObjectCardId := gameObjectCards[cardId]
		card, ok := cardData[gameObjectCardId]
		if !ok {
			log.Fatal("Card with game ID " + strconv.Itoa(cardId) + " and MTGA ID " + strconv.Itoa(gameObjectCardId) + " not found in card data.")
		}
		handCardNames = append(handCardNames, card.Name)
	}

	wonOrLost := ""
	if playerSeats[winner] == userName {
		wonOrLost = "won"
	} else {
		wonOrLost = "lost"
	}
	fmt.Println("You " + wonOrLost + " with cards \"" + strings.Join(handCardNames, "\", \"") + "\".")

}
