package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func test() {
	c := cron.New()
	i := 1
	entryId, err := c.AddFunc("* * * 1-10/10,1-5/8 4-5,4-5", func() {
		fmt.Println(time.Now(), "每秒执行一次", i)
		i++
	})
	c.Start()
	fmt.Println(entryId, err)
	for _, v := range c.Entries() {
		fmt.Println(v.ID, v.Job, v.Next, v.Prev, v.WrappedJob, v.Schedule)
	}
	time.Sleep(time.Minute * 5)
}
