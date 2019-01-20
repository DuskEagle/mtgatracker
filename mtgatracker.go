package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"strconv"
)

func GetTransactionStrings(mtga_log string) []string {
	s := strings.Split(mtga_log, "\n")
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

func GetTransactions(mtga_log string) []Transaction {
	transaction_strings := GetTransactionStrings(mtga_log)
	transactions := make([]Transaction, 0, len(transaction_strings))
	for _, transaction_string := range transaction_strings {
		var t Transaction
		err := json.Unmarshal([]byte(transaction_string), &t)
		if err != nil {
			log.Print(transaction_string)
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
	gre_messages := transaction.GreToClientEvent.GreToClientMessages

	for _, gre_message := range gre_messages {
		turn_info := gre_message.GameStateMessage.TurnInfo
		if turn_info != nil {
			if turn_info.TurnNumber == turn {
				zones := gre_message.GameStateMessage.Zones
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
	card_data_string, err := ioutil.ReadFile("CardData.json")
	if err != nil {
		log.Fatal(err)
	}
	tmp_card_data := map[string]Card{}
	err = json.Unmarshal(card_data_string, &tmp_card_data)
	if err != nil {
		log.Fatal(err)
	}
	card_data := map[int]Card{}
	for id, _ := range tmp_card_data {
		int_id, err := strconv.Atoi(id)
		if err != nil {
			log.Fatal(err)
		}
		card_data[int_id] = tmp_card_data[id]
	}
	return card_data
} 

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//mtga_log, err := os.Open("C:\\Users\\joel\\AppData\\LocalLow\\Wizards Of The Coast\\MTGA\\output_log.txt")
	mtga_log, err := ioutil.ReadFile("C:\\Users\\joel\\mtgaoutput\\output_log.txt")
	if err != nil {
		log.Fatal(err)
	}
	card_data := GetCardData()
	transactions := GetTransactions(string(mtga_log))
	hand := GetHand(transactions[0], 1, 18)
	game_objects := GetGameObjects(transactions)

	hand_card_names := make([]string, 0, len(hand))
	for _, card_id := range hand {
		card, ok := card_data[card_id]
		if !ok {
			log.Fatal("Card with ID " + strconv.Itoa(card_id) + " not found in card data.")
		}
		hand_card_names = append(hand_card_names, card.Name)
	}

	fmt.Println(hand_card_names)
}
