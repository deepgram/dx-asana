package ui

import (
	"encoding/json"
	"fmt"

	"github.com/deepgram/dx-asana/internal/asana"
)

// JSONOutput is the envelope every --json command writes to stdout.
type JSONOutput struct {
	Success bool                   `json:"success"`
	Data    interface{}            `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Errors  []errorDetail          `json:"errors,omitempty"`
	Status  int                    `json:"status,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

type errorDetail struct {
	Message string `json:"message"`
	Help    string `json:"help,omitempty"`
}

// PrintJSON writes a single-resource or list response. err may be nil.
func PrintJSON(data interface{}, err error) {
	PrintJSONWithMeta(data, nil, err)
}

// PrintJSONWithMeta writes the response envelope with optional metadata
// (count, query echo, etc.) and structured error details when err is an
// asana.APIError.
func PrintJSONWithMeta(data interface{}, meta map[string]interface{}, err error) {
	output := JSONOutput{Success: err == nil, Meta: meta}

	if err != nil {
		output.Error = err.Error()
		if apiErr, ok := asana.AsAPIError(err); ok {
			output.Status = apiErr.StatusCode
			if len(apiErr.Errors) > 0 {
				output.Errors = make([]errorDetail, len(apiErr.Errors))
				for i, e := range apiErr.Errors {
					output.Errors[i] = errorDetail{Message: e.Message, Help: e.Help}
				}
			}
		}
	} else {
		output.Data = data
	}

	jsonBytes, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(jsonBytes))
}
