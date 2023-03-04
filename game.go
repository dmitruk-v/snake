package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

type GameConfig struct {
	Width           int
	Height          int
	MaxScore        int
	FruitSpawnDelay time.Duration
	Difficulty      Difficulty
}

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

type Difficulty int

const (
	Easy       Difficulty = 5
	Medium     Difficulty = 7
	Hard       Difficulty = 9
	Impossible Difficulty = 12
)

var DifficultiesMap = map[Difficulty]string{
	Easy:       "EASY",
	Medium:     "MEDIUM",
	Hard:       "HARD",
	Impossible: "IMPOSSIBLE",
}

type Result struct {
	success bool
	spent   time.Duration
}

type Game struct {
	config      GameConfig
	gameRand    *rand.Rand
	isGameOver  bool
	isCompleted bool
	snake       Node
	tail        []Node
	score       int
	fruits      [][]int
	exitCh      chan struct{}
	keyCh       chan rune
	wg          sync.WaitGroup
	ticker      *time.Ticker
}

func NewGame(config GameConfig) *Game {
	fruits := make([][]int, config.Width)
	for i := 0; i < config.Width; i++ {
		fruits[i] = make([]int, config.Height)
	}
	return &Game{
		config:   config,
		gameRand: rand.New(rand.NewSource(time.Now().Unix())),
		fruits:   fruits,
		exitCh:   make(chan struct{}),
		keyCh:    make(chan rune, 1),
		ticker:   time.NewTicker(config.FruitSpawnDelay),
	}
}

func (g *Game) init() {
	// Initiate
	g.placeFruit()
	g.snake = Node{X: g.config.Width / 2, Y: g.config.Height / 2, Dir: Right}
	// Hide cursor
	fmt.Print("\x1b[?25l")
}

func (g *Game) Run() Result {
	// Init and run
	g.init()

	g.wg.Add(1)
	go func() {
		keyEventCh, err := keyboard.GetKeys(0)
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = keyboard.Close()
			g.wg.Done()
		}()
		for {
			select {
			case key := <-keyEventCh:
				g.keyCh <- key.Rune
			case <-g.exitCh:
				close(g.keyCh)
				return
			}
		}
	}()

	// Game loop
	begin := time.Now()
	for {
		// Signal to gorutines to exit
		if g.isGameOver || g.isCompleted {
			close(g.exitCh)
			g.ticker.Stop()
			g.wg.Wait()
		}
		// Game over case
		if g.isGameOver {
			return Result{
				success: false,
				spent:   time.Since(begin),
			}
		}
		// Success game case
		if g.isCompleted {
			return Result{
				success: true,
				spent:   time.Since(begin),
			}
		}
		// g.clearTerm("clear")
		g.gameInput()
		g.gameLogic()
		g.gameDraw()

		time.Sleep(time.Second / time.Duration(g.config.Difficulty))
	}
}

func (g *Game) gameInput() {
	select {
	case key := <-g.keyCh:
		// React to keypress
		switch {
		case key == 'w':
			if g.snake.Dir != Down {
				g.snake.Dir = Up
			}
		case key == 's':
			if g.snake.Dir != Up {
				g.snake.Dir = Down
			}
		case key == 'a':
			if g.snake.Dir != Right {
				g.snake.Dir = Left
			}
		case key == 'd':
			if g.snake.Dir != Left {
				g.snake.Dir = Right
			}
		case key == 0:
			g.gameOver()
			return
		}
	default:
		// skip
	}
}

func (g *Game) gameLogic() {
	// Place fruit on timer
	select {
	case <-g.ticker.C:
		g.placeFruit()
	default:
	}
	// Eat fruit logic
	if g.fruits[g.snake.X][g.snake.Y] != 0 {
		g.score++
		g.fruits[g.snake.X][g.snake.Y] = 0
		if g.score == g.config.MaxScore {
			g.gameCompleted()
			return
		}
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
	if g.doesWallsCollision() || g.doesSelfEat() {
		g.gameOver()
	}
}

func (g *Game) gameDraw() {
	fmt.Print("\x1b[1J\x1b[0;0H")
	for y := 0; y < g.config.Height; y++ {
		for x := 0; x < g.config.Width; x++ {
			// Print top and bottom walls
			if y == 0 || y == g.config.Height-1 {
				fmt.Print("#")
				if x == g.config.Width-1 {
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
			if x == g.config.Width-1 {
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
				fmt.Print("\x1b[38;5;46mF\x1b[0m")
				continue
			}
			// Print empty space
			if !isTail {
				fmt.Print("\x1b[38;5;241m.\x1b[0m")
			}
		}
	}
	diffName := DifficultiesMap[g.config.Difficulty]
	statsTpl := "\x1b[38;5;226mSCORE: %v/%v\x1b[0m, DIFFICULTY: \x1b[38;5;207m%v\x1b[0m\n"
	fmt.Printf(statsTpl, g.score, g.config.MaxScore, diffName)
	if g.isGameOver {
		fmt.Println("\x1b[38;5;198mGAME OVER!\x1b[0m")
	} else if g.isCompleted {
		fmt.Println("\x1b[38;5;51mCOMPLETED!\x1b[0m")
	}
}

func (g *Game) gameOver() {
	g.isGameOver = true
	fmt.Print("\x1b[?25h")
}

func (g *Game) gameCompleted() {
	g.isCompleted = true
	fmt.Print("\x1b[?25h")
}

func (g *Game) placeFruit() {
	x := 1 + g.gameRand.Intn(g.config.Width-2)
	y := 1 + g.gameRand.Intn(g.config.Height-2)
	g.fruits[x][y] = 1
}

func (g *Game) doesWallsCollision() bool {
	return g.snake.X <= 0 || g.snake.X >= g.config.Width-1 || g.snake.Y <= 0 || g.snake.Y >= g.config.Height-1
}

func (g *Game) doesSelfEat() bool {
	for i := 0; i < len(g.tail); i++ {
		if g.snake.X == g.tail[i].X && g.snake.Y == g.tail[i].Y {
			return true
		}
	}
	return false
}
