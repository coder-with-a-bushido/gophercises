package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	csv := flag.String("csv", "problems.csv", "a csv filename used for problems")
	limit := flag.Int("limit", 30, "time limit for quiz in seconds")
	flag.Parse()

	problemsMap := parseCsv(*csv)
	count, correctCount := 0, 0
	isGameOver := false
	pass := make(chan bool)

	var someKey string
	fmt.Println("Press Any Key to Start the Quiz!")
	fmt.Scanln(&someKey)
	ticker := time.NewTicker(time.Duration(*limit) * time.Second)

	for qn, ans := range *problemsMap {
		count++
		func() {
			go askQues(qn, ans, &count, pass)
			select {
			case isCorrect := <-pass:
				if isCorrect {
					correctCount++
				}
				ticker.Reset(time.Duration(*limit) * time.Second)
				return
			case <-ticker.C:
				fmt.Println("\nTime Out!")
			}
			isGameOver = true
			ticker.Stop()
		}()
		if isGameOver {
			break
		}

	}
	fmt.Printf("You have answered %d out of %d questions correctly!", correctCount, count)
}

func askQues(qn string, ans string, count *int, pass chan bool) {
	fmt.Printf("Problem #%d: %s = ", *count, qn)
	var userAns string
	fmt.Scanln(&userAns)
	if userAns == ans {
		pass <- true
	} else {
		pass <- false
	}
}

func parseCsv(path string) *map[string]string {
	// open csv file
	f, f_err := os.Open(path)
	//if file err
	if f_err != nil {
		log.Fatal(f_err)
	}
	// close file at end of fn
	defer f.Close()

	problemsMap := make(map[string]string)

	// csv file read
	csvReader := csv.NewReader(f)
	for {
		line, l_err := csvReader.Read()

		// if EOF break loop
		if l_err == io.EOF {
			break
		}
		if l_err != nil {
			log.Fatal(l_err)
		}

		problemsMap[line[0]] = line[1]

	}

	return &problemsMap
}
