package tool

import (
	"reflect"
	"summer/summer_exercise/corm/schema"
)

// GetColumnDiff 返回dest中有，src中没有的字段
func GetColumnDiff(dest, src []*schema.Field) map[string]*schema.Field {
	var modifyMap = make(map[string]*schema.Field)
	fieldMap := make(map[string]*schema.Field)
	//先将src的字段全部塞到map里面去
	for _, v := range src {
		fieldMap[v.Name] = v
	}
	//如果src中没有dest中的字段的话，就在最后的切片里面添加该字段
	for _, v := range dest {
		_, ok := fieldMap[v.Name]
		if !ok {
			//如果在src中没有该字段的话，那么就把该字段添加进去
			modifyMap[v.Name] = v
		}
		equal := reflect.DeepEqual(v, fieldMap[v.Name])
		//两个值如果不深度相等，那么就取新的结构体中的字段
		if !equal {
			modifyMap[v.Name] = v
		}
	}
	return modifyMap
}
