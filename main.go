package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Card struct {
	Value string
	Suit  string
}

type Deck []Card

func (d Deck) Draw() (Card, Deck) {
	card, newDeck := d[0], d[1:]
	return card, newDeck
}

func NewDeck(numDecks int) Deck {
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	suits := []string{"Clubs", "Diamonds", "Hearts", "Spades"}
	var deck Deck
	for i := 0; i < numDecks; i++ {
		for _, suit := range suits {
			for _, value := range values {
				deck = append(deck, Card{Value: value, Suit: suit})
			}
		}
	}
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

func CardValue(card Card) int {
	switch card.Value {
	case "A":
		return 11
	case "K", "Q", "J":
		return 10
	default:
		val, _ := strconv.Atoi(card.Value)
		return val
	}
}

func HandValue(hand []Card) int {
	value := 0
	aces := 0
	for _, card := range hand {
		value += CardValue(card)
		if card.Value == "A" {
			aces++
		}
	}
	for value > 21 && aces > 0 {
		value -= 10
		aces--
	}
	return value
}

func PrintHand(hand []Card, hideSecondCard bool) {
	for i, card := range hand {
		if hideSecondCard && i == 1 {
			fmt.Print("[ ], ")
			continue
		}
		fmt.Printf("%s, ", card.Value)
	}
	fmt.Println()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Blackjack!")
	fmt.Print("Enter number of decks (1-6): ")
	decksInput, _ := reader.ReadString('\n')
	decksInput = strings.TrimSpace(decksInput)
	numDecks, err := strconv.Atoi(decksInput)
	if err != nil || numDecks < 1 || numDecks > 6 {
		fmt.Println("Invalid input. Using default of 6 decks.")
		numDecks = 6
	}

	deck := NewDeck(numDecks)
	balance := 20 // Default balance

gameLoop:
	for {
		if len(deck) < 20 {
			deck = NewDeck(numDecks) // Re-shuffle deck if running low
		}

		playerHand, dealerHand := []Card{}, []Card{}
		var card Card

		// Initial deal
		card, deck = deck.Draw()
		playerHand = append(playerHand, card)
		card, deck = deck.Draw()
		dealerHand = append(dealerHand, card)
		card, deck = deck.Draw()
		playerHand = append(playerHand, card)
		card, deck = deck.Draw()
		dealerHand = append(dealerHand, card)

		fmt.Println("Game has started.")
		PrintHand(playerHand, false)
		PrintHand(dealerHand, true)
		fmt.Println("Your move: (h)it, (s)tand")

		for {
			action, _ := reader.ReadString('\n')
			action = strings.TrimSpace(action)

			switch action {
			case "h": // Hit
				card, deck = deck.Draw()
				playerHand = append(playerHand, card)
				fmt.Println("You drew:", card.Value)
				if HandValue(playerHand) > 21 {
					fmt.Println("Bust! Your hand:", HandValue(playerHand))
					balance-- // Assume bet size of 1
					continue gameLoop
				}
			case "s": // Stand
				fmt.Print("Dealer's hand: ")
				PrintHand(dealerHand, false)
				for HandValue(dealerHand) < 17 {
					card, deck = deck.Draw()
					dealerHand = append(dealerHand, card)
					fmt.Println("Dealer draws:", card.Value)
				}
				if HandValue(dealerHand) > 21 {
					fmt.Println("Dealer busts. You win!")
					balance++ // Assume bet size of 1
				} else if HandValue(dealerHand) >= HandValue(playerHand) {
					fmt.Println("Dealer wins. Your loss.")
					balance--
				} else {
					fmt.Println("You win!")
					balance++
				}
				continue gameLoop
			default:
				fmt.Println("Invalid action. Please choose (h)it or (s)tand.")
			}

			fmt.Println("Your hand value:", HandValue(playerHand))
			fmt.Println("Waiting for your move: (h)it, (s)tand")
		}
	}
}
