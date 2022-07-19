package clock

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Cron struct {
	IsRunning bool   //是否已经启动
	ID        int    //启动定时任务的ID
	Parse     Parser //解析字段
	Lock      sync.Mutex
	nextID    int
	Tasks     []*Task
	add       chan *Task
}

type Task struct {
	TaskID   int
	Schedule SepSchedule
	Next     time.Time //下一次运行的时间
	Prev     time.Time //上一次运行的时间
	Job      func()
}
type sortByTime []*Task

//New 使用函数式选项模式来编写配置文件以自定义添加配置和默认配置
func New(Opt ...Option) *Cron {
	cron := &Cron{
		add:       make(chan *Task, 100),
		IsRunning: false,
		Parse:     Parser{IsWithSecond: false},
	}
	for _, c := range Opt {
		c(cron)
	}
	return cron
}

func (c *Cron) Run() {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if c.IsRunning {
		//如果闹钟已经启动，则无事发生
		return
	}
	//如果还没有启动的话就更改cron状态
	c.IsRunning = true
	go c.startUpTask()
}

// AddFunction AddFunc 通过解析sep的表达式来获取定时任务的间隔等等，function为自定义函数
//返回启动Cron的任务ID和nil（如果error为空的话）
func (c *Cron) AddFunction(sep string, function func()) (int, error) {
	sepS, err := c.Parse.ParseSep(sep)
	fmt.Println(sepS)
	if err != nil {
		return 0, fmt.Errorf("%v 解析失败: %v", sep, err)
	}
	return c.Schedule(sepS, function), nil
}

//Schedule 将解析的时间存入线程池（数组）
//返回任务的id
func (c *Cron) Schedule(schedule SepSchedule, job func()) int {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.nextID++
	task := &Task{
		TaskID:   c.nextID,
		Schedule: schedule,
		Job:      job,
	}
	if !c.IsRunning {
		//未启动
		c.Tasks = append(c.Tasks, task)
	} else {
		//如果已经启动，则塞到开始任务的管道里面
		c.add <- task
	}
	return c.nextID
}

func (c *Cron) startUpTask() {
	nowTime := time.Now()
	for _, task := range c.Tasks {
		task.Next = task.Schedule.Next(nowTime)
		fmt.Println(task.TaskID, task.Next)
	}
	//这里借鉴了golang cron包的 `run` 实现方法 https://github.com/robfig/cron/blob/master/cron.go
	for {
		//先按照时间顺序，把各个任务的执行时间先排序，然后依次执行各个函数任务
		//使用sort包排序
		sort.Sort(sortByTime(c.Tasks))

		//todo 使用更好的排序方法进行时间排序
		var newTimer *time.Timer

		if len(c.Tasks) == 0 {
			newTimer = time.NewTimer(time.Hour * 24 * 100000)
		} else {
			//用下一次执行程序的时间与现在时间计算间隔得到等待时间
			newTimer = time.NewTimer(c.Tasks[0].Next.Sub(nowTime))
		}

		for {
			select {
			case nowTime = <-newTimer.C:
				nowTime = nowTime.In(time.Local)

				for _, t := range c.Tasks {
					if t.Next.IsZero() {
						//下一次执行的时间为零即任务完成之后跳出循环
						break
					}
					t.setup()
					t.Prev = t.Next
					t.Next = t.Schedule.Next(nowTime)
				}
			case newTask := <-c.add:
				nowTime = time.Now()
				newTimer.Stop()
				newTask.Next = newTask.Schedule.Next(nowTime)
				c.Tasks = append(c.Tasks, newTask)
				//刷新时间
			}
			break
		}
	}
}

func (t *Task) setup() {
	go func() {
		i := 0
		fmt.Println(i)
		i++
		t.Job()
	}()
}

func (s sortByTime) Len() int      { return len(s) }
func (s sortByTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortByTime) Less(i, j int) bool {
	if s[i].Next.IsZero() {
		return false
	}
	if s[j].Next.IsZero() {
		return true
	}
	return s[i].Next.Before(s[j].Next)
}
