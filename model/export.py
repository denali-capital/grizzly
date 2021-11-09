import tensorflow as tf

from .killerinstinct import KillerInstinct


def main():
    ki = KillerInstinct()

    # Execute dummy invocation of learn and predict graphs
    ki.learn(
        tf.zeros((ki.batch_size, 1, ki.num_features), dtype=tf.float32),
        tf.zeros((ki.batch_size), dtype=tf.int32)
    )
    ki.predict(tf.zeros((ki.batch_size, 1, ki.num_features), dtype=tf.float32))

    tf.saved_model.save(
        ki,
        "KillerInstinct",
        signatures={
            "learn": ki.learn,
            "predict": ki.predict
        }
    )


if __name__ == "__main__":
    main()
