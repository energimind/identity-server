package config

import (
	"fmt"
	"reflect"
	"strings"
)

// Dump dumps the configuration to a string slice.
// Sensible fields are masked.
func Dump(cfg *Config) []string {
	sections := make([]string, 0)

	sectionType := reflect.TypeOf(*cfg)
	sectionValue := reflect.ValueOf(*cfg)

	for i := 0; i < sectionType.NumField(); i++ {
		sectionName := sectionType.Field(i).Name
		fieldValue := sectionValues(sectionValue.Field(i).Interface())

		sections = append(sections, fmt.Sprintf("%s {%s}", sectionName, fieldValue))
	}

	return sections
}

func sectionValues(section any) string {
	t := reflect.TypeOf(section)
	v := reflect.ValueOf(section)

	fields := make([]string, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.Field(i).Interface()

		if isProtectedField(fieldName) {
			fieldValue = "****"
		}

		fields = append(fields, fmt.Sprintf("%s=%v", fieldName, fieldValue))
	}

	return strings.Join(fields, " ")
}

func isProtectedField(fieldName string) bool {
	return strings.Contains(fieldName, "Password") || strings.Contains(fieldName, "Secret")
}
