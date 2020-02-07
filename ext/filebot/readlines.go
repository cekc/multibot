package filebot

import (
	"bufio"
	"io"
	"log"
)

func readlines(reader io.Reader) <-chan string {
	scanner := bufio.NewScanner(reader)
	lines := make(chan string)

	go func() {
		defer close(lines)

		for scanner.Scan() {
			lines <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			log.Println("IO error:", err)
		}
	}()

	return lines
}
