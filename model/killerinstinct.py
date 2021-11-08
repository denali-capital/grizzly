import tensorflow as tf
import tensorflow.keras as k

class KillerInstinct(tf.Module):

    """
    Grizzly's Killer Instinct.
    Neural network with
        Inputs                  Output
        Price delta             Classification score [0, 1]
        Liquidity Exch1
        Liquidity Exch2
        Average Latency Exch1
        Average Latency Exch2
        Volatility Exch1
        Volatility Exch2
    Deciding for us whether to make a play
    """

    def __init__(self):
        super().__init__()

        self.num_features = 7

        self.batch_size = 32
        self._model = k.Sequential(
            [
                k.layers.Input(
                    shape=(1, self.num_features),
                    batch_size=self.batch_size
                ),
                k.layers.Dense(4),
                k.layers.Dense(1, activation='sigmoid')
            ]
        )

        self._global_step = tf.Variable(0, dtype=tf.int32, trainable=False)
        self._optimizer = k.optimizers.Adam()
        self._loss = k.losses.BinaryCrossEntropy()

    @tf.function
    def learn():
        pass

    @tf.function
    def predict():
        pass
