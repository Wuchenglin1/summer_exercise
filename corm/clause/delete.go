package clause

import "fmt"

//长度为1，表名string
func _delete(values ...any) (string, []any) {
	//DELETE FROM <table_name>
	return fmt.Sprintf("DELETE FROM %s ", values[0]), []any{}
}
