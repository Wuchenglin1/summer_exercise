package clause

import "fmt"

//长度为1，orderBy的参数
func _orderBy(values ...any) (string, []any) {

	return fmt.Sprintf("ORDER BY %s", values[0]), []any{}
}
