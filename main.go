package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

type Stats struct {
	Games []Game
}

type Game struct {
	Playerhand Hand
	Dealerhand Hand
	WonFlag    int
}

type Hand struct {
	Cards []Card
}

type Card struct {
	Name  string
	Suit  string
	Rank  string
	Value int //SPLIT INTO RANK AND SUIT VALUES??? MODIFY createCard()???
}

var (
	suits []string       = []string{"♥", "♦", "♠", "♣"}
	ranks map[int]string = map[int]string{
		1:  "A",
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "10",
		11: "J",
		12: "Q",
		13: "K",
	}
	cardstore []Card
)

func main() {
	play := false
	for !play {
		fmt.Println("\nPLAY (1)  QUIT(2)  SEE STATS(3)")
		cardstore = []Card{}
		var input int
		fmt.Scanln(&input)

		for input != 1 && input != 2 && input != 3 {
			fmt.Scanln(&input)
		}

		if input == 1 {
			round()
		} else if input == 2 {
			play = true
		} else if input == 3 {
			displayStats()
		}
	}
}

func displayStats() {
	stats := readStats()
	wins := 0
	losses := 0
	draws := 0

	for _, game := range stats.Games {
		if game.WonFlag == 1 {
			wins++
		} else if game.WonFlag == 2 {
			losses++
		} else if game.WonFlag == 3 {
			draws++
		}
	}

	fmt.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nS T A T S")
	fmt.Println("Games played:", len(stats.Games))
	fmt.Println("Won:         ", wins)
	fmt.Println("Lost:        ", losses)

	if wins != 0 {
		fmt.Printf("Win %%:        %.2f%%\n", float64(wins)/float64(len(stats.Games))*100)
	} else {
		fmt.Println("Win %:        0%")
	}
}

func readStats() (stats Stats) {
	if _, err := os.Stat("stats.txt"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("No save file.")
	} else {
		file, err := os.ReadFile("stats.txt")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(file, &stats)
		if err != nil {
			panic(err)
		}
	}
	return stats
}

