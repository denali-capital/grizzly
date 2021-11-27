package nn

import (
    "log"

    "github.com/denali-capital/grizzly/types"
    tg "github.com/galeone/tfgo"
    tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type KillerInstinct struct {
    model *tg.Model
}

func NewKillerInstinct() *KillerInstinct {
    return &KillerInstinct{
        model: tg.LoadModel("KillerInstinct", []string{"serve"}, nil)
    }
}

func learn(model *tg.Model, data *tf.Tensor, labels *tf.Tensor) float32 {
    loss := model.Exec(
        []tf.Output{
            model.Op("StatefulPartitionedCall", 0),
        },
        map[tf.Output]*tf.Tensor{
            model.Op("learn_data", 0): data,
            model.Op("learn_labels", 0): labels,
        },
    )[0]

    return loss.Value().(float32)
}

func (k *KillerInstinct) Learn(observations []types.Observation) float32 {
    size := len(observations)
    var data [size][1][7]float32
    var labels [size]int32

    for i := 0; i < size; i++ {
        data[i][0] = [7]float32{
            observations[i].PriceDelta,
            observations[i].Liquidity1,
            observations[i].Liquidity2,
            observations[i].Latency1,
            observations[i].Latency2,
            observations[i].Volatility1,
            observations[i].Volatility2,
        }
        if label := observations[i].label; label == 0 || label == 1 {
            labels[i] = label
        } else {
            log.Fatalln("label must be one of {0, 1}")
        }
    }

    dataTensor, err := tf.NewTensor(data)
    if err != nil {
        log.Fatalln(err)
    }
    labelTensor, err := tf.NewTensor(labels)
    if err != nil {
        log.Fatalln(err)
    }

    return learn(k.model, dataTensor, labelTensor)
}

func predict(model *tg.Model, data *tf.Tensor) []float32 {
    rawPredictions := model.Exec(
        []tf.Output{
            model.Op("StatefulPartitionedCall_1", 0),
        },
        map[tf.Output]*tf.Tensor{
            model.Op("predict_data", 0): data,
        },
    )[0].Value().([][]float32)

    predictions := make([]float32, len(rawPredictions))
    for i := 0; i < len(rawPredictions); i++ {
        predictions[i] = rawPredictions[i][0]
    }

    return predictions
}

func (k *KillerInstinct) Predict(observations []types.Observation) []float32 {
    size := len(observations)
    var data [size][1][7]float32

    for i := 0; i < size; i++ {
        data[i][0] = [7]float32{
            observations[i].PriceDelta,
            observations[i].Liquidity1,
            observations[i].Liquidity2,
            observations[i].Latency1,
            observations[i].Latency2,
            observations[i].Volatility1,
            observations[i].Volatility2,
        }
    }

    dataTensor, err := tf.NewTensor(data)
    if err != nil {
        log.Fatalln(err)
    }

    return predict(k.model, dataTensor)
}