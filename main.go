package main

import (
	"fmt"
	"sync"
	"time"
)

// inspiration: https://gist.github.com/nmjmdr/d3637b726b564033d318
var wg sync.WaitGroup
var free = make([]chan bool, 5)
var done = make([]chan bool, 5)

func fork(i int) {
	for {
		free[i] <- true

		select {
		case <-done[i]:
			break
		}
	}
}

func phil(i int) {
	var e = 0
	fmt.Println(i, " started eating")

	for {
		left := false
		right := false

		select {
		case left = <-free[i]:
			break
		default:
			left = false
			break
		}

		select {
		case right = <-free[(i+1)%5]:
			break
		default:
			right = false
			break
		}

		if left && !right {
			done[i] <- true
		} else if !left && right {
			done[(i+1)%5] <- true
		} else if left && right {
			fmt.Println(i, " eating", e+1, "/3")
			time.Sleep(1000 * time.Millisecond)
			done[i] <- true
			done[(i+1)%5] <- true
			e = e + 1
			fmt.Println(i, " thinking")
			if e == 3 {
				break
			}
		}
	}

	fmt.Println(i, " finished eating")
	wg.Done()
}

func main() {
	fmt.Println("empty table")
	wg.Add(5)

	for i := 0; i < 5; i++ {
		free[i] = make(chan bool)
		done[i] = make(chan bool)
	}

	for i := 0; i < 5; i++ {
		go fork(i)
		go phil(i)
	}

	wg.Wait()
	fmt.Println("empty table")
}
