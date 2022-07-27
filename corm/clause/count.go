package clause

import "fmt"

//长度为1，只有一个表名
func _count(values ...any) (string, []any) {
	return fmt.Sprintf("SELECT COUNT(*) FROM %v", values[0]), nil
}
