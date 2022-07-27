package clause

import (
	"fmt"
	"strings"
)

/*
生成语句 VALUES(?,?,?,...)
*/

//长度为n，传入的是[]any，存储插入的值，可以同时传多个[]any来解析，返回sql语句和参数
func _values(values ...any) (string, []any) {
	var (
		rm   string          //request mark => rm :)
		sql  strings.Builder //构建sql语句的builder
		args []any           //需要插入的参数
	)
	//先写一个 VALUES 埋伏他一手
	sql.WriteString("VALUES")
	for i, value := range values {
		v := value.([]any) //这里是将单个要插入的value的值反射出来
		//如果是多个传入值的可以先对rm进行判断是否有值，因为类型是一样，所以只需要生成一次rm就行了
		if rm == "" {
			//如果没有值的话就生成值
			rm = generateRM(len(v))
		}
		//再将问号写到builder里面
		sql.WriteString(rm)
		if i+1 != len(values) {
			//如果还没有遍历完的话，就加个小逗号 ", "
			sql.WriteString(",")
		}
		//将参数存到args里面
		args = append(args, v...)
	}
	return sql.String(), args
}

//生成语句(?,?,...)，先生成括号，再在括号里生成num个?
func generateRM(num int) string {
	var str []string
	for i := 0; i < num; i++ {
		//	有多少问号就加多少个问号
		str = append(str, "?")
	}
	//封装成(?,?,...)
	return fmt.Sprintf("(%v)", strings.Join(str, ","))
}
