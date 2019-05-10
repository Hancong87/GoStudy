package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	id int
	randomno int
}

type Result struct {
	job Job
	sumofdigits int
	wokerId int
}

var jobs = make(chan Job, 10)
var results = make(chan Result, 10)

func digits(number int) int {
	sum := 0
	no := number
	for no != 0 {
		digit := no % 10
		sum += digit
		no /= 10
	}
	time.Sleep(2 * time.Second)
	return sum
}

func worker(wg *sync.WaitGroup, workerId int) {
	for job := range jobs {
		output := Result{job, digits(job.randomno), workerId}
		if workerId % 2 == 0 {
			time.Sleep(2 * time.Second)
		}
		results <- output
	}
	fmt.Println("finishing: worker id " , workerId)
	wg.Done()
}

func createWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg, i)
	}
	wg.Wait()
	close(results)
}

func allocate(noOfJobs int) {
	for i:= 0; i < noOfJobs; i++ {
		randomno := rand.Intn(999)
		job := Job{i, randomno}
		jobs <- job
	}
	close(jobs)
}

func result(done chan bool) {
	for result := range results {
		fmt.Printf("Job id %d, Worker id %d, input random no %d, sum of digits %d\n", result.job.id, result.wokerId, result.job.randomno, result.sumofdigits)
	}
	done <- true
}

func main() {
	startTime := time.Now()
	noOfJobs := 100
	go allocate(noOfJobs)
	done := make(chan bool)

	// why cant I put go result(done) after createWorkerPool ???
	go result(done)
	noOfWorkders := 10
	createWorkerPool(noOfWorkders)
	//go result(done)

	<-done
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time taken ", diff.Seconds(), "seconds")
}

