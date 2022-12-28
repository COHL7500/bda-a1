package main

import (
	"fmt"
	"sync"
	"time"
)

// inspiration: https://gist.github.com/nmjmdr/d3637b726b564033d318
var wg sync.WaitGroup

// 5 forks

// Channel for checking whether a fork is available.
var forkIsAvailable = make([]chan bool, 5)

// Channel for checking whether a fork is done.
var doneWithFork = make([]chan bool, 5)

// Function simulating a particular "i" fork.
func fork(i int) {

	// While fork is not done, set it to available...
	for {
		forkIsAvailable[i] <- true

		// When fork is done is set to true, break the loop.
		select {
		case <-doneWithFork[i]:
			break
		}
	}
}

// Function simulating a particular "i" philosopher.
func philosopher(i int) {

	// Right fork is always one number higher than current i, thus i + 1.
	// We find modulo of 5, since there's no more than 5 forks, thus it will "reset" to 0.
	var rightFork = (i + 1) % 5

	// Left fork is always current i, thus i.
	// Left fork (i) also corresponds to the philosopher's number (i).
	var leftFork = i

	var biteCount = 0

	fmt.Printf("Philosopher %v: started eating...", i)

	for {
		// Initiate eating by setting left and right fork availability to false.
		leftForkFree := false
		rightForkFree := false

		// Set fork to available/unavailable depending on whatever info channel has

		select {
		case leftForkFree = <-forkIsAvailable[leftFork]:
		default:
			break
		}

		select {
		case rightForkFree = <-forkIsAvailable[rightFork]:
		default:
			break
		}

		// whichever fork is free, use and get done with that fork.
		if leftForkFree && !rightForkFree {
			doneWithFork[leftFork] <- true

		} else if !leftForkFree && rightForkFree {
			doneWithFork[rightFork] <- true

			// if both are free, commence the eating...
		} else if leftForkFree && rightForkFree {

			biteCount = biteCount + 1

			fmt.Printf("Philosopher %v: Eating %v / 3", i, biteCount)

			// Spend 1 second eating.
			time.Sleep(1000 * time.Millisecond)

			// After 1 second, inform that you're done eating.
			doneWithFork[leftFork] <- true
			doneWithFork[rightFork] <- true

			// Whenever 3 bites have been taking, the philsopher is done eating!
			if biteCount == 3 {
				break
			}

			fmt.Printf("Philosopher %v: Thinking...", i)
		}
	}

	fmt.Printf("Philosopher %v: Finished eating!", i)

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

	// When all in the wait group has announced they're finished (aka the philosophers), announce they're finished.
	wg.Wait()
	fmt.Println("------------\nFINISHED: All philosophers have taken 3 bites.")
}
