package main

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
	"os/exec"
	"sync"
)

var wg sync.WaitGroup

func main() {
	err := tb.Init()
	if err != nil {
		panic(err)
	}
	defer tb.Close()

	c := make(chan string)

	wg.Add(1)
	go writer(c)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go speaker(c)
	}

	wg.Wait()
}

// writer reads char from Stdin and then write it to channel
func writer(c chan string) {
	for {
		event := tb.PollEvent()
		switch {
		case event.Key == tb.KeyEsc:
			c <- "bye"
			// notify all reader
			close(c)
			wg.Done()
			return
		case event.Key == tb.KeySpace:
			c <- "space"
		case event.Key == tb.KeyEnter:
			c <- "enter"
		default:
			c <- string(event.Ch)
		}
	}
}

// speaker reads char from channel and then say it
func speaker(c chan string) {
	for w := range c {
		cmd := exec.Command("say", w)
		fmt.Println(w)
		cmd.Output()
	}
	wg.Done()
}