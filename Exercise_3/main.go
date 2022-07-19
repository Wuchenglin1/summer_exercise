package main

import (
	"fmt"
	"summer/summer_exercise/Exercise_3/clock"
	"time"
)

func main() {
	cron := clock.New(clock.WithSecond(true))
	i, err := cron.AddFunction("1/5 * * * * *", func() {
		fmt.Println("啦啦啦啦啦啦我是")
	})
	if err != nil {
		fmt.Println(i, err)
		return
	}
	i, err = cron.AddFunction("2/5 * * * * *", func() {
		fmt.Println("呜呜呜呜呜呜呜呜")
	})
	if err != nil {
		fmt.Println(i, err)
		return
	}
	cron.Run()
	for {
		fmt.Println(time.Now())
		time.Sleep(time.Second * 1)
	}
}
