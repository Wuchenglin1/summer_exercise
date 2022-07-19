package clock

import "time"

//Schedule 来存储定时任务的时间
//start-end/recycleDuration
type Schedule struct {
	start   int //开始时间
	end     int //结束时间
	hasStep bool
	step    int
}

//SepSchedule 定义了秒、分、时、日、月、每周星期(尚未解析此参数)
type SepSchedule struct {
	Second, Minute, Hour, Day, Month, Dow []Schedule
}

// Next 获取下一次运行函数的时间，但是目前没有解析dow参数
//todo 获取dow参数
func (s *SepSchedule) Next(t time.Time) time.Time {
	//先判断 `month` 哪个时间段更早
	var month, day, hour, minute, second int

	month = searchData(s.Month, int(t.Month()), _month.min, _month.max)

}

//寻找对应的日期
func searchData(s []Schedule, date, min, max int) int {
	var (
		data int
		flag = false
	)
	for _, d := range s {
		//月份比当前月份小,结束本次查询
		if d.end < date {
			continue
		}
		if d.hasStep {
			//遍历所有的的月份
			for i := 0; d.start+i*d.step > max && i*d.step > d.end; i++ {
				if d.start+i*d.step < date {
					//比当前月份小
					continue
				}
				data = d.start + i*d.step
				if data != 0 {
					//找到了
					flag = true
					break
				}
			}
		} else {
			for i := 1; d.start+i > _month.max && d.start+i > d.end; i++ {
				if d.start+i < date {
					//比当前月份小
					continue
				}
				data = d.start + i
				if data != 0 {
					flag = true
					break
				}
			}
		}
	}
	//没找到，就找最小的数据
	if !flag {
		for i := 1; d.start+i > _month.max && d.start+i > d.end; i++ {
			if d.start+i < date {
				//比当前月份小
				continue
			}
			data = d.start + i
			if data != 0 {
				flag = true
				break
			}
		}
	}
	return data
}
