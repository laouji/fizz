package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/laouji/fizz/pkg/cache"
	"github.com/sirupsen/logrus"
)

type StatsOutput struct {
	Query    url.Values `json:"query"`
	HitCount float64    `json:"hit_count"`
}

func Stats(
	cache *cache.Client,
	logger logrus.FieldLogger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		results, err := cache.GetEndpointHitCount(r.Context(), 0, 1)
		if err != nil {
			logger.WithError(err).Error("failed to get endpoint hit count")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		out := []StatsOutput{}
		for _, z := range results {
			rawQuery, ok := z.Member.(string)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Errorf("failed convert %T to string", z.Member)
				return
			}
			query, err := url.ParseQuery(rawQuery)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.WithError(err).Errorf("failed to parse query %s", rawQuery)
				return
			}
			out = append(out, StatsOutput{
				Query:    query,
				HitCount: z.Score,
			})
		}

		if err := json.NewEncoder(w).Encode(out); err != nil {
			logger.WithError(err).Error("failed to encode stats response JSON")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
}
