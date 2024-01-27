package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/laouji/fizz/pkg/fizzbuzz"
	"github.com/sirupsen/logrus"
)

type FizzBuzzInput struct {
	Int1  Number `validate:"required,number,max=7"` // max refers to number of runes
	Int2  Number `validate:"required,number,max=7"` // max refers to number of runes
	Str1  string `validate:"required,printascii,max=30"`
	Str2  string `validate:"required,printascii,max=30"`
	Limit Number `validate:"required,number,max=7"` // max refers to number of runes
}

func FizzBuzz(
	logger logrus.FieldLogger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		input := extractFizzBuzzInput(r)
		v := validator.New(validator.WithRequiredStructEnabled())
		if err := v.Struct(input); err != nil {
			errs := err.(validator.ValidationErrors)
			validationErrorResponse(w, logger, errs)
			return
		}

		out, field, err := fizzbuzz.Calculate(
			input.Int1.Int(),
			input.Int2.Int(),
			input.Str1,
			input.Str2,
			input.Limit.Int(),
		)
		if err != nil {
			calculateErrorResponse(w, logger, field, err)
			return
		}

		if err := json.NewEncoder(w).Encode(out); err != nil {
			logger.WithError(err).Error("failed to encode fizzbuzz response JSON")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
}

func extractFizzBuzzInput(r *http.Request) *FizzBuzzInput {
	queryParams := r.URL.Query()
	params := &FizzBuzzInput{
		Str1:  queryParams.Get("str1"),
		Str2:  queryParams.Get("str2"),
		Int1:  Number(queryParams.Get("int1")),
		Int2:  Number(queryParams.Get("int2")),
		Limit: Number(queryParams.Get("limit")),
	}
	return params
}

func calculateErrorResponse(
	w http.ResponseWriter,
	logger logrus.FieldLogger,
	field string,
	err error,
) {
	res := []ValidationError{{
		Msg:    "input validation failed",
		Field:  strings.ToLower(field),
		Reason: err.Error(),
	}}

	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.WithError(err).Error("failed to encode response JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
