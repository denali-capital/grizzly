package bootstrap

import (
	"github.com/denali-capital/grizzly/types"
)

func primalInstinct(i uint, observation type.Observation, channel chan types.PredictionResponse) float32 {
	// implement logistic regression for probabilities
}

func Predict(observations []types.Observation) []float32 {
	channel := make(chan types.PredictionResponse)
	for i, observation := range observations {
		go primalInstinct(i, observation, channel)
	}

	predictions := make([]float32, len(observations))
	for i := 0; i < len(observations); i++ {
		response := <- channel
		predictions[response.Index] = response.Prediction
	}
	return predictions
}
