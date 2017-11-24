package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var filename string
var random bool
var timeLimit int

func init() {
	flag.StringVar(&filename, "csv", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.BoolVar(&random, "randomize", false, "randomize the question order")
	flag.IntVar(&timeLimit, "limit", 30, "time limit for the quiz, in seconds")
	flag.Parse()
}

func main() {
	file, err := os.Open(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "quiz: error in opening file (%v)\n", err)
		os.Exit(1)
	}

	rdr := csv.NewReader(file)
	rdr.FieldsPerRecord = 2
	records, err := rdr.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "quiz: read error (%v)", err)
	}

	problems := parseCsvLines(records)
	stdin := bufio.NewReader(os.Stdin)

	time := time.NewTimer(time.Duration(timeLimit) * time.Second)

	correct := 0
Outer:
	for i, prob := range problems {
		fmt.Printf("Problem #%d: %v = ", i+1, prob.q)
		ansChan := make(chan string)
		go func() {
			s, readErr := stdin.ReadString('\n')
			if readErr != nil {
				s = "\nEOF\n"
			} else {
				s = strings.TrimSpace(s)
			}
			ansChan <- s
		}()
		select {
		case <-time.C:
			println()
			break Outer
		case ans := <-ansChan:
			if ans == prob.a {
				correct++
			} else if ans == "\nEOF\n" {
				println()
			}
		}
	}
	fmt.Printf("You got %d correct answers out of %d questions\n", correct, len(records))
}

type problem struct {
	q, a string
}

func parseCsvLines(records [][]string) []problem {
	ret := make([]problem, len(records))
	perms := make([]int, len(records))
	if random {
		perms = rand.Perm(len(records))
	} else {
		for i, _ := range perms {
			perms[i] = i
		}
	}
	for i, rec := range records {
		ret[perms[i]] = problem{
			q: rec[0],
			a: strings.TrimSpace(rec[1]),
		}
	}
	return ret
}
