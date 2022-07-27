package clause

import (
	"fmt"
	"strings"
)

//长度为2，第一个参数是tableName，第二个参数是存fields的[]string，返回一个空的[]any{}
func _select(values ...any) (string, []any) {
	tableName := values[0]
	//给每个field之间加上 ","
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %v", fields, tableName), []any{}
}
