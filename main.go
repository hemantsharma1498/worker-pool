package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const MAXINPUTTRIES = 5

type WorkerPool struct {
	Jobs    chan File
	Results chan Result
	Workers int
}

type File struct {
	Name string
	Size int
}

type Result struct {
	FileName string
	Status   string
}

func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{Jobs: make(chan File), Results: make(chan Result), Workers: workers}
}

func (wp *WorkerPool) Run() {
	for i := 0; i < wp.Workers; i++ {
		go func(workerID int) {
			for job := range wp.Jobs {
				fmt.Printf("Worker %d processing %s\n", workerID, job.Name)
				result := ProcessFile(job)
				wp.Results <- result
			}
		}(i)
	}
}

func ProcessFile(file File) Result {
	processTime := time.Duration(file.Size) * time.Millisecond
	time.Sleep(processTime)

	num := rand.Intn(100)
	if num%2 == 0 && num%4 == 0 {
		return Result{FileName: file.Name, Status: "Error"}
	}
	return Result{FileName: file.Name, Status: "Processed"}
}

func getInput() (int, error) {
	var poolSize int
	var try int
	for try = 0; try < MAXINPUTTRIES; try++ {
		fmt.Print("Enter pool size: ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		poolSize, err = strconv.Atoi(input)
		if err != nil {
			fmt.Println("Please enter a valid number")
			continue
		}

		if poolSize <= 0 {
			fmt.Println("Pool size must be a positive number")
			continue
		}
		break
	}
	if try == MAXINPUTTRIES {
		return -1, errors.New("encountered an error or too many invalid attemps")
	}
	return poolSize, nil
}

func main() {
	fmt.Print("Simulating worker pool processing sample files\n")
	poolSize, err := getInput()
	if err != nil {
		fmt.Println("Exhausted input tries")
		os.Exit(0)
	}

	wp := NewWorkerPool(poolSize)
	go wp.Run()

	go func() {
		for i := 0; i < 10; i++ {
			wp.Jobs <- File{Name: fmt.Sprintf("file_%d.txt", i), Size: rand.Intn(100) + 10}
		}
		close(wp.Jobs)
	}()

	for i := 0; i < 10; i++ {
		result := <-wp.Results
		fmt.Printf("File %s: %s\n", result.FileName, result.Status)
	}
}
