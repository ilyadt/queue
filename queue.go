// Реализовать брокер очередей в виде веб сервиса. Сервис должен обрабатывать 2 метода:
//
// 1.	PUT /queue?v=....
// 	положить сообщение в очередь с именем queue (имя очереди может быть любое), пример:
//
// 	curl -XPUT http://127.0.0.1/color?v=red
// 	curl -XPUT http://127.0.0.1/color?v=green
// 	curl -XPUT http://127.0.0.1/name?v=alex
// 	curl -XPUT http://127.0.0.1/name?v=anna
//
// 	в ответ {пустое тело + статус 200 (ok)}
// 	в случае отсутствия параметра v - пустое тело + статус 400 (bad request)
//
// 2.	GET /queue
//
// 	забрать (по принципу FIFO) из очереди с названием queue сообщение и вернуть в теле http запроса, пример (результат, который должен быть при выполненных put’ах выше):
//
// 	curl http://127.0.0.1/color => red
// 	curl http://127.0.0.1/color => green
// 	curl http://127.0.0.1/color => {пустое тело + статус 404 (not found)}
// 	curl http://127.0.0.1/color => {пустое тело + статус 404 (not found)}
// 	curl http://127.0.0.1/name => alex
// 	curl http://127.0.0.1/name => anna
// 	curl http://127.0.0.1/name => {пустое тело + статус 404 (not found)}
//
// 	при GET-запросах сделать возможность задавать аргумент timeout
//
// 	curl http://127.0.0.1/color?timeout=N
//
// 	если в очереди нет готового сообщения получатель должен ждать либо до момента прихода сообщения либо до истечения таймаута (N - кол-во секунд). В случае, если сообщение так и не появилось - возвращать код 404.
// 	получатели должны получать сообщения в том же порядке как от них поступал запрос, если 2 получателя ждут сообщения (используют таймаут), то первое сообщение должен получить тот, кто первый запросил.
//
// Порт, на котором будет слушать сервис, должен задаваться в аргументах командной строки.
// Запрещается пользоваться какими либо сторонними пакетами кроме стандартных библиотек. (задача в написани кода, а не в использовании чужого)
// Желательно (но не обязательно) весь код расположить в одном go-файле (предполагается, что решение будет не больше 200 строк кода) для удобства проверки, никаких дополнительных файлов readme и т.п. не требуется, создание классической структуры каталогов (cmd/internal/...) не требуется.
// Лаконичность кода будет восприниматься крайне положительно, не нужна "гибкость" больше, чем требуется для решения именно этой задачи, не нужны логи процесса работы программы (только обработка ошибок), никакого дебага и т.д... чем меньше кода - тем лучше!
// Оцениваться корректность реализации (заданные условия выполняются), архитектурная составляющая (нет лишних действий в программе, только решающие задачи программы), лаконичность и понятность кода (субъективно, конечно, но думайте о том, насколько будет понятен ваш код для других, это куда более важно в командной разработке, чем сложный "крутой" код).

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Queue --->

var ErrNoElements = errors.New("no elements")

func NewQueue() *Queue {
	return &Queue{
		waitChan: make(map[string]chan string),
		queues:   make(map[string][]string),
	}
}

type Queue struct {
	waitChan map[string]chan string
	queues   map[string][]string
	mu       sync.Mutex // TODO for improving performance each queue can be locked separately
}

func (q *Queue) Add(name, el string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Try to put element to wait queue if anyone is waiting
	select {
	case q.waitChan[name] <- el:
		return
	default:
	}

	q.queues[name] = append(q.queues[name], el)
}

func (q *Queue) Get(name string) (res string, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queues[name]) == 0 {
		return "", ErrNoElements
	}

	res, q.queues[name] = q.queues[name][0], q.queues[name][1:]

	return res, nil
}

