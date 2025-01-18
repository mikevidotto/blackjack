package main

import (
	"fmt"
	"math/rand/v2"
)

type Hand struct {
	cards []Card
}

type Card struct {
	name    string
	value   int
	cardval int
}

var (
	suits []string       = []string{"hearts", "diamonds", "spades", "clubs"}
	ranks map[int]string = map[int]string{
		1:  "Ace",
		2:  "Two",
		3:  "Three",
		4:  "Four",
		5:  "Five",
		6:  "Six",
		7:  "Seven",
		8:  "Eight",
		9:  "Nine",
		10: "Jack",
		11: "Queen",
		12: "King",
	}
)

func main() {
	valuestore := []int{}
	round(valuestore)
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func round(valuestore []int) {
	//deal a random card to user, then dealer, then user, then dealer.
	//pick random number
	//repeat 3 other times without repeat numbers.
	player, dealer, valuestore := deal(valuestore)
	//reveal the second card dealt to the dealer.

	win := true

	for getHandval(player) < 21 && getHandval(dealer) < 21 {
		displayHands(player, dealer)
		fmt.Println("Hit(1) or Stand(2)?")
		var input int
		fmt.Scanln(&input)

		for input != 1 && input != 2 {
			fmt.Scanln(&input)
		}

		if input == 1 {
			//hit
			fmt.Println("Hit.")
			card, _ := createCard(valuestore)
			player.cards = append(player.cards, card)
		} else {
			fmt.Println("Stand.")
			//stand
			for getHandval(dealer) < 21 {
				card, _ := createCard(valuestore)
				dealer.cards = append(dealer.cards, card)
				if getHandval(dealer) > getHandval(player) || getHandval(dealer) == 21 {
					fmt.Println("Dealer wins.")
					displayHands(player, dealer)
					win = false
					break
				}
			}
			if !win {

			} else {
				fmt.Println("You win.")
				break
			}
		}
	}
	if !win {
		fmt.Println("Dealer wins.")
	}
	//give option to hit or stay.
	//if you go over 21, you lose.
	//if dealer goes over 21 they lose.
	//whoever is closer to 21 wins.
}

func displayHands(player Hand, dealer Hand) {
	if len(dealer.cards) < 3 {
		fmt.Println("Your hand: ")
		for _, card := range player.cards {
			fmt.Println(card.name)
		}
		fmt.Println(getHandval(player))
		fmt.Println("Dealer: ")
		fmt.Println(dealer.cards[1].name)
		fmt.Println(getHandval(dealer))
	} else {
		fmt.Println("Your hand: ")
		for _, card := range player.cards {
			fmt.Println(card.name)
		}
		fmt.Println(getHandval(player))
		fmt.Println("Dealer: ")
		for _, card := range dealer.cards {
			fmt.Println(card.name)
		}
		fmt.Println(getHandval(dealer))
	}
}

func getHandval(hand Hand) (handval int) {
	for _, card := range hand.cards {
		handval += card.value
	}
	return handval
}

func checkDuplicate(card Card, valuestore []int) bool {
	for i := 0; i < len(valuestore); i++ {
		if card.cardval == valuestore[i] {
			return false
		}
	}
	return true
}

func createCard(valuestore []int) (Card, []int) {
	//GENERATE 4 DIFFERENT CARDS.
	var card Card
	suitval := randRange(1, len(suits))
	rankval := randRange(1, len(ranks))
	card.cardval = rankval + suitval
	card.name = ranks[rankval] + " of " + suits[suitval]
	card.value = rankval
	if checkDuplicate(card, valuestore) == false {
		return createCard(valuestore)
	}

	valuestore = append(valuestore, card.cardval)
	return card, valuestore
}

func deal(valuestore []int) (Hand, Hand, []int) {
	var player Hand
	var dealer Hand
	//create cards that are different.
	card1, valuestore := createCard(valuestore)
	card2, valuestore := createCard(valuestore)
	card3, valuestore := createCard(valuestore)
	card4, valuestore := createCard(valuestore)

	player.cards = append(player.cards, card1)
	dealer.cards = append(dealer.cards, card2)
	player.cards = append(player.cards, card3)
	dealer.cards = append(dealer.cards, card4)

	return player, dealer, valuestore
}

func hit() {

}

func stay() {

}
