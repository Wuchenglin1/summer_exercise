package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct {
}

const (
	tinyint     = "tinyint"
	smallint    = "smallint"
	integer     = "int"
	bigint      = "bigint"
	tinyUint    = "tinyint unsigned"
	smallUint   = "smallint unsigned"
	unsignedInt = "int unsigned"
	bigUint     = "bigint unsigned"
	varchar     = "varchar(256)"
	float       = "float"
	double      = "double"
	boolean     = "bool"
	dataTime    = "datetime"
)

//检查mysql是否实现了dialect的方法
var _ Dialector = (*mysql)(nil)

func init() {
	RegisterDialect("mysql", mysql{})
}

//DataType 查询对应field对应mysql中的类型
func (m mysql) DataType(field reflect.Value) string {
	switch field.Kind() {
	case reflect.Int8:
		return tinyint
	case reflect.Int16:
		return smallint
	case reflect.Int32:
		return integer
	case reflect.Int64, reflect.Int:
		return bigint
	case reflect.Uint8:
		return tinyUint
	case reflect.Uint16:
		return smallUint
	case reflect.Uint32:
		return unsignedInt
	case reflect.Uint64, reflect.Uint:
		return bigUint
	case reflect.String:
		return varchar
	case reflect.Float32:
		return float
	case reflect.Float64:
		return double
	case reflect.Bool:
		return boolean
	case reflect.Struct:
		if _, ok := field.Interface().(time.Time); ok {
			return dataTime
		}
	}
	panic(fmt.Sprintf("unsupport type : %v(name: %v)", field.Kind(), field.Type().Name()))
}

//TableExistSql 返回mysql查询在数据库中是否存在该表名的表
func (m mysql) TableExistSql(tableName string) (string, []any) {
	args := []any{tableName}
	return "SELECT table_name FROM information_schema.TABLES WHERE table_name = ?", args
}