func (q *Queue) GetWait(name string) <-chan string {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queues[name]) == 0 {
		if q.waitChan[name] == nil {
			q.waitChan[name] = make(chan string)
		}

		return q.waitChan[name]
	}

	res := make(chan string, 1)
	res <- q.queues[name][0] // queue always has at least one element (previous check)
	close(res)

	q.queues[name] = q.queues[name][1:]

	return res
}

// <--- Queue

// Client Request Queue -->

type ClientReqQueue struct {
	Queue    *Queue
	clientsQ sync.Map
	mu       sync.Mutex
}

func (w *ClientReqQueue) Add(job *Job) {
	ch, ok := w.clientsQ.Load(job.QueueName)

	// Queue initialisation
	if !ok {
		w.mu.Lock()

		ch, ok = w.clientsQ.Load(job.QueueName)
		if !ok {
			ch = make(chan *Job)
			w.clientsQ.Store(job.QueueName, ch)

			go w.process(ch.(chan *Job))
		}

		w.mu.Unlock()
	}

	// Pushing request to queue async
	go func() {
		ch.(chan *Job) <- job
	}()
}

func (w *ClientReqQueue) process(ch chan *Job) {
	for job := range ch {
		// Priority for discarding jobs which are cancelled
		select {
		case <-job.cancelled:
			continue
		default:
		}

		select {
		case <-job.cancelled:
			continue
		case val := <-w.Queue.GetWait(job.QueueName):
			job.Result = val
			close(job.done)

			// Dropped value from queue (race condition)
			// Job cancelled and at the same time value received
			select {
			case <-job.cancelled:
				log.Printf("dropped value `%s` queue=%s", job.Result, job.QueueName)
			default:
			}
		}
	}
}

func NewClientReqQueue(q *Queue) *ClientReqQueue {
	return &ClientReqQueue{
		Queue: q,
	}
}

type Job struct {
	QueueName string
	Result    string

	done      chan struct{}
	cancelled chan struct{}
}

func (r *Job) Done() <-chan struct{} {
	return r.done
}

func (r *Job) Cancel() {
	close(r.cancelled)
}

func NewJob(qName string) *Job {
	return &Job{
		QueueName: qName,
		done:      make(chan struct{}),
		cancelled: make(chan struct{}),
	}
}

/// <---  Client Request Queue

// Handler --->

type Handler struct {
	Queue    *Queue
	ClientsQ *ClientReqQueue
}

func (h *Handler) Request(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		queueName := r.URL.Path[1:]
		if queueName == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		val := r.URL.Query().Get("v")

		if val == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.Queue.Add(queueName, val)
	case http.MethodGet:
		queueName := r.URL.Path[1:]
		if queueName == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		timeout := 0

		if v := r.URL.Query().Get("timeout"); v != "" {
			timeout, _ = strconv.Atoi(v)
		}

		if timeout < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if timeout == 0 {
			el, err := h.Queue.Get(queueName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			_, _ = w.Write([]byte(el))

			return
		}

		// Create new job
		job := NewJob(queueName)

		// Place job in a queue
		h.ClientsQ.Add(job)

		// Creating context with timeout for job cancellation
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
		defer cancel()

		select {
		case <-ctx.Done():
			job.Cancel()
			w.WriteHeader(http.StatusNotFound)
		case <-job.Done():
			_, _ = w.Write([]byte(job.Result))
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// <--- Handler

func main() {
	p := flag.String("p", "2802", "listening port -p 80")

	flag.Parse()

	q := NewQueue()

	h := &Handler{
		ClientsQ: NewClientReqQueue(q),
		Queue:    q,
	}

	http.HandleFunc("/", h.Request)

	addr := "127.0.0.1:" + *p

	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("error starting the server: %s\n", err.Error())

		os.Exit(1)
	}

	fmt.Printf("server started listening on addr: %s\n", addr)

	if err := http.Serve(l, nil); err != nil {
		fmt.Printf("http server error %s\n", err.Error())

		os.Exit(1)
	}
}
