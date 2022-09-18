package main

import (
	"fmt"
	"sync"
	"time"
)

// inspiration: https://gist.github.com/nmjmdr/d3637b726b564033d318
var wg sync.WaitGroup

var forkFree = make([]chan bool, 5)
var forkDone = make([]chan bool, 5)

func fork(i int) {
	for {
		forkFree[i] <- true

		select {
		case <-forkDone[i]:
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

		case isLeftFree = <-forkFree[leftFork]:
			break

		default:
			isLeftFree = false
			break
		}

		select {

		case isRightFree = <-forkFree[rightFork]:
			break

		default:
			isRightFree = false
			break
		}

		if isLeftFree && !isRightFree {
			forkDone[leftFork] <- true

		} else if !isLeftFree && isRightFree {
			forkDone[rightFork] <- true

		} else if isLeftFree && isRightFree {
			biteCount = biteCount + 1
			fmt.Println("Philosopher", i, ": Eating", biteCount, "/ 3")
			time.Sleep(1000 * time.Millisecond)
			forkDone[leftFork] <- true
			forkDone[rightFork] <- true
			if biteCount == 3 {
				break
			}
			fmt.Println("Philosopher", i, ": Thinking...")
		}
	}

	fmt.Println("Philosopher", i, ": Finished eating!")
	wg.Done()
}

func main() {
	fmt.Println("Eating commences!")
	wg.Add(5)

	for i := 0; i < 5; i++ {
		forkFree[i] = make(chan bool)
		forkDone[i] = make(chan bool)
	}

	for i := 0; i < 5; i++ {
		go fork(i)
		go philosopher(i)
	}

	wg.Wait()
	fmt.Println("------------\nFINISHED: All philosophers have taken 3 bites.")
}
