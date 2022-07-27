package schema

import (
	"go/ast"
	"reflect"
	"summer/summer_exercise/corm/dialect"
)

type Schema struct {
	Model      any               //存储原结构体
	Name       string            //表名
	FieldsArr  []*Field          //用切片存储Field，遍历 Fields 时会乱序
	Fields     map[string]*Field //记录字段名和 Field 的映射关系，以后查询字段时直接在map中查找
	FieldNames []string          //包含所有的字段名
}

//Field 结构体对应的数据库中每个column的信息
type Field struct {
	Name string
	Type string
	Tag  string //标签 表示 `corm:"<tagName>:<tag>"`
}

//GetField 获取Schema中某个具体的字段信息
func (s *Schema) GetField(name string) *Field {
	return s.Fields[name]
}

// ParseType Parse 将structure结构体中的每一个字段解析成对应dialect中每一个字段并存储到返回的Schema中
func ParseType(structure any, dialect dialect.Dialector) *Schema {
	v := reflect.ValueOf(structure)
	//处理指针
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	schema := &Schema{
		Model: structure,
	}
	//获取类型名 比如 type A struct
	//var a A ==> 这样就可以得到他的类型名是A了，方便以后操作
	schema.Name = v.Type().Name()
	schema.Fields = make(map[string]*Field)

	//处理结构体的每一个字段
	for i := 0; i < v.Type().NumField(); i++ {
		//获取结构体第 i 个字段的信息，赋值到v上，方便操作
		field := v.Type().Field(i)
		//是否是一个嵌入字段或者是否是首字母大写
		//=> 嵌入字段就是一个字段只含有字段类型而没有指定字段的名字
		if !field.Anonymous && ast.IsExported(field.Name) {
			f := &Field{
				Name: field.Name,
				//用field的type去创建一个类型为field.type指向field的reflect.Value 然后再从不同方言的DataType方法获取对应方言数据库中的类型
				Type: dialect.DataType(reflect.Indirect(reflect.New(field.Type))),
			}
			//获取tag标签
			tag, ok := field.Tag.Lookup("corm")
			//对tag的进一步处理 todo
			if ok {
				f.Tag = tag
			}
			//将field装到schema中
			schema.Fields[f.Name] = f
			schema.FieldsArr = append(schema.FieldsArr, f)
			schema.FieldNames = append(schema.FieldNames, f.Name)
		}
	}
	return schema
}

// ParseValue 将schema中结构体的所有值都解析成dialect数据库中对应的值
func (s *Schema) ParseValue(dest any) []any {
	//先获取结构体的value
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() == reflect.Ptr {
		//如果为指针，则拿到他指向的值
		destValue = destValue.Elem()
	}
	var fieldValues []any
	//这里遍历如果是map的话会是乱序，会随运气的好坏报错
	for _, field := range s.FieldsArr {
		//通过结构体的名称来找到值
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
