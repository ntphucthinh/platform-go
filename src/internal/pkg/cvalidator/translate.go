package cvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	trans ut.Translator
)

func TranslateError(err error) map[string]string {
	if err == nil {
		return map[string]string{}
	}

	// 1. Check validation errors (binding tags)
	var valErrors validator.ValidationErrors
	if errors.As(err, &valErrors) {
		errMessages := map[string]string{}
		for _, e := range valErrors {
			fieldName := e.Field()
			if fieldName == "" {
				fieldName = toLowerFirst(e.StructField())
			}
			msg := e.Translate(trans)
			if msg == "" || trans == nil {
				msg = fmt.Sprintf("%s is invalid", fieldName)
			}
			errMessages[fieldName] = msg
		}
		return errMessages
	}

	// 2. Check JSON unmarshal type error via standard error type
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		fieldName := unmarshalTypeError.Field
		if fieldName == "" {
			fieldName = "request"
		} else {
			parts := strings.Split(fieldName, ".")
			fieldName = toLowerFirst(parts[len(parts)-1])
		}
		expectedType := "a valid value"
		if unmarshalTypeError.Type != nil {
			expectedType = translateType(unmarshalTypeError.Type.String())
		}
		return map[string]string{
			fieldName: fmt.Sprintf("%s must be provided as %s", fieldName, expectedType),
		}
	}

	// 3. Fallback regex for string errors or syntax error
	field, msg := parseJSONUnmarshalError(err.Error())
	return map[string]string{
		field: msg,
	}
}

func toLowerFirst(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// parseJSONUnmarshalError parses JSON unmarshal errors and returns field name and translated message.
func parseJSONUnmarshalError(errMsg string) (string, string) {
	// Example: "json: cannot unmarshal string into Go struct field PutAdminUsersRequest.userIDs of type uint"
	pattern := `json: cannot unmarshal (\w+) into Go struct field ([A-Za-z0-9_.]+) of type (.+)`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(errMsg)

	const matchNum = 4
	// Return generic message if the error is syntax error or does not match type mismatch pattern
	if len(matches) != matchNum {
		return "error", "Invalid request format"
	}

	fullFieldPath := matches[2]
	expectedType := matches[3]

	// Extract the field name after the last dot.
	parts := strings.Split(fullFieldPath, ".")
	fieldName := toLowerFirst(parts[len(parts)-1])

	expectedTypeEN := translateType(expectedType)

	return fieldName, fmt.Sprintf(
		"%s must be provided as %s",
		fieldName,
		expectedTypeEN,
	)
}

func translateType(typeName string) string {
	// 1. Remove pointer prefix (e.g. *string -> string, *int -> int)
	typeName = strings.TrimPrefix(typeName, "*")

	// 2. Identify array or slice (e.g. []string, []api.CreateHomestayRoomInput)
	if strings.HasPrefix(typeName, "[]") {
		return "an array"
	}

	// 3. Identify package-qualified structs/types
	if strings.Contains(typeName, ".") {
		switch {
		case strings.HasSuffix(typeName, "Time"):
			return "a valid datetime (ISO 8601 / RFC 3339)"
		case strings.HasSuffix(typeName, "Date"):
			return "a valid date (YYYY-MM-DD)"
		case strings.HasSuffix(typeName, "Email"):
			return "a valid email"
		case strings.HasSuffix(typeName, "UUID"):
			return "a valid UUID"
		default:
			return "an object"
		}
	}

	// 4. Primitive type mapping table
	typeTranslations := map[string]string{
		"string":  "a string",
		"int":     "an integer",
		"int8":    "an integer",
		"int16":   "an integer",
		"int32":   "an integer",
		"int64":   "an integer",
		"uint":    "an integer",
		"uint8":   "an integer",
		"uint16":  "an integer",
		"uint32":  "an integer",
		"uint64":  "an integer",
		"float":   "a number",
		"float64": "a number",
		"float32": "a number",
		"bool":    "a boolean",
		"array":   "an array",
		"slice":   "an array",
		"map":     "an object",
		"object":  "an object",
		"null":    "a null value",
	}

	if translation, exists := typeTranslations[typeName]; exists {
		return translation
	}

	return "a value of type " + typeName
}
