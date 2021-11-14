package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Me costo un poco entenderlo, pero ya le agarre la onda!
// 
// Unos puntos importantes.
// 
// El Dispacher obtendra una cola, que contedra una cola de cada uno de los workers.
// Una vez que el dispacher reciba un trabajo, esta agarrara de su cola a un worker, y mandara la tarea por ese medio
// El worker, procesara el trabajao y volvera a agregar su cola de trabajo a la cola del dispacher
// Tambien cree un repositorio en GitHub, con comentarios y con un CodeTour viendo que hace el codigo paso a paso. Pueden clonar mi repo y visualizarlo instalando la extension: https://github.com/CarlosTrejo2308/platzigoserver
// 
// Espero que les sirva para entenderlo mejor!


type Job struct {
	Name string
	Delay time.Duration
	Number int
}

type Worker struct {
	Id int
	JobQueue chan Job
	WorkerPool chan chan Job
	QuitChan chan bool
}

func NewWorker(id int, workerPool chan chan Job) * Worker {
	return &Worker{
		Id: id,
		WorkerPool: workerPool,
		JobQueue: make(chan Job),
		QuitChan: make(chan bool),
	}
}

func (w Worker) Start( ) {
	go func () {
		for {
			w.WorkerPool <- w.JobQueue

			select {
			case job := <- w.JobQueue:
				fmt.Printf("-----------------------------------")		
				fmt.Printf("Worker with id %d Started\n", w.Id)		
				fib := Fibonacci(job.Number)
				time.Sleep(job.Delay)
				fmt.Printf("Worker with id %d Finished with result %d\n", w.Id, fib)	
			case <- w.QuitChan:
				fmt.Printf("Worker with id %d Stopped %d\n", w.Id)	
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int
	JobQueue chan Job
}

func NewDispatcher(jobQueue chan Job, maxWorkers int) *Dispatcher {
  worker := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		JobQueue: jobQueue,
		MaxWorkers: maxWorkers,
		WorkerPool: worker,
	}
}


func (d *Dispatcher) Dispatch() {
	for {
		select{
		case job := <- d.JobQueue:
			go func() {
				workerJobQueue := <-d.WorkerPool
				workerJobQueue <- job
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

func Fibonacci(n int) int {
	if n <= 1{
		return n
	}

	return Fibonacci(n-1) + Fibonacci(n -2)
}

func RequestHandler(w http.ResponseWriter, r *http.Request, jobQueue chan Job) {
	if r.Method != "POST" { // GET, PUT, DELETE
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	delay, err := time.ParseDuration(r.FormValue("delay"))
	if err != nil {
		http.Error(w, "Invalid Delay", http.StatusBadRequest)
		return
	}

	value, err := strconv.Atoi(r.FormValue("value"))
	if err != nil {
		http.Error(w, "Invalid Value", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")

	if name == "" {
		http.Error(w, "Invalid Name", http.StatusBadRequest)
		return
	}

	job := Job{Name: name, Delay: delay, Number: value}
	jobQueue <- job
	w.WriteHeader(http.StatusCreated)
}

// func RequestHandler(jobQueue chan Job) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request)
// }
func main() {
	const (
		maxWorkers = 4
		maxQueueSize = 20
		port = ":4000"
	)

	//Buffered Channel: cantidad máxima que nuestro programa va a manejar
	jobQueue := make(chan Job, maxQueueSize);

	dispatcher := NewDispatcher(jobQueue, maxWorkers)
	dispatcher.Run()

	// http://localhost:8081/fib
	http.HandleFunc("/fib", func(w http.ResponseWriter, r *http.Request) {
		RequestHandler(w, r, jobQueue)
	})

	log.Fatal(http.ListenAndServe(port, nil))
}
