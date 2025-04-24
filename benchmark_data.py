import pandas
import pandas as pd
import numpy as np
import os
import json

test_results_folder = "testresults/testresults_10"
main_results_file = test_results_folder + "/algorithm_results.csv"

with open(main_results_file, 'r') as f:
    #load in pandas
    lines = f.readlines()

# read data in pd.DataFrame
data = []
for line in lines[1:]:
    data.append(line.split(","))
data = pandas.read_csv(main_results_file)

#get average of each algorithm
algorithms = data['algorithm'].unique()
averages = {}
for algorithm in algorithms:
    algorithm_data = data[data['algorithm'] == algorithm]
    query_time_ms = algorithm_data['query_time_ms']
    response_time_ms = algorithm_data['response_time_ms']
    total_time_ms = algorithm_data['total_time_ms']

    averages[algorithm] = {
        'query_time_ms': query_time_ms.mean(),
        'response_time_ms': response_time_ms.mean(),
        'total_time_ms': total_time_ms.mean()
    }
    # save to json
    with open(test_results_folder + "/averages.json", 'w') as a:
        json.dump(averages, a, indent=4)
    # save to csv
    averages_df = pd.DataFrame(averages).T
    averages_df.to_csv(test_results_folder + "/averages.csv", index=True)


