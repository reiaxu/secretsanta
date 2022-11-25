package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

func shuffleCards(n int) []int {
	cards := make([]int, n)
	for i := range cards {
		cards[i] = i + 1
	}

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return cards
}

func cutAndShiftCards(cards []int) []int {
	card := cards[0]
	cards = cards[1:]
	cards = append(cards, card)
	return cards
}

func drawCards(names map[string]([2]int), shift []int, cards []int) map[string]([2]int) {
	assignedIndex := make([]int, 0)

	for index, key := range names {
		// go through index of cards and assign to random key in map which doesnt have value [0, 0]
		if key == [2]int{0, 0} {
			randIndex := -1
			for randIndex == -1 || slices.Contains(assignedIndex, randIndex) {
				randIndex = rand.Intn(len(cards))
			}

			names[index] = [2]int{shift[randIndex], cards[randIndex]}
			assignedIndex = append(assignedIndex, randIndex)
		}
	}

	return names
}

func main() {
	rand.Seed(time.Now().UnixNano())

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("input santa names <3")
	fmt.Println("---------------------")

	names := make(map[string]([2]int))
	for {
		line, err := reader.ReadString('\n')
		line = strings.Replace(line, "\n", "", -1)
		if err != nil {
			log.Fatal(err)
		}
		if len(strings.TrimSpace(line)) == 0 {
			break
		}
		names[line] = [2]int{0, 0}
	}

	if len(names) <= 2 {
		fmt.Println("there's too few of you!")
	} else {
		shuffledCards := shuffleCards(len(names))

		shift := cutAndShiftCards(shuffledCards)

		assigned := drawCards(names, shift, shuffledCards)

		fmt.Println("what's your name?")
		fmt.Println("---------------------")

		for {
			yourName, err := reader.ReadString('\n')
			yourName = strings.Replace(yourName, "\n", "", -1)
			if err != nil {
				log.Fatal(err)
			}
			yourName = strings.TrimSpace(yourName)
			if len(yourName) == 0 {
				break
			}
			fmt.Printf("hello, %s \n", yourName)

			if assigned[yourName] == [2]int{0, 0} {
				fmt.Println("you're not in the santa list!")
			} else {
				// yourIndex := assigned[yourName][0]
				// fmt.Printf("your index is %d \n", yourIndex)
				assigneeIndex := assigned[yourName][1]
				// fmt.Printf("you're buying a gift for %d \n", assigneeIndex)
				for key, value := range assigned {
					if value[0] == assigneeIndex {
						fmt.Printf("you're buying a gift for %s", key)
					}
				}
			}

			time.Sleep(time.Second * 2)
			fmt.Printf("\r                                                                              \n")
		}
	}
}
