// package routine
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 1)
	for {
		select {
		case ch <- 0:
		case ch <- 1:
		default:
			fmt.Println("neither 0 nor 1 has been sent")
		}
		i := <-ch
		fmt.Printf("receive: %d\n", i)
		time.Sleep(time.Second / 5)
	}
}
