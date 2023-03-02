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

type Game struct {
	width      int
	height     int
	isGameOver bool
	snake      Node
	tail       []Node
	score      int
	fruits     [][]int
	exitCh     chan struct{}
}

func NewGame(width, height int) *Game {
	fruits := make([][]int, width)
	for i := 0; i < width; i++ {
		fruits[i] = make([]int, height)
	}
	return &Game{
		width:  width,
		height: height,
		fruits: fruits,
		exitCh: make(chan struct{}),
	}
}

func (g *Game) Init() {
	// Initiate
	g.placeFruit()
	g.snake = Node{X: g.width / 2, Y: g.height / 2, Dir: Right}
	g.tail = []Node{
		// {X: g.snake.X - 1, Y: g.snake.Y, Dir: g.snake.Dir},
	}
}

func (g *Game) Run() {
	// Init and run
	g.Init()

	// Keyboard init
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	go func() {
		for event := range keysEvents {
			switch event.Rune {
			case 'w':
				if g.snake.Dir != Down {
					g.snake.Dir = Up
				}
			case 's':
				if g.snake.Dir != Up {
					g.snake.Dir = Down
				}
			case 'a':
				if g.snake.Dir != Right {
					g.snake.Dir = Left
				}
			case 'd':
				if g.snake.Dir != Left {
					g.snake.Dir = Right
				}
			case 0:
				close(g.exitCh)
				return
			}
		}
	}()

	// Game loop
	for !g.isGameOver {
		g.clearTerm("clear")
		g.gameLogic()
		g.gameDraw()
		fmt.Printf("SCORE: %v\n", g.score)
		time.Sleep(time.Second / 5)
	}
	fmt.Println("GAME OVER!")
}

func (g *Game) gameLogic() {
	// Eat fruit logic
	if g.fruits[g.snake.X][g.snake.Y] != 0 {
		g.score++
		g.fruits[g.snake.X][g.snake.Y] = 0
		g.placeFruit()
		// Add tail
		last := g.snake
		if len(g.tail) > 0 {
			last = g.tail[len(g.tail)-1]
		}
		node := Node{X: last.X, Y: last.Y, Dir: last.Dir}
		switch last.Dir {
		case Up:
			node.Y = last.Y + 1
		case Down:
			node.Y = last.Y - 1
		case Left:
			node.X = last.X + 1
		case Right:
			node.X = last.X - 1
		}
		g.tail = append(g.tail, node)
	}
	// Move tail logic
	prev := g.snake
	for i := 0; i < len(g.tail); i++ {
		tmp := g.tail[i]
		g.tail[i].X = prev.X
		g.tail[i].Y = prev.Y
		g.tail[i].Dir = prev.Dir
		prev = tmp
	}
	// Move snake head
	switch g.snake.Dir {
	case Up:
		g.snake.Y--
	case Down:
		g.snake.Y++
	case Left:
		g.snake.X--
	case Right:
		g.snake.X++
	}
	// Check game over
	select {
	case <-g.exitCh:
		g.isGameOver = true
	default:
		if g.doesWallsCollision() {
			g.isGameOver = true
		}
		if g.doesSelfEat() {
			g.isGameOver = true
		}
	}
}

func (g *Game) gameDraw() {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			// Print top and bottom walls
			if y == 0 || y == g.height-1 {
				fmt.Print("#")
				if x == g.width-1 {
					fmt.Print("\n")
				}
				continue
			}
			// Print left wall
			if x == 0 {
				fmt.Print("#")
				continue
			}
			// Print right wall
			if x == g.width-1 {
				fmt.Print("#\n")
				continue
			}
			// Print snake head and tail
			if g.snake.X == x && g.snake.Y == y {
				fmt.Print("0")
				continue
			}
			var isTail bool
			for i := 0; i < len(g.tail); i++ {
				if g.tail[i].X == x && g.tail[i].Y == y {
					fmt.Print("O")
					isTail = true
				}
			}
			// Print fruit
			if g.fruits[x][y] == 1 {
				fmt.Print("F")
				continue
			}
			// Print empty space
			if !isTail {
				fmt.Print(" ")
			}
		}
	}
}

func (g *Game) clearTerm(cmdName string) {
	c := exec.Command(cmdName)
	c.Stdout = os.Stdout
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) placeFruit() {
	x := 1 + rand.Intn(g.width-2)
	y := 1 + rand.Intn(g.height-2)
	g.fruits[x][y] = 1
}

func (g *Game) doesWallsCollision() bool {
	return g.snake.X <= 0 || g.snake.X >= g.width-1 || g.snake.Y <= 0 || g.snake.Y >= g.height-1
}

func (g *Game) doesSelfEat() bool {
	for i := 0; i < len(g.tail); i++ {
		if g.snake.X == g.tail[i].X && g.snake.Y == g.tail[i].Y {
			return true
		}
	}
	return false
}
