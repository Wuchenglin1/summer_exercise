package clock

import (
	"time"
)

//Schedule 来存储定时任务的时间
//start-end/recycleDuration
type Schedule struct {
	start   uint //开始时间
	end     uint //结束时间
	hasStep bool
	step    uint
}

//SepSchedule 定义了秒、分、时、日、月、每周星期(尚未解析此参数)
type SepSchedule struct {
	Second, Minute, Hour, Day, Month, Dow []Schedule
}

// Next 获取下一次运行函数的时间，但是目前没有解析dow参数
//todo 获取dow参数
func (s *SepSchedule) Next(t time.Time) time.Time {
	var (
		loc       = t.Location()
		added     = false
		yearLimit = t.Year() + 5
	)

	//参考了cron框架源码，在王鑫的讲解下写出的代码
	//https://github.com/robfig/cron/blob/master/spec.go
	//真的这个结构写不出来了:cry: 大致意思就是通过不断地循环来迭代出一个距离现在最近的时间点来开启任务
WRAP:
	if t.Year() > yearLimit {
		return t
	}

	for !s.match(s.Month, uint(t.Month())) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)
		}
		t = t.AddDate(0, 1, 0)

		if t.Month() == time.January {
			goto WRAP
		}
	}

	for !s.match(s.Day, uint(t.Day())) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
		}
		t = t.AddDate(0, 0, 1)

		if t.Hour() != 0 {
			if t.Hour() > 12 {
				t = t.Add(time.Duration(24-t.Hour()) * time.Hour)
			} else {
				t = t.Add(time.Duration(-t.Hour()) * time.Hour)
			}
		}
		if t.Day() == 1 {
			goto WRAP
		}
	}
	for !s.match(s.Hour, uint(t.Hour())) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, loc)
		}
		t = t.Add(time.Hour)
		// need recheck
		if t.Hour() == 0 {
			goto WRAP
		}
	}
	for !s.match(s.Minute, uint(t.Minute())) {
		if !added {
			added = true
			t = t.Truncate(time.Minute)
		}
		t = t.Add(1 * time.Minute)

		if t.Minute() == 0 {
			goto WRAP
		}
	}
	for !s.match(s.Second, uint(t.Second())) {
		if !added {
			added = true
			t = t.Truncate(time.Second)
		}
		t = t.Add(1 * time.Second)
		if t.Second() == 0 {
			goto WRAP
		}
	}
	return t.In(loc)
}

//math 查看schedule是否存在该key
func (s *SepSchedule) match(schedule []Schedule, key uint) bool {
	for _, sd := range schedule {
		if sd.hasStep {
			for i := uint(0); i < sd.end && sd.start+i*sd.step < sd.end; i++ {
				if sd.start+i*sd.step == key {
					return true
				}
			}
		} else {
			for i := uint(0); i < sd.end; i++ {
				if sd.start+i == key {
					return true
				}
			}
		}
	}
	return false
}
