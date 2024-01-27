package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Number string

func (n Number) Int() int {
	val, err := strconv.Atoi(string(n))
	if err != nil {
		return 0
	}
	return val
}

type ValidationError struct {
	err    error
	Msg    string `json:"msg"`
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func (e *ValidationError) Error() string {
	return e.err.Error()
}

func validationErrorResponse(
	w http.ResponseWriter,
	logger logrus.FieldLogger,
	errs validator.ValidationErrors,
) {
	res := make([]ValidationError, 0, len(errs))

	for _, vErr := range errs {
		logger.Debug(vErr)
		res = append(res, ValidationError{
			Msg:    "input validation failed",
			Field:  strings.ToLower(vErr.Field()),
			Reason: vErr.Tag(),
		})
	}

	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.WithError(err).Error("failed to encode response JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