func saveStats(stats Stats) {
	json, err := json.Marshal(stats)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("stats.txt")
	if err != nil {
		panic(err)
	}
	f.Write(json)
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func round() {
	//deal a random card to user, then dealer, then user, then dealer.
	//pick random number
	//repeat 3 other times without repeat numbers.
	game := Game{}
	endgame := false
	standoff := false
	player, dealer := deal()
	//reveal the second card dealt to the dealer.

	for !endgame && !standoff {
		displayHands(player, dealer)
		fmt.Println("Hit(1) or Stand(2)?")
		var input int
		fmt.Scanln(&input)

		for input != 1 && input != 2 {
			fmt.Scanln(&input)
		}

		if input == 1 {
			fmt.Println("hit.")
			card := createCard()
			player.Cards = append(player.Cards, card)
			if one, two := getHandval2(player); one >= 21 || two >= 21 {
				endgame = true
			}
		} else {
			fmt.Println("stand.")
			standoff = true
		}
	}

	if standoff {
		//deal card to dealer unti he gets a higher cardval than player, or until he draws 21 or over.
		game.WonFlag = 0
		for !endgame {
			if one, two := getHandval2(dealer); one < 21 && two < 21 {
				fmt.Println(one, two)
				done, dtwo := getHandval2(dealer)
				pone, ptwo := getHandval2(player)
				if (done > pone && done > ptwo) || (dtwo > pone && dtwo > ptwo) {
					endgame = true
					fmt.Println("1Dealer wins.")
					game.WonFlag = 2
				}
			} else if one, two := getHandval2(dealer); one > 21 && two > 21 {
				fmt.Println(one, two)
				endgame = true
				fmt.Println("You win.")
				game.WonFlag = 1
			} else {
				fmt.Println(one, two)
				if one, two := getHandval2(dealer); one == 21 || two == 21 {
					endgame = true
					if one, two := getHandval2(player); one == 21 || two == 21 {
						fmt.Println("Draw.")
						game.WonFlag = 3
					} else {
						fmt.Println("Dealer wins.")
						game.WonFlag = 2
					}
				}
			}
			if game.WonFlag != 0 {
				displayHands(player, dealer)
				time.Sleep(2 * time.Second)
			} else {
				card := createCard()
				dealer.Cards = append(dealer.Cards, card)
				displayHands(player, dealer)
				time.Sleep(2 * time.Second)
			}
			switch game.WonFlag {
			case 1:
				fmt.Println("You win!")
			case 2:
				fmt.Println("Dealer wins.")
			case 3:
				fmt.Println("Draw.")
			}
		}
	} else {
		if one, two := getHandval2(player); one == 21 || two == 21 {
			fmt.Println(one, two)
			//check if dealer gets 21. otherwise, you win.
			done, dtwo := getHandval2(dealer)
			for done < 21 && dtwo < 21 {
				fmt.Println(done, dtwo)
				card := createCard()
				dealer.Cards = append(dealer.Cards, card)
				done, dtwo = getHandval2(dealer)
				displayHands(player, dealer)
			}
			dealerone, dealertwo := getHandval2(dealer)
			fmt.Println("DEALER:", dealerone, dealertwo)
			if dealerone == 21 || dealertwo == 21 {
				//tie.
				fmt.Println("2Draw.")
				game.WonFlag = 3
			} else {
				displayHands(player, dealer)
				fmt.Println("You win!")
				game.WonFlag = 1
			}
		} else {
			displayHands(player, dealer)
			fmt.Println("3Dealer wins")
			game.WonFlag = 2
		}

	}
	game.Playerhand = player
	game.Dealerhand = dealer
	stats := readStats()
	stats.Games = append(stats.Games, game)
	saveStats(stats)
	//give option to hit or stay.
	//if you go over 21, you lose.
	//if dealer goes over 21 they lose.
	//whoever is closer to 21 wins.
}

func displayHands(player Hand, dealer Hand) {
	if len(dealer.Cards) < 3 {
		fmt.Println("\n\n\n\n\n\n\n\n     D E A L E R ")
		for _ = range dealer.Cards {
			fmt.Print(" ---------  ")
		}
		fmt.Println()
		for i, card := range dealer.Cards {
			if i < 1 {
				fmt.Printf("| %-8s| ", "?")
			} else {
				fmt.Printf("| %-8s| ", card.Rank)
			}
		}
		fmt.Println()
		for _ = range dealer.Cards {
			fmt.Print("|         | ")
		}
		fmt.Println()
		for i, card := range dealer.Cards {
			if i < 1 {
				fmt.Printf("|    %-2s   | ", "?")
			} else {
				fmt.Printf("|    %-2s   | ", card.Suit)
			}
		}
		fmt.Println()
		for _ = range dealer.Cards {
			fmt.Print("|         | ")
		}
		fmt.Println()
		for i, card := range dealer.Cards {
			if i < 1 {
				fmt.Printf("|       %-2s| ", "?")
			} else {
				fmt.Printf("|       %-2s| ", card.Rank)
			}
		}
		fmt.Println()
		for _ = range dealer.Cards {
			fmt.Print(" ---------  ")
		}
		fmt.Println()
		fmt.Println("     Y O U R  H A N D ")
		for _ = range player.Cards {
			fmt.Print(" ---------  ")
		}
		fmt.Println()
		for _, card := range player.Cards {
			fmt.Printf("| %-8s| ", card.Rank)
		}
		fmt.Println()
		for _ = range player.Cards {
			fmt.Print("|         | ")
		}
		fmt.Println()
		for _, card := range player.Cards {
			fmt.Printf("|    %-2s   | ", card.Suit)
		}
		fmt.Println()
		for _ = range player.Cards {
			fmt.Print("|         | ")
		}
		fmt.Println()
		for _, card := range player.Cards {
			fmt.Printf("|       %-2s| ", card.Rank)
		}
		fmt.Println()
		for _ = range player.Cards {
			fmt.Print(" ---------  ")
		}
		fmt.Println()

	} else {
		fmt.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n     D E A L E R ")
		for _ = range dealer.Cards {
			fmt.Print(" ---------  ")
		}
		fmt.Println()
		for _, card := range dealer.Cards {
			fmt.Printf("| %-8s| ", card.Rank)
		}
		fmt.Println()
		for _ = range dealer.Cards {
			fmt.Print("|         | ")
		}
		fmt.Println()
		for _, card := range dealer.Cards {
			fmt.Printf("|    %-2s   | ", card.Suit)
		}
		fmt.Println()
		for _ = range dealer.Cards {
			fmt.Print("|         | ")
		}
		fmt.Println()
		for _, card := range dealer.Cards {
			fmt.Printf("|       %-2s| ", card.Rank)
		}
		fmt.Println()
		for _ = range dealer.Cards {
			fmt.Print(" ---------  ")
		}
		fmt.Println()
		fmt.Println("     Y O U R  H A N D ")
		for _ = range player.Cards {
			fmt.Print(" ---------  ")
		}
		fmt.Println()
		for _, card := range player.Cards {
			fmt.Printf("| %-8s| ", card.Rank)
		}
		fmt.Println()
		for _ = range player.Cards {
			fmt.Print("|         | ")
		}
		fmt.Println()
		for _, card := range player.Cards {
			fmt.Printf("|    %-2s   | ", card.Suit)
		}
		fmt.Println()
		for _ = range player.Cards {
			fmt.Print("|         | ")
		}
		fmt.Println()
		for _, card := range player.Cards {
			fmt.Printf("|       %-2s| ", card.Rank)
		}
		fmt.Println()
		for _ = range player.Cards {
			fmt.Print(" ---------  ")
		}
		fmt.Println()
	}
}

func getHandval(hand Hand) (handval int) {
	for _, card := range hand.Cards {
		if card.Value >= 10 {
			handval += 10
		} else {
			handval += card.Value
		}
	}
	return handval
}

func getHandval2(hand Hand) (handval int, handval2 int) {
	aces := 0
	for _, card := range hand.Cards {
		if card.Value == 1 {
			//ace
			if aces >= 1 {
				//two aces
				handval += 1
				handval2 += 1
			} else {
				handval += 1
				handval2 += 11
				aces += 1
			}
		} else if card.Value >= 10 {
			handval += 10
			handval2 += 10
		} else {
			handval += card.Value
			handval2 += card.Value
		}
	}
	return handval, handval2
}

func checkDuplicate(card Card) bool {
	for i := 0; i < len(cardstore); i++ {
		if card.Rank == cardstore[i].Rank && card.Suit == cardstore[i].Suit {
			return false
		}
	}
	return true
}

func createCustomCard(suit int, rank int) Card {
	//GENERATE 4 DIFFERENT Cards.
	var card Card
	suitval := suit
	rankval := rank
	card.Rank = ranks[rankval]
	card.Suit = suits[suitval]
	card.Name = ranks[rankval] + " of " + suits[suitval]
	card.Value = rankval
	if !checkDuplicate(card) {
		return createCard()
	}

	cardstore = append(cardstore, card)
	return card
}

func createCard() Card {
	//GENERATE 4 DIFFERENT Cards.
	var card Card
	suitval := randRange(1, len(suits))
	rankval := randRange(1, len(ranks))
	card.Rank = ranks[rankval]
	card.Suit = suits[suitval]
	card.Name = ranks[rankval] + " of " + suits[suitval]
	card.Value = rankval
	if !checkDuplicate(card) {
		return createCard()
	}

	cardstore = append(cardstore, card)
	return card
}

func deal() (Hand, Hand) {
	var player Hand
	var dealer Hand
	//create Cards that are different.
	card1 := createCard()
	card2 := createCard()
	card3 := createCard()
	card4 := createCard()

	player.Cards = append(player.Cards, card1)
	displayHands(player, dealer)
	time.Sleep(2 * time.Second)

	dealer.Cards = append(dealer.Cards, card2)
	displayHands(player, dealer)
	time.Sleep(2 * time.Second)

	player.Cards = append(player.Cards, card3)
	displayHands(player, dealer)
	time.Sleep(2 * time.Second)

	dealer.Cards = append(dealer.Cards, card4)
	displayHands(player, dealer)
	time.Sleep(2 * time.Second)

	return player, dealer

}
