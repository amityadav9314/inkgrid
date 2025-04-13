// Package utils contains all the generic helper methods
// string.go contains all the string specific methods
package utils

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

// StringInSlice Returns true if string is present in the given slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IntegerInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Difference Set Difference: A - B
func Difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func StringCheckEmptyIfNil(a *string) string {
	if a == nil {
		return ""
	}
	return *a
}

func StringToSlice(str string, delimiter string) []string {
	return strings.Split(strings.ReplaceAll(str, " ", ""), delimiter)
}

// RemoveNonASCIICharacters removes non-ASCII characters from the string
func RemoveNonASCIICharacters(input string) string {
	var builder strings.Builder

	for _, char := range input {
		if char < 128 {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}

func ProtoToJson(message protoreflect.ProtoMessage) string {
	jsonString := protojson.Format(message)
	return jsonString
}

func GetPanicMsg(msg any) string {
	if err, ok := msg.(error); ok {
		fmt.Printf("Error msg from panic: %v\n", err)
		msg := err.Error()
		return msg
	} else {
		return ""
	}
}

func GetLastPart(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
