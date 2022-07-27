package clause

import (
	"fmt"
	"strings"
)

/*
生成语句 INSERT INTO <TableName> (<column1>,<column2>,...)
*/

//长度为2，第一个参数是表名，第二个参数是一个[]string，表示插入的column，的[]any为空
func _insert(values ...any) (string, []any) {
	tableName := values[0]
	//将所有插入的值用,进行分割
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []any{}
}
