package parser

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
	current_str_arr := []rune{}
	processing_json := false
	brace_count := 0
	for _, line := range s {
		for _, char := range line {
			if char == '{' {
				brace_count += 1
				processing_json = true
			} 
			if processing_json {
				current_str_arr = append(current_str_arr, char)
			}
			if char == '}' {
				brace_count -= 1
				if brace_count == 0 {
					processing_json = false
					all_transaction_strings = append(all_transaction_strings, string(current_str_arr))
					current_str_arr = []rune{}
				}
			}
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
			if err.Error() != "json: cannot unmarshal number into Go struct field Transaction.timestamp of type string" {
				log.Print(fmt.Sprintf("Transaction number: %d", line))
				log.Print(fmt.Sprintf("Transaction contents: %s", transaction_string))
				log.Fatal(err)
			}
		}
		transactions = append(transactions, t)
	}
	return transactions
}

func GetTurn(transaction Transaction) (int, bool) {
	greMessages := transaction.GreToClientEvent.GreToClientMessages
	for _, greMessage := range greMessages {
		turnInfo := greMessage.GameStateMessage.TurnInfo
		if turnInfo != nil {
			return turnInfo.TurnNumber, true
		}
	}
	return 0, false
}

/*
 * Returns a map of player id to map of turns to array of card ids.
 */
func GetHands(transaction Transaction) (map[int][]int, bool) {
	hand := map[int][]int{}
	greMessages := transaction.GreToClientEvent.GreToClientMessages
	for _, greMessage := range greMessages {
		zones := greMessage.GameStateMessage.Zones
		if zones != nil {
			for _, zone := range *zones {
				if zone.Type == "ZoneType_Hand" {
					hand[zone.OwnerSeatId] = zone.ObjectInstanceIds
				}
			}
			return hand, true
		}
	}
	return nil, false
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
func GetGameObjectCards(transaction Transaction) (map[int]int, bool) {
	cardMap := map[int]int{}
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
	if len(cardMap) > 0 {
		return cardMap, true
	}
	return nil, false
}

/*
 * Get the winner's seat Id
 */
func GetWinner(transaction Transaction) (int, bool) {
	greMessages := transaction.GreToClientEvent.GreToClientMessages
	for _, greMessage := range greMessages {
		gameInfo := greMessage.GameStateMessage.GameInfo
		if gameInfo != nil {
			results := gameInfo.Results
			for _, result := range results {
				if result.Result == "ResultType_WinLoss" {
					return result.WinningTeamId, true
				}
			}
		}
	}
	return 0, false
}

/*
 * Returns a map of seatId to player names
 */
func GetPlayerSeats(transaction Transaction) (map[int]string, bool) {
	seatMap := map[int]string{}
	gameRoom := transaction.MatchGameRoomStateChangedEvent
	if gameRoom != nil {
		players := gameRoom.GameRoomInfo.GameRoomConfig.ReservedPlayers
		if players != nil {
			for _, player := range players {
				seatMap[player.SystemSeatId] = player.PlayerName
			}
			return seatMap, true
		}
	}
	return nil, false
}

func HandToCardNames(hand []int, cardData map[int]Card, gameObjectCards map[int]int) []string {
	handCardNames := make([]string, 0, len(hand))
	for _, cardId := range hand {
		gameObjectCardId := gameObjectCards[cardId]
		card, ok := cardData[gameObjectCardId]
		if !ok {
			log.Fatal("Card with game ID " + strconv.Itoa(cardId) + " and MTGA ID " + strconv.Itoa(gameObjectCardId) + " not found in card data.")
		}
		handCardNames = append(handCardNames, card.Name)
	}
	return handCardNames
}

func PrintWinner(userName string, playerSeats map[int]string, winner int, handCardNames []string) {
	wonOrLost := ""
	if playerSeats[winner] == userName {
		wonOrLost = "won"
	} else {
		wonOrLost = "lost"
	}
	fmt.Println("You " + wonOrLost + " with cards \"" + strings.Join(handCardNames, "\", \"") + "\".")
}

func GetGames(transactions []Transaction, cardData map[int]Card, userName string) []GameResult {
	playerSeats := map[int]string{}
	userPlayerSeat := 0
	turn := 0
	hand := map[int]struct{}{}
	gameObjectCards := map[int]int{}

	for i, transaction := range transactions {
		newPlayerSeats, playerSeatsOk := GetPlayerSeats(transaction)
		newTurn, turnOk := GetTurn(transaction)
		newHand, handOk := GetHands(transaction)
		newGameObjectCards, gameObjectCardsOk := GetGameObjectCards(transaction)
		winner, winnerOk := GetWinner(transaction)

		if playerSeatsOk {
			playerSeats = newPlayerSeats
			for i, playerName := range playerSeats {
				if playerName == userName {
					userPlayerSeat = i
				}
			}
			if userPlayerSeat == 0 {
				log.Fatal(fmt.Sprintf("Could not find userName's player ID while handling transaction ID %s", transaction.TransactionId))
			}
		}
		if turnOk {
			turn = newTurn
			_ = turn // Remove declared and not used during development
		}
		if handOk {
			for _, k := range newHand[userPlayerSeat] {
				hand[k] = struct{}{}
			}
		}
		if gameObjectCardsOk {
			for k, v := range newGameObjectCards {
				gameObjectCards[k] = v
			}
		}
		if winnerOk {
			var wonOrLost string
			if playerSeats[winner] == userName {
				wonOrLost = "won"
			} else {
				wonOrLost = "lost"
			}
			handArray := []int{}
			for k, _ := range hand {
				handArray = append(handArray, k)
			}
			handCardNames := HandToCardNames(handArray, cardData, gameObjectCards)
			return append(GetGames(transactions[i+1:], cardData, userName), GameResult{Result: wonOrLost, Hand: handCardNames})
		}
	}
	return []GameResult{}
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	userName := "Jolteon#55824"

	//mtgaLog, err := ioutil.ReadFile("C:\\Users\\joel\\AppData\\LocalLow\\Wizards Of The Coast\\MTGA\\output_log.txt")
	mtgaLog, err := ioutil.ReadFile("C:\\Users\\joel\\output_log.txt")
	if err != nil {
		log.Fatal(err)
	}

	cardData := GetCardData()
	transactions := GetTransactions(string(mtgaLog))

	games := GetGames(transactions, cardData, userName)
	log.Println(len(games))
	for _, game := range games {
		gr, err := json.Marshal(game)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(gr))
	}
}
