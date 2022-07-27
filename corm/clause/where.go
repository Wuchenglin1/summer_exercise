package clause

import "fmt"

//长度为1，为一个[]any，其中values[1]是desc，values[1:]是args
func _where(values ...any) (string, []any) {
	//where复合查询 todo
	v := values[0].([]any)
	desc, vas := v[0], v[1:]
	return fmt.Sprintf("WHERE %s", desc), vas
}
