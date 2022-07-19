package clock

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//bounds second minute hour day month dotw 的范围
type bounds struct {
	max, min uint //最大，最小
}

type Parser struct {
	IsWithSecond bool //是否开启 秒 级别的时钟 如果开启的话，最后一位读的是秒

}

var (
	_second = bounds{max: 59, min: 0} //秒
	_minute = bounds{max: 59, min: 0} //每分钟
	_hour   = bounds{max: 59, min: 0} //每时
	_day    = bounds{max: 31, min: 1}
	_month  = bounds{max: 12, min: 1} //每月
	_dow    = bounds{max: 7, min: 1}  //每周星期 day of the week
)

/*

Minutes：分钟，取值范围[0-59]；
Hours：小时，取值范围[0-23]；
Day of month：每月的第几天，取值范围[1-31]；
Month：月，取值范围[1-12]或者使用月份名字缩写[JAN-DEC]；
Day of the week：周历，取值范围[0-6]或名字缩写[JUN-SAT]。


*/

//ParseSep 解析cron表达式
func (p Parser) ParseSep(sep string) (SepSchedule, error) {
	var (
		ch       = make(chan error, 6)
		err      error
		schedule = SepSchedule{}
	)
	if len(sep) == 0 {
		return schedule, errors.New("sep为空")
	}
	//先去掉所有空格
	fields := strings.Fields(sep)
	//先检查长度是否正确
	if p.IsWithSecond {
		//开启秒级 6 个参数
		if len(fields) != 6 {
			return schedule, fmt.Errorf("已经开启了秒级定时器，但sep参数个数不为6，而是%v", len(fields))
		}
		fields = []string{fields[1], fields[2], fields[3], fields[4], fields[5], fields[0]}
	} else {
		//不开启秒级 5 个参数
		if len(fields) != 5 {
			return schedule, fmt.Errorf("没有开启秒级定时器，但是sep参数个数不为5，而是%v", len(fields))
		}
	}

	f := func(exp string, b bounds) (s []Schedule) {
		s, err = GetField(exp, b)
		if err != nil {
			ch <- err
			return nil
		}
		return s
	}

	//获取时间间隔

	var (
		second []Schedule
		minute = f(fields[0], _minute)
		hour   = f(fields[1], _hour)
		day    = f(fields[2], _day)
		month  = f(fields[3], _month)
		dotw   = f(fields[4], _dow)
	)
	if p.IsWithSecond {
		second = f(fields[5], _second)
	}
	close(ch)
	if len(ch) != 0 {
		return schedule, <-ch
	}
	return SepSchedule{
		Second: second,
		Minute: minute,
		Hour:   hour,
		Day:    day,
		Month:  month,
		Dow:    dotw,
	}, err
}

//GetField 只支持解析特殊字符 `,`  `/`  `-`  `*`
//这里出现了错误就直接panic掉了
func GetField(field string, b bounds) (s []Schedule, err error) {
	/*
			思路：
		1.先去掉字符中的 `,` 将多个值或范围逐个处理
		2.再去掉 `/` 以获取步长（两次触发的时间间隔）
		3.再去掉 `-`来获得时间范围
	*/

	//首先先去掉逗号
	fields := strings.FieldsFunc(field, func(r rune) bool {
		return r == ','
	})

	s = make([]Schedule, len(fields))

	for k, exp := range fields {
		//对每一个值都单独处理
		var (
			step       = strings.Split(exp, "/")     //步长，相当于每隔多久执行一次
			minAndMax  = strings.Split(step[0], "-") //去掉减号后的最小和最大值
			fieldIsOne = len(minAndMax) == 1         //如果等于1的话就相当于 `exp` 中没有 `-`
		)

		//先解析 `*` 和 `-` 特殊字符
		if minAndMax[0] == "*" {
			//不管是 */x 还是 * 都属于这种情况
			//每 field 都执行一次这个任务
			s[k].start = b.min
			s[k].end = b.max
			s[k].step = 1
		} else {
			//剩下的就是 x/y  x-y/z
			s[k].start, err = parseInt(minAndMax[0])
			if err != nil {
				return nil, err
			}
			switch len(minAndMax) {
			case 1:
				// x 或者 x/y 的形式 结束时间和开始时间设为相同的
				// 在 field 内的 x 执行一次任务
				s[k].step = 1
				s[k].end, err = parseInt(minAndMax[0])
				if err != nil {
					return nil, err
				}
			case 2:
				// x-y 或者 x-y/z 的形式 解析 y
				s[k].end, err = parseInt(minAndMax[1])
				if err != nil {
					return nil, err
				}
				if s[k].end == s[k].start {
					return nil, fmt.Errorf("`-`参数两边数字不能相同")
				}
			default:
				return nil, fmt.Errorf("%v `-`参数太多啦", exp)
			}
		}
		//再解析 `/` 特殊字符，即循环时间
		switch len(step) {
		case 1:
			//没有 `/`
			s[k].hasStep = false
			//s[k].end = b.max
		case 2:
			//有 `/`
			s[k].hasStep = true
			s[k].step, err = parseInt(step[1])
			if err != nil {
				return nil, err
			}
			if fieldIsOne {
				s[k].end = b.max
			}
		default:
			return nil, fmt.Errorf("%v 的 `/` 参数太多啦", exp)
		}
		if s[k].start > s[k].end {
			return nil, fmt.Errorf("开始时间 %v 不可以大于结束时间 %v ！", s[k].start, s[k].end)
		}
		if s[k].start < b.min {
			return nil, fmt.Errorf("开始时间 %v 不可以小于最小时间 %v ", s[k].start, b.min)
		}
		if s[k].end > b.max {
			return nil, fmt.Errorf("结束的时间 %v 不可以大于最大时间 %v", s[k].end, b.max)
		}
		if s[k].step == 0 {
			return nil, fmt.Errorf("步长不可以为0")
		}
	}
	return s, err
}

func parseInt(exp string) (uint, error) {
	n, err := strconv.Atoi(exp)
	if err != nil {
		return 0, fmt.Errorf("不能将 %v 解析成数字 错误 : %v", exp, err)
	}
	if n < 0 {
		return 0, fmt.Errorf("%v 是一个小于0的数字，无法设置时长 ", n)
	}
	return uint(n), nil
}
