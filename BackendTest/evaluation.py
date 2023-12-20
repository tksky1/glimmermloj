import time
import pickle
import numpy as np
import sys

def model_test(route):
    sys.path.append(route)
    from test_scripts import test
    
    with open('test_dataset_5.pickle', 'rb') as file:
        test_dict = pickle.load(file, encoding='iso-8859-1')
    test_all_x = test_dict['x']
    test_all_labels = test_dict['y']
    temp = np.column_stack((test_all_x, test_all_labels))
    index = np.random.choice(3000, 300, replace=False)
    test_all = temp[index,:]
    test_x = test_all[:,:-1]
    test_y = test_all[:,-1]
    start_dt = time.time()
    result = test(test_x)
    end_dt = time.time()
    precision = round(np.sum(result == test_y)/300 ,2)
    speed = round(end_dt - start_dt, 2)
    return str(precision) + ' ' + str(speed)
