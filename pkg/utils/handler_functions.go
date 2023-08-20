package utils

import (
	"TESTGO/config"
	"TESTGO/pkg/api/errors"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
)

// ResponseBodyToStruct is a function that convert response body to map[string]interface{}
func ResponseBodyToStruct(resp *resty.Response) (interface{}, error) {
	// response body to map[string]interface{} for map erro response body to struct
	var resultMap map[string]interface{}
	bodyBytes := resp.String()
	// if response body is empty read body from resp.RawResponse.Body (Mock Seekster API)
	if resp.String() == "" {
		bodyBytes2, err3 := io.ReadAll(resp.RawResponse.Body)
		if err3 != nil {
			return nil, err3
		}
		bodyBytes = string(bodyBytes2)
	}
	// convert response body to map[string]interface{}
	err := json.Unmarshal([]byte(bodyBytes), &resultMap)
	if err != nil {
		return nil, err
	}
	return resultMap, nil

}

func BindAndValidateInput(c *gin.Context, obj interface{}) error {
	if err := c.BindJSON(obj); err != nil {
		return errors.ErrBadRequest
	}
	// Validate input
	if err := config.Validate.Struct(obj); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			errorMessages := handleValidationErrors(err, &config.Translator)
			return errors.NewAppError(http.StatusBadRequest, errors.ErrValidationInput.Code, strings.Join(errorMessages, " , "))
		} else {
			return errors.ErrValidationInput
		}
	}
	return nil
}

func HandleAPIError(resp *resty.Response, err error) error {
	if err != nil {
		// convert response body to map[string]interface{}
		resultMap, err := ResponseBodyToStruct(resp)
		if err != nil {
			return errors.ErrParseJSON
		}
		return errors.NewExternalAPIError(resp.StatusCode(), resp.Status(), "", "", resultMap)
	}
	return nil
}
