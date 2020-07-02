package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var banks []bank
var self string

func init() {
	banks = append(banks, &nbcb{id: "1005"})
	banks = append(banks, &nbcb{id: "1004"})
}

func formatter(banks []bank) string {
	var output []string
	for _, b := range banks {
		output = append(output, b.formatter())
	}
	return strings.Join(output, " ")
}

func main() {
	self, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get self path: %v", err)
	}
	f, err := os.OpenFile(filepath.Join(filepath.Dir(self), "monitor.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(io.MultiWriter(f))

	for _, b := range banks {
		b.update()
	}

	for _, b := range banks {
		go func(b bank) {
			for {
				b.update()
				time.Sleep(1 * time.Second)
			}
		}(b)
	}

	go func() {
		for {
			s := fmt.Sprintf("%s\r", formatter(banks))
			log.Println(strings.ReplaceAll(strings.ReplaceAll(s, "↑", " "), "↓", " "))
			time.Sleep(15 * time.Second)
		}
	}()

	for {
		fmt.Printf("%-20s %s\r", time.Now().Format("2006-01-02 15:04:05"), formatter(banks))
		time.Sleep(100 * time.Millisecond)
	}
}
