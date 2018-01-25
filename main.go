package main

import (
	"fmt"
	"time"
)

func main() {
	testRunner()
	testWorker()
}

func testRunner() {
	r := NewRunner(10 * time.Second)

	r.Add(createTask(), createTask(), createTask())

	if err := r.Start(); nil != err {
		fmt.Println("Stopped with error: ", err)
	}
}

func createTask() func(id int) {
	return func(id int) {
		fmt.Println("Hello from ", id)
		time.Sleep(200 * time.Millisecond)
	}
}

func testWorker() {
	var w = NewWorkPool(5)

	w.Run(makeWorker(1))
	w.Run(makeWorker(2))
	w.Run(makeWorker(3))
	w.Run(makeWorker(4))
	w.Run(makeWorker(5))
	w.Run(makeWorker(6))
	w.Run(makeWorker(7))

	w.ShutDown()
}

func makeWorker(sec time.Duration) *WorkerItem {
	return &WorkerItem{sec * time.Second}
}
