import tensorflow as tf

from .killerinstinct import KillerInstinct


def main():
    ki = KillerInstinct()

    # Execute dummy invocation of learn and predict graphs
    ki.learn()
    ki.predict()

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
