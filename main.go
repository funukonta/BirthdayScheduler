package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func hello(name string) {
	message := fmt.Sprintf("Hi, %v", name)
	fmt.Println(message)
}

func runCronJobs() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Seconds().Do(func() {
		hello("John Doe")
	})

	s.StartBlocking()
}

func main() {
	runCronJobs()
}
