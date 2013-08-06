// Package tails offers an implementation of tail in Go
package tails

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Tail monitors and returns new lines added to a file by other processes.
// It sends every new line on the passed channel.
func Tail(path string, ch chan string, log bool) {
	currentSize := getFileSize(path)

	for {
		// Only read if the file size has changed
		if getFileSize(path) == currentSize {
			if log {
				fmt.Println("Waiting for file changes...")
			}

			time.Sleep(5000 * time.Millisecond)
			continue
		} else if getFileSize(path) < currentSize {
			currentSize = 0
		}

		fi, _ := os.Open(path)

		// Information
		//pos, _ := fi.Seek(0, os.SEEK_END)
		//pos, _ := fi.Seek(currentSize, 0)

		fi.Seek(currentSize, 0)
		scanner := bufio.NewScanner(fi)

		for scanner.Scan() {
			ch <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error scanning input: ", err)
		}

		fi.Close()
		currentSize = getFileSize(path)
	}
}

// TODO Error handling
func getFileSize(path string) int64 {
	info, _ := os.Stat(path)
	return info.Size()
}