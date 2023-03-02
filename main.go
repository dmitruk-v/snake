package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

type Node struct {
	X   int
	Y   int
	Dir Direction
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

const (
	width  = 20
	height = 10
)

func main() {

	var (
		score  = 0
		fruits = [width][height]int{}
		snake  Node
		tail   []Node
		exitCh = make(chan struct{})
	)

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	// Initiate
	placeFruit(width-1, height-1, &fruits)
	snake = Node{
		X:   width / 2,
		Y:   height / 2,
		Dir: Right,
	}
	tail = []Node{
		{X: snake.X - 1, Y: snake.Y, Dir: snake.Dir},
		{X: snake.X - 2, Y: snake.Y, Dir: snake.Dir},
		{X: snake.X - 3, Y: snake.Y, Dir: snake.Dir},
	}

	for {

		// Process keypress
		select {
		case event := <-keysEvents:
			log.Println(event)
			switch event.Rune {
			case 'w':
				if snake.Dir != Down {
					snake.Dir = Up
				}
			case 's':
				if snake.Dir != Up {
					snake.Dir = Down
				}
			case 'a':
				if snake.Dir != Right {
					snake.Dir = Left
				}
			case 'd':
				if snake.Dir != Left {
					snake.Dir = Right
				}
			case 0:
				close(exitCh)
				return
			}
		case <-exitCh:
			fmt.Println("Exiting...")
			return
		default:
		}

		// Choose where to move
		switch snake.Dir {
		case Up:
			snake.Y--
		case Down:
			snake.Y++
		case Left:
			snake.X--
		case Right:
			snake.X++
		}

		// Draw logic
		clearTerm("clear")
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Print header and footer
				if y == 0 || y == height-1 {
					fmt.Print("#")
					if x == width-1 {
						fmt.Print("\n")
					}
					continue
				}
				// Print middle part
				if x == 0 {
					fmt.Print("#")
				} else if x == width-1 {
					fmt.Print("#\n")
				} else {
					// Draw snake, fruit or empty space
					isCollide := checkWallsCollision(width, height, snake.X, snake.Y)
					if isCollide {
						close(exitCh)
						fmt.Println("Game over!")
						return
					}

					if fruits[snake.X][snake.Y] != 0 {
						_ = tail
						score++
						fruits[snake.X][snake.Y] = 0
						placeFruit(width, height, &fruits)
					}

					if snake.X == x && snake.Y == y {
						fmt.Print("H")
					} else if fruits[x][y] == 1 {
						fmt.Print("F")
					} else {
						prev := snake
						for i := 0; i < len(tail); i++ {
							if tail[i].X == x && tail[i].Y == y {
								fmt.Print("T")
							}
							tail[i].X = prev.X
							tail[i].Y = prev.Y
							tail[i].Dir = prev.Dir
							prev = tail[i]
						}
						fmt.Print(" ")
					}
					// --------------------------------
				}
			}
		}
		fmt.Printf("SCORE: %v\n", score)
		time.Sleep(time.Second / 5)
	}
}

func clearTerm(clearCmd string) {
	c := exec.Command(clearCmd)
	c.Stdout = os.Stdout
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}

func placeFruit(width, height int, fruits *[width][height]int) {
	x := 1 + rand.Intn(width-2)
	y := 1 + rand.Intn(height-2)
	fruits[x][y] = 1
}

func checkWallsCollision(width, height, x, y int) bool {
	return x <= 0 || x >= width-1 || y <= 0 || y >= height-1
}
