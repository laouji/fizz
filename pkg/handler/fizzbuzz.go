package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func FizzBuzz(
	logger logrus.FieldLogger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		out := make([]string, 0)

		// TODO: fizzbuzz logic

		err := json.NewEncoder(w).Encode(out)
		if err != nil {
			logger.WithError(err).Error("Failed to encode repos response JSON")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
}
