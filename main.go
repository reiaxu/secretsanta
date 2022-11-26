package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"math/rand"
)

func makeDeck(numCards int) []int {
	cards := make([]int, numCards)
	for i := range cards {
		cards[i] = i
	}

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}

func cutAndShiftCards(cards []int) []int {
	if len(cards) <= 1 {
		return cards
	}

	shiftedCards := make([]int, 0, len(cards))
	shiftedCards = append(shiftedCards, cards[len(cards)-1])
	shiftedCards = append(shiftedCards, cards[0:len(cards)-1]...)

	fmt.Println(shiftedCards)
	return shiftedCards
}

// Exchange is a mapping of the giver to the receiver
type Exchanges map[string]string

func doExchange(names []string) Exchanges {
	exchanges := make(Exchanges)
	deck := makeDeck(len(names))
	shiftedDeck := cutAndShiftCards(deck)
	for i, giverIndex := range deck {
		giverName := names[giverIndex]
		receiverName := names[shiftedDeck[i]]
		exchanges[giverName] = receiverName
	}
	return exchanges
}

type WriteFlusher interface {
	Write(p []byte) (n int, err error)
	Flush() error
}

func showExchanges(r io.Reader, w WriteFlusher, exchanges Exchanges) {
	fmt.Println("Now it's time to reveal the exchanges...")
	fmt.Println("what's your name?")
	fmt.Println("---------------------")

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		giverName := strings.TrimSpace(scanner.Text())
		if len(giverName) == 0 {
			// Empty name entered
			break
		}

		recieverName, giverInNames := exchanges[giverName]
		if !giverInNames {
			fmt.Fprintln(w, "Sorry, your name was not in the list...")
			continue
		}

		fmt.Fprintf(w, "hello, %s \n", giverName)

		buyingMsg := fmt.Sprintf("you're buying a gift for %s", recieverName)
		fmt.Fprintf(w, "%s", buyingMsg)
		if err := w.Flush(); err != nil {
			panic(err)
		}

		time.Sleep(time.Second * 2)
		fmt.Fprintf(w, "\r%s\n", strings.Repeat(" ", len(buyingMsg)+1))
		if err := w.Flush(); err != nil {
			panic(err)
		}
	}
}

func acceptNames(r io.Reader) ([]string, error) {
	fmt.Println("input santa names <3")
	fmt.Println("---------------------")

	scanner := bufio.NewScanner(r)
	names := make(map[string]struct{})
	for scanner.Scan() {
		newName := strings.TrimSpace(scanner.Text())
		if len(newName) == 0 {
			// No more names
			break
		}

		_, nameExists := names[newName]
		if nameExists {
			fmt.Println("This name has already been entered. Try again...")
			continue
		}

		names[newName] = struct{}{}
	}

	if len(names) <= 2 {
		return nil, errors.New("there's too few people")
	}

	nameList := make([]string, 0, len(names))
	for name := range names {
		nameList = append(nameList, name)
	}
	return nameList, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	reader := bufio.NewReader(os.Stdin)
	names, err := acceptNames(reader)
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
		return
	}

	exchanges := doExchange(names)

	writer := bufio.NewWriter(os.Stdout)
	showExchanges(reader, writer, exchanges)
}
