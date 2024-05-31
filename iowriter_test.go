package main

import (
	"testing"
	"time"
)

func TestFileWriter_WriteLine(t *testing.T) {
	ch := make(chan string)
	for i := 0; i < 5; i++ {
		go func() {
			for i := 0; i < 50000; i++ {
				ch <- "test"
			}
		}()
	}
	WriteFromChannel(ch)

	time.Sleep(3 * time.Second)
	close(ch)
	time.Sleep(1 * time.Second)
}
