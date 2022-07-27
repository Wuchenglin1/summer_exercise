package dialect

import (
	"reflect"
)

var dialectsMap = map[string]Dialector{}

type Dialector interface {
	DataType(Type reflect.Value) string
	TableExistSql(TableName string) (string, []any)
}

//RegisterDialect 注册支持的数据库列表
func RegisterDialect(name string, d Dialector) {
	dialectsMap[name] = d
}

// GetDialect 获取支持的方言列表
func GetDialect(name string) (d Dialector, ok bool) {
	d, ok = dialectsMap[name]
	return
}
