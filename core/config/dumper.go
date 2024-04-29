package config

import (
	"fmt"
	"reflect"
	"strings"
)

// Section represents a configuration section.
type Section struct {
	Name   string
	Values []string
}

// Sections dumps the configuration to a slice of sections.
// Sensible fields are masked.
// This is useful for logging.
func Sections(cfg *Config) []Section {
	sections := make([]Section, 0)

	sectionType := reflect.TypeOf(*cfg)
	sectionValue := reflect.ValueOf(*cfg)

	for i := range sectionType.NumField() {
		sectionName := sectionType.Field(i).Name
		fieldValue := sectionValues(sectionValue.Field(i).Interface())

		sections = append(sections, Section{
			Name:   sectionName,
			Values: fieldValue,
		})
	}

	return sections
}

func sectionValues(section any) []string {
	t := reflect.TypeOf(section)
	v := reflect.ValueOf(section)

	values := make([]string, 0, t.NumField())

	for i := range t.NumField() {
		fieldName := t.Field(i).Name
		fieldValue := v.Field(i).Interface()

		if isProtectedField(fieldName) {
			fieldValue = "****"
		}

		values = append(values, fmt.Sprintf("%s=%v", fieldName, fieldValue))
	}

	return values
}

func isProtectedField(fieldName string) bool {
	return strings.Contains(fieldName, "Password") || strings.Contains(fieldName, "Secret")
}
