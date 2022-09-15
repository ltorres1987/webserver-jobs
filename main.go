package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Job struct {
	Name   string
	Delay  time.Duration
	Number int
}

type Worker struct {
	Id         int
	JobQueue   chan Job
	WorkerPool chan chan Job
	QuitChan   chan bool
}

type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int
	JobQueue   chan Job
}

func NewWorker(id int, workerPool chan chan Job) *Worker {
	return &Worker{
		Id:         id,
		JobQueue:   make(chan Job),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {

			w.WorkerPool <- w.JobQueue //agregar trabajo al grupo
			select {
			case job := <-w.JobQueue:
				fmt.Println("Start worker # : ", w.Id)
				fib := Fibonnaci(job.Number)
				time.Sleep(job.Delay)
				fmt.Println("end worker # : ", w.Id, " fibonacci: ", fib)

			case <-w.QuitChan:
				fmt.Println("finalizado worker # : ", w.Id)
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func Fibonnaci(n int) int {

	if n <= 1 {
		return n
	}
	return Fibonnaci(n-1) + Fibonnaci(n-2)
}

func NewDispatcher(maxWorkers int, jobQueue chan Job) *Dispatcher {
	worker := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		MaxWorkers: maxWorkers,
		JobQueue:   jobQueue,
		WorkerPool: worker,
	}
}

func (d *Dispatcher) Dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			go func() {
				workerJobQueue := <-d.WorkerPool // obtener trabajador del grupo
				workerJobQueue <- job            //Los trabajadores leerán de JOB
				fmt.Println("dispatcher: ", job)
			}()
		}
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(i, d.WorkerPool)
		worker.Start()
	}

	go d.Dispatch()
}

func ReuqestHandler(w http.ResponseWriter, r *http.Request, jobQueue chan Job) {

	if r.Method == "POST" {

		delay, err := time.ParseDuration(r.FormValue("delay"))
		if err != nil {
			http.Error(w, "Invalid delay", http.StatusBadRequest)
			return
		}

		value, err := strconv.Atoi(r.FormValue("value"))
		if err != nil {
			http.Error(w, "Invalid value", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		if name == "" {
			http.Error(w, "Invalid name", http.StatusBadRequest)
			return
		}

		job := Job{
			Name:   name,
			Delay:  delay,
			Number: value,
		}

		jobQueue <- job
		w.WriteHeader(http.StatusCreated)
	} else {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {

	const (
		maxWorker    = 4
		maxQueueSize = 20
		port         = ":8081"
	)
	jobQueue := make(chan Job, maxQueueSize)
	dispatcher := NewDispatcher(maxWorker, jobQueue)

	dispatcher.Run()

	http.HandleFunc("/fib", func(writer http.ResponseWriter, request *http.Request) {
		ReuqestHandler(writer, request, jobQueue)
	})

	log.Fatal(http.ListenAndServe(port, nil))
}
