package config

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

type ValidationErrors []ValidationError

func (ve *ValidationErrors) Error() string {
	strErrors := []string{}
	for _, v := range *ve {
		strErrors = append(strErrors, fmt.Sprintf("field %s: %s", v.Field, v.Message))
	}

	return strings.Join(strErrors, "; ")
}
