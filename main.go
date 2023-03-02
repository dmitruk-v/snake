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

type Snake struct {
	HX  int
	HY  int
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
	width  = 50
	height = 20
)

func main() {

	var (
		fruits = [width][height]int{}
		snake  = Snake{
			HX:  width / 2,
			HY:  height / 2,
			Dir: Right,
		}
		exitCh = make(chan struct{})
	)

	keysEvents, err := keyboard.GetKeys(1)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {

		// Process keypress
		select {
		case event := <-keysEvents:
			switch event.Key {
			case keyboard.KeyArrowUp:
				if snake.Dir != Down {
					snake.Dir = Up
				}
			case keyboard.KeyArrowDown:
				if snake.Dir != Up {
					snake.Dir = Down
				}
			case keyboard.KeyArrowLeft:
				if snake.Dir != Right {
					snake.Dir = Left
				}
			case keyboard.KeyArrowRight:
				if snake.Dir != Left {
					snake.Dir = Right
				}
			case keyboard.KeyEsc:
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
			snake.HY--
		case Down:
			snake.HY++
		case Left:
			snake.HX--
		case Right:
			snake.HX++
		}

		// Draw logic
		clearTerm("clear")
		for y := 0; y < height; y++ {
			for x := 0; y < width; x++ {
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
					fmt.Print("\n")
				} else {
					// Draw snake, fruit or empty space
					if x == snake.HX && y == snake.HY {
						fmt.Print("S")
					} else {
						if fruits[x][y] == 1 {
							fmt.Print("F")
						} else {
							fmt.Print(" ")
						}
					}
					// --------------------------------
				}
			}
		}

		placeFruit(width, height, fruits)
		time.Sleep(time.Second / 10)
	}
}

func clearTerm(clearCmd string) {
	c := exec.Command(clearCmd)
	c.Stdout = os.Stdout
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}

func placeFruit(width, height int, fruits [width][height]int) {
	rx := 1 + rand.Intn(width-1)
	ry := 1 + rand.Intn(height-1)
	fruits[ry][rx] = 1
}
