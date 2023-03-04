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

const delayBetweenGames = 3 * time.Second

type Level struct {
	width           int
	height          int
	maxScore        int
	diff            Difficulty
	fruitSpawnDelay time.Duration
}

var levels = []Level{
	{width: 30, height: 13, maxScore: 5, diff: Easy, fruitSpawnDelay: 5 * time.Second},
	{width: 35, height: 15, maxScore: 10, diff: Medium, fruitSpawnDelay: 5 * time.Second},
	{width: 40, height: 17, maxScore: 15, diff: Hard, fruitSpawnDelay: 10 * time.Second},
	{width: 45, height: 19, maxScore: 20, diff: Impossible, fruitSpawnDelay: 10 * time.Second},
}

func run() error {

	for i, level := range levels {
		cfg := GameConfig{
			Width:           level.width,
			Height:          level.height,
			MaxScore:        level.maxScore,
			Difficulty:      level.diff,
			FruitSpawnDelay: level.fruitSpawnDelay,
		}
		game := NewGame(cfg)
		result := game.Run()
		if result.success {
			fmt.Printf("%v completed in %v\n", DifficultiesMap[level.diff], result.spent.Truncate(time.Millisecond))
			if i < len(levels)-1 {
				nextDiff := levels[i+1].diff
				fmt.Printf("%v in %v seconds...\n", DifficultiesMap[nextDiff], delayBetweenGames)
			}
			time.Sleep(delayBetweenGames)
		} else {
			return nil
		}
	}
	return nil
}
