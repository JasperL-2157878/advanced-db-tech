import requests, time, json, csv

get_route_url = "http://localhost:8080/api/v1/route/"

popular_nodes = {
    'erkestraat': 10560298930511,
    'ieperleestraat': 10560299626068,
}

start_node = popular_nodes['erkestraat']
end_node = popular_nodes['ieperleestraat']

algorithms = [
    "alg/dijkstra", "alg/astar",
    "alg/bddijkstra", "alg/bdastar",
    "opt/none", "opt/tnr",
    "opt/ch", "opt/chtnr",
]

params = {
    "from": start_node,
    "to": end_node
}

wait_10_secs = False
# Prepare CSV file
for alg in algorithms:
    safe_alg = alg.replace("/", "_")
    with open(f'testresults/algorithm_results_{safe_alg}.csv', 'w', newline='') as csvfile:
        fieldnames = ['iteration', 'algorithm', 'query_time_ms', 'response_time_ms', 'total_time_ms']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()

        n = 1
        for i in range(n):
            print(f"[{alg}] Iteration {i + 1} of {n}")
            url = get_route_url + alg
            if wait_10_secs:
                time.sleep(10)
            response = requests.get(url, params=params)

            if response.status_code == 200:
                data = response.json()
                query_time_ms = data["query_time_ns"] / 1_000_000
                response_time_ms = data["response_time_ns"] / 1_000_000
                total_time_ms = query_time_ms + response_time_ms
            else:
                query_time_ms = None
                response_time_ms = None
                print(f"Error on {alg} run {i}: {response.status_code}")

            writer.writerow({
                'iteration': i + 1,
                'algorithm': alg,
                'query_time_ms': query_time_ms,
                'response_time_ms': response_time_ms,
                'total_time_ms': total_time_ms
            })

#merge all csv files from testresults into one
for alg in algorithms:
    safe_alg = alg.replace("/", "_")
    with open(f'testresults/algorithm_results_{safe_alg}.csv', 'r') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            with open('testresults/algorithm_results.csv', 'a', newline='') as outfile:
                writer = csv.DictWriter(outfile, fieldnames=reader.fieldnames)
                writer.writerow(row)
