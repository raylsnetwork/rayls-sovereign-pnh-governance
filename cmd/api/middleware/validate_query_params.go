package middleware

import (
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/utils"
)

// paramInfo holds metadata about a query parameter
type paramInfo struct {
	name        string
	allowedEnum []string // empty if no enum
}

// getParamInfoFromStruct extracts parameter names and their enum constraints from struct tags
func getParamInfoFromStruct(v any) map[string]paramInfo {
	params := make(map[string]paramInfo)
	t := reflect.TypeOf(v)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		formTag := field.Tag.Get("form")
		if formTag == "" || formTag == "-" {
			continue
		}

		// Extract parameter name from form tag
		parts := strings.Split(formTag, ",")
		paramName := parts[0]

		info := paramInfo{name: paramName}

		// Extract enum values if present
		enumTag := field.Tag.Get("enums")
		if enumTag != "" {
			info.allowedEnum = strings.Split(enumTag, ",")
		}

		params[paramName] = info
	}

	return params
}

// ValidateQueryParams creates a middleware that validates query parameters against struct form tags
func ValidateQueryParams(filterStruct any) gin.HandlerFunc {
	paramInfoMap := getParamInfoFromStruct(filterStruct)

	// Build list of allowed param names
	allowedParams := make([]string, 0, len(paramInfoMap))
	for name := range paramInfoMap {
		allowedParams = append(allowedParams, name)
	}

	return func(c *gin.Context) {
		queryParams := c.Request.URL.Query()

		var unknownParams []string
		var emptyParams []string

		for param, values := range queryParams {
			// Check if parameter is allowed
			info, exists := paramInfoMap[param]
			if !exists {
				unknownParams = append(unknownParams, param)
				continue
			}

			// Check if parameter has an empty value
			if slices.Contains(values, "") {
				emptyParams = append(emptyParams, param)
				continue
			}

			// Check enum values if the parameter has enum constraints
			if len(info.allowedEnum) > 0 {
				for _, value := range values {
					if !slices.Contains(info.allowedEnum, value) {
						errResponse := utils.NewInvalidEnumValueError(
							param,
							value,
							strings.Join(info.allowedEnum, ", "),
						)
						c.JSON(http.StatusBadRequest, errResponse)
						c.Abort()
						return
					}
				}
			}
		}

		if len(unknownParams) > 0 {
			errResponse := utils.NewUnknownParamsError(
				strings.Join(unknownParams, ", "),
				strings.Join(allowedParams, ", "),
			)
			c.JSON(http.StatusBadRequest, errResponse)
			c.Abort()
			return
		}

		if len(emptyParams) > 0 {
			errResponse := utils.NewEmptyParamsError(strings.Join(emptyParams, ", "))
			c.JSON(http.StatusBadRequest, errResponse)
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateQueryEncoding creates a middleware that checks if the query string is properly URL-encoded
func ValidateQueryEncoding() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawQuery := c.Request.URL.RawQuery
		if rawQuery == "" {
			c.Next()
			return
		}

		_, err := url.ParseQuery(rawQuery)
		if err != nil {
			errResponse := utils.NewMalformedQueryError(err.Error())
			c.JSON(http.StatusBadRequest, errResponse)
			c.Abort()
			return
		}

		c.Next()
	}
}
