package clause

import (
	"fmt"
	"strings"
)

//长度为2，第一个参数是表名string，第二个是map类型的参数，表示待更新的键值对
func _update(values ...any) (string, []any) {
	var (
		keys []string
		vas  []any
	)
	tableName := values[0]
	m := values[1].(map[string]any)
	for k, v := range m {
		//<keyName>= ?
		keys = append(keys, k+"= ?")
		//vas[k] = ["<value>"]
		vas = append(vas, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vas
}
