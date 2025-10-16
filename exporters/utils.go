package exporters

import (
	"reflect"
	"slices"
)

func getStructFields(item reflect.Value) []string {
	t := item.Type()

	fields := []string{}
	for i := range t.NumField() {
		fields = append(fields, t.Field(i).Name)
	}

	return fields
}

func getStructValues(dataValue reflect.Value, fields []string) []any {
	values := []any{}

	for _, field := range fields {
		f := dataValue.FieldByName(field)
		if !f.IsValid() {
			values = append(values, "N/A")
			continue
		}

		if f.Kind() == reflect.Ptr && !f.IsNil() {
			f = f.Elem()
		}

		values = append(values, f.Interface())
	}

	return values
}

func contains(slice []string, value string) bool {
	return slices.Contains(slice, value)
}
