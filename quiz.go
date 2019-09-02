package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Question struct {
	Text   string
	Answer int
}

var correctAnswers int

var filename = flag.String("f", "problems.csv", "Quiz questions file")
var shuffle = flag.Bool("s", false, "Shuffle questions order")
var quizDuration = 20 * time.Second

func main() {
	flag.Parse()
	input := bufio.NewScanner(os.Stdin)
	fmt.Printf("Welcome to my quiz!\nYou have %s to finish this quiz\n --- press enter key to start ---", quizDuration)
	input.Scan()
	questions := getQuestions(*filename)
	ch := make(chan bool)
	go askQuestions(questions, input, *shuffle, ch)
	select {
	case <-ch:
	case <-time.After(quizDuration):
		fmt.Println("\nTime Out!")
	}
	fmt.Printf("Quiz Finished\nYour results: %d / %d", correctAnswers, len(questions))
}

func getQuestions(filename string) []Question {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "opening file: %v\n", err)
	}
	reader := csv.NewReader(bufio.NewReader(file))
	var questions []Question
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		question := Question{
			Text: row[0],
		}
		question.Answer, _ = strconv.Atoi(row[1])
		questions = append(questions, question)
	}
	return questions
}

func askQuestions(questions []Question, input *bufio.Scanner, shuffle bool, ch chan<- bool) {
	var userAnswer int
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
	}

	for _, question := range questions {
		fmt.Printf("%s ? ", question.Text)
		if input.Scan() {
			userAnswer, _ = strconv.Atoi(input.Text())
		}
		if userAnswer == question.Answer {
			correctAnswers++
		}
	}
	ch <- true
}
