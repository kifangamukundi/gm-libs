package binders

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateBindJSONRequest(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		handleValidationErrors(c, err, "json")
		return false
	}

	return true
}

func ValidateBindFormRequest(c *gin.Context, req any) bool {
	if err := BindForm(c, req); err != nil {
		handleValidationErrors(c, err, "form")
		return false
	}

	return true
}

func ValidateBindXMLRequest(c *gin.Context, req any) bool {
	if err := BindXML(c, req); err != nil {
		handleValidationErrors(c, err, "xml")
		return false
	}

	return true
}

func handleValidationErrors(c *gin.Context, err error, tagType string) {
	response := gin.H{
		"error": "Invalid request format",
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		validationMap := make(map[string][]string)

		validationMessages := map[string]string{
			"required": "%s is required",
			"min":      "%s must have at least %s characters",
			"max":      "%s must have no more than %s characters",
			"email":    "%s must be a valid email address",
			"url":      "%s must be a valid URL",
			"len":      "%s must be exactly %s characters",
		}

		for _, fieldError := range validationErrors {
			fieldName := getFieldName(fieldError, tagType)
			message := formatValidationError(fieldError, validationMessages)
			validationMap[fieldName] = append(validationMap[fieldName], message)
		}

		response["validation_errors"] = validationMap

		c.JSON(http.StatusBadRequest, response)

		return
	}

	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		response["error"] = fmt.Sprintf("Invalid type for field %s", e.Field)
	case *json.SyntaxError:
		response["error"] = "Invalid JSON syntax"
	default:
		response["error"] = err.Error()
	}

	c.JSON(http.StatusBadRequest, response)
}

func getFieldName(err validator.FieldError, tagType string) string {
	t := reflect.TypeOf(err.Value())

	if t.Kind() != reflect.Struct {
		return err.StructField()
	}

	field, _ := t.FieldByName(err.StructField())

	switch tagType {
	case "json":
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			return jsonTag
		}
	case "form":
		if formTag := field.Tag.Get("form"); formTag != "" {
			return formTag
		}
	case "xml":
		if xmlTag := field.Tag.Get("xml"); xmlTag != "" {
			if parts := strings.Split(xmlTag, ","); len(parts) > 0 {
				return parts[0]
			}

			return xmlTag
		}
	}

	return strings.ToLower(err.StructField())
}

func formatValidationError(err validator.FieldError, messages map[string]string) string {
	message := messages[err.Tag()]
	if message == "" {
		message = "Field '%s' is invalid"
	}

	if err.Param() != "" {
		return fmt.Sprintf(message, err.Field(), err.Param())
	}

	return fmt.Sprintf(message, err.Field())
}
