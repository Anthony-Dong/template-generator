package utils

import (
	"reflect"
)

func ObjToMap(bean interface{}, filter func(fieldName string) string) map[string]interface{} {
	value := reflect.Indirect(reflect.ValueOf(bean))
	if value.Kind() != reflect.Struct {
		panic("the bean mush struct")
	}
	_type := value.Type()
	fieldNum := value.NumField()
	_map := make(map[string]interface{}, fieldNum)
	for x := 0; x < fieldNum; x++ {
		field := _type.Field(x)
		_map[filter(field.Name)] = value.Field(x).Interface()
	}
	return _map
}
