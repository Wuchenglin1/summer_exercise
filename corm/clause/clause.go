/*
 * @Author: wcl
 * @Data: 2022/7/25 16:15
 * @Desc: 解析sql语句
 */

package clause

import (
	"strings"
)

//Clause 用来存储分割了的sql语句和参数
type Clause struct {
	sql     map[Type]string
	sqlArgs map[Type][]any
}

//Type sql分割后语句对应的Type
type Type int

//注册需要满足的方法，传入的值的要求在每个函数前面都已经写了，返回的值第一个参数是sql语句，第二个要插入的参数
type generator func(values ...any) (string, []any)

//这里已经注册了的方法
var generators map[Type]generator

const (
	INSERT Type = iota + 1
	VALUES
	SELECT
	DELETE
	UPDATE
	LIMIT
	WHERE
	ORDERBY
	COUNT
)

func init() {
	//这里相当于一个注册函数，向generators里面注册各种分割了的sql语句的方法，后面会用到
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[DELETE] = _delete
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[COUNT] = _count
}

// Set 根据对应的Type来调用对应的generators[Type]生成sql语句
func (c *Clause) Set(name Type, values ...any) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlArgs = make(map[Type][]any)
	}
	//先生成sql语句
	sql, args := generators[name](values...)
	//然后存储sql语句
	c.sql[name] = sql
	//存储对应的参数
	c.sqlArgs[name] = args
}

// Build 这里需要对传入的orders进行排序，因为map中存的数据是无序的
func (c *Clause) Build(orders ...Type) (string, []any) {
	var (
		sqls []string
		args []any
	)
	//先遍历所有的orders即所有的sql语句
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			//如果执行的sql语句中存在order，将sql语句叠加到最后的sql语句里面去
			sqls = append(sqls, sql)
			//将参数全部叠加到最后的参数列表里去，防止有[]any没有参数还来占位
			if len(c.sqlArgs) != 0 {
				args = append(args, c.sqlArgs[order]...)
			}
		}
	}
	//对sql语句加空格 " " 来隔开
	return strings.Join(sqls, " "), args
}
