from keras.models import load_model
import tensorflow as tf
import numpy as np
import os

def test(test_x):
    X_test = tf.cast(test_x/255.0,tf.float32)
    X_test = tf.reshape(X_test, [-1, 32, 32, 3])
    model = load_model(os.path.split(os.path.realpath(__file__))[0]+'/CIFAR10_CNN_weights.h5')
    distribution = model.predict(X_test)
    result = tf.argmax(distribution, axis=1)
    return np.array(result)