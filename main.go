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

func philosopher(i int) {
	var rightFork = (i + 1) % 5
	var leftFork = i
	var biteCount = 0
	fmt.Println("Philosopher", i, ": started eating...")

	for {
		isLeftFree := false
		isRightFree := false

		select {

		case isLeftFree = <-free[leftFork]:
			break

		default:
			isLeftFree = false
			break
		}

		select {

		case isRightFree = <-free[rightFork]:
			break

		default:
			isRightFree = false
			break
		}

		if isLeftFree && !isRightFree {
			done[leftFork] <- true

		} else if !isLeftFree && isRightFree {
			done[rightFork] <- true

		} else if isLeftFree && isRightFree {
			fmt.Println("Philosopher", i, ": Eating", biteCount+1, "/ 3")
			time.Sleep(1000 * time.Millisecond)
			done[leftFork] <- true
			done[rightFork] <- true
			biteCount = biteCount + 1
			fmt.Println("Philosopher", i, ": Thinking...")
			if biteCount == 3 {
				break
			}
		}
	}

	fmt.Println("Philosopher", i, ": Finished eating!")
	wg.Done()
}

func main() {
	fmt.Println("Eating commences...")
	wg.Add(5)

	for i := 0; i < 5; i++ {
		free[i] = make(chan bool)
		done[i] = make(chan bool)
	}

	for i := 0; i < 5; i++ {
		go fork(i)
		go philosopher(i)
	}

	wg.Wait()
	fmt.Println("------------\nFINISHED: All philosophers have taken 3 bites.")
}
