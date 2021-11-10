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

        # set this to number of asset pairs we are trading?
        self.batch_size = 32
        # how to do batch_size while allowing inference with only 1 at a time
        self._model = k.Sequential(
            [
                k.layers.Input(
                    shape=(self.num_features,),
                    batch_size=self.batch_size
                ),
                k.layers.Dense(4, activation='relu'),
                k.layers.Dense(1, activation='sigmoid')
            ]
        )

        self._global_step = tf.Variable(0, dtype=tf.int32, trainable=False)
        self._optimizer = k.optimizers.Adam()
        self._loss = k.losses.BinaryCrossentropy()

    @tf.function(
        input_signature=[
            tf.TensorSpec(shape=(None, 7), dtype=tf.float32),
            tf.TensorSpec(shape=(None,), dtype=tf.int32)
        ]
    )
    def learn(self, data, labels):
        self._global_step.assign_add(1)
        with tf.GradientTape() as tape:
            loss = self._loss(labels, self._model(data))
            tf.print(self._global_step, ": loss: ", loss)

        gradient = tape.gradient(loss, self._model.trainable_variables)
        self._optimizer.apply_gradients(zip(gradient, self._model.trainable_variables))
        return {
            'loss': loss
        }

    @tf.function(
        input_signature=[
            tf.TensorSpec(shape=(None, 7), dtype=tf.float32),
        ]
    )
    def predict(self, data):
        predictions = self._model(data)
        # deliberately return probability so it can be used in calculation of how much to buy/sell
        return {
            'predictions': predictions
        }
