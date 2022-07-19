package main

import (
	"fmt"
	"log"
	"summer/summer_exercise/Exercise_3/clock"
	"time"
)

func main() {
	cron := clock.New(clock.WithSecond(true))
	i, err := cron.AddFunction("1 * * * * *", func() {
		fmt.Println("5s")
	})
	if err != nil {
		log.Fatalf("error : %v", err)
	}
	fmt.Println(i)
	time.Sleep(time.Second * 10)
}
