package main

import (
	"fmt"
	"sync"
	"time"
)

// inspiration: https://gist.github.com/nmjmdr/d3637b726b564033d318
var wg sync.WaitGroup

var forkIsAvailable = make([]chan bool, 5)
var doneWithFork = make([]chan bool, 5)

func fork(i int) {

	for {
		forkIsAvailable[i] <- true

		select {
		case <-doneWithFork[i]:
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
		isLeftForkFree := false
		isRightForkFree := false

		select {

		case isLeftForkFree = <-forkIsAvailable[leftFork]:
		default:
			break
		}

		select {

		case isRightForkFree = <-forkIsAvailable[rightFork]:
		default:
			break
		}

		if isLeftForkFree && !isRightForkFree {
			doneWithFork[leftFork] <- true

		} else if !isLeftForkFree && isRightForkFree {
			doneWithFork[rightFork] <- true

		} else if isLeftForkFree && isRightForkFree {
			biteCount = biteCount + 1
			fmt.Println("Philosopher", i, ": Eating", biteCount, "/ 3")
			time.Sleep(1000 * time.Millisecond)
			doneWithFork[leftFork] <- true
			doneWithFork[rightFork] <- true
			if biteCount == 3 {
				break
			}
			fmt.Println("Philosopher", i, ": Thinking...")
		}
	}

	fmt.Println("Philosopher", i, ": Finished eating!")
	wg.Done()
}

/*
Why does the program not deadlock?

We prevent a deadlock by making the philosophers drop their forks IF they don't have 2 in their hands.
Furthermore, due to delay between each philosopher picking up a fork, we inherently prevent a deadlock.
They constantly check if they 2 in their hand. If not, they release the fork once again until they can
pick up two.
*/

func main() {
	fmt.Println("Eating commences!")
	wg.Add(5)

	for i := 0; i < 5; i++ {
		forkIsAvailable[i] = make(chan bool)
		doneWithFork[i] = make(chan bool)
	}

	for i := 0; i < 5; i++ {
		go fork(i)
		go philosopher(i)
	}

	wg.Wait()
	fmt.Println("------------\nFINISHED: All philosophers have taken 3 bites.")
}
