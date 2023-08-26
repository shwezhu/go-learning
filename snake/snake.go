package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"math/rand"
	"time"
)

const (
	width = 20
	height = 20
)

const (
	right = 0
	left  = 1
	up    = 2
	down  = 3
)

type node struct {
	x int
	y int
}

var (
	board 		= make([][]bool, height)
	snake 		= []node{{x: width / 2, y: height / 2}, {x: width / 2 + 1, y: height / 2}}
	food        = node{1, 3}
	direction   = left
	refreshRate = time.Millisecond * 300
	gameOver    = false
	quitSignal  = make(chan struct{})
)

func generateFood() {
	food.x = rand.Intn(width)
	food.y = rand.Intn(height)
}

func draw() {
	fmt.Print("\033[H\033[2J")
	// original: 'for y, _ := range board { ... }'
	for y := range board {
		for x, xValue := range board[y] {
			if !xValue && (food.x != x || food.y != y) {
				fmt.Print("□ ")
			} else if !xValue && food.x == x && food.y == y {
				fmt.Print("★ ")
			} else {
				fmt.Print("■ ")
			}
		}
		fmt.Println()
	}
}

func update() {
	for _, v := range snake{
		board[v.y][v.x] = true
	}

	// draw map and snake and food
	draw()

	// let snake move, change the position of snake for next calling of draw()
	// by changing the x,y value of snake[0] and snake[len(snake)-1] to make snake "move"
	length := len(snake)
	// ate food, increase length
	if snake[0] == food {
		snake = append(snake, node{})
		generateFood()
	} else {
		// clean the tail on map
		board[snake[length - 1].y][snake[length - 1].x] = false
	}
	// "move" all node of snake except the head
	if length > 1 {
		// Make the value of each (i)th element equal to the (i-1)th value
		for i := len(snake) - 1; i >= 1; i-- {
			snake[i] = snake[i-1]
		}
	}

	// change the head position
	switch direction {
	case right:
		snake[0].x++
	case left:
		snake[0].x--
	case up:
		snake[0].y--
	case down:
		snake[0].y++
	}

	// decide if snake hit the wall
	x := snake[0].x
	y := snake[0].y
	if x > width - 1 || x < 0 || y > height - 1 || y < 0 {
		gameOver = true
	}
}

// init() function will be called once automatically.
func init() {
	// original: 'for i, _ := range board { ... }'
	for i := range board {
		board[i] = make([]bool, width)
	}
}

func monitorInput() {
	for {
		select {
		case <-quitSignal:
			close(quitSignal)
			return
		default:
			// 每次退出游戏后 需要再按一次按键才能退出就是因为这个GetKey
			// 你按下q之后, gameOver确实被改成了true
			// 且在main中执行到了quit <- true
			// 但handleInput是个无限循环, 在‘quit <- true’和main的无限循环退出之间
			// handleInput又开始了下一次循环, 即此时quit还没被写入数据, 所以跳到了下面这段代码
			// keyboard.GetSingleKey(), 所以程序卡住了, 你需要多按一次按键才能游戏结束,
			// 还没想到合适的办法解决
			char, _, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}

			switch char {
			case 'w', 'W':
				if direction != down {
					direction = up
				}
			case 'a', 'A':
				if direction != right {
					direction = left
				}
			case 's', 'S':
				if direction != up {
					direction = down
				}
			case 'd', 'D':
				if direction != left {
					direction = right
				}
			case 'q', 'Q':
				gameOver = true
			}
		}
	}
}

func main() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	// call anonymous function, often used in these sensorial: defer、go
	defer func() {
		err := keyboard.Close()
		if err != nil {
			return
		}
	}()

	go monitorInput()
	for !gameOver {
		update()
		time.Sleep(refreshRate)
	}

	quitSignal <- struct{}{}
	fmt.Printf("Game Over!\n")
}