package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

const delayBetweenGames = 3

func run() error {
	game := NewGame(30, 13, 5, Easy)
	result := game.Run()
	if result.success {
		fmt.Println("EASY completed in", result.spent.Truncate(time.Millisecond))
		fmt.Printf("MEDIUM in %v seconds...\n", delayBetweenGames)
		time.Sleep(time.Second * delayBetweenGames)
	} else {
		return nil
	}

	game = NewGame(35, 15, 10, Medium)
	result = game.Run()
	if result.success {
		fmt.Println("MEDIUM completed in", result.spent.Truncate(time.Millisecond))
		fmt.Printf("HARD in %v seconds...\n", delayBetweenGames)
		time.Sleep(time.Second * delayBetweenGames)
	} else {
		return nil
	}

	game = NewGame(40, 17, 15, Hard)
	result = game.Run()
	if result.success {
		fmt.Println("HARD completed in", result.spent.Truncate(time.Millisecond))
		fmt.Printf("IMPOSSIBLE in %v seconds...\n", delayBetweenGames)
		time.Sleep(time.Second * delayBetweenGames)
	} else {
		return nil
	}

	game = NewGame(45, 19, 30, Impossible)
	result = game.Run()
	if result.success {
		fmt.Println("IMPOSSIBLE completed in", result.spent.Truncate(time.Millisecond))
		fmt.Println("OMG YOU HAVE DONE THE WHOLE GAME! AMAIZING!")
	}
	return nil
}
