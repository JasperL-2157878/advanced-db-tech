import random
import requests, time, json, csv

URL = "http://localhost:8080/api/v1/route/"
MIN = 10560298646904
MAX = 10560299999998
WAIT = False
ALGORITHMS = [
    "alg/dijkstra", "alg/astar",
    "alg/bddijkstra", "alg/bdastar",
    "opt/none", "opt/tnr",
    "opt/ch", "opt/chtnr",
]


def get_files_and_writers():
    fieldnames = ['iteration', 'algorithm', 'query_time_ms', 'response_time_ms', 'total_time_ms']
    files = {
        'alg/dijkstra': open(f'testresults/alg_dijkstra.csv', 'w', newline=''),
        'alg/astar': open(f'testresults/alg_astar.csv', 'w', newline=''),
        'alg/bddijkstra': open(f'testresults/alg_bddijkstra.csv', 'w', newline=''),
        'alg/bdastar': open(f'testresults/alg_bdastar.csv', 'w', newline=''),
        'opt/none': open(f'testresults/opt_none.csv', 'w', newline=''),
        'opt/tnr': open(f'testresults/opt_tnr.csv', 'w', newline=''),
        'opt/ch': open(f'testresults/opt_ch.csv', 'w', newline=''),
        'opt/chtnr': open(f'testresults/opt_chtnr.csv', 'w', newline='')
    }
    
    writers = {
        'alg/dijkstra': csv.DictWriter(files['alg/dijkstra'], fieldnames=fieldnames),
        'alg/astar': csv.DictWriter(files['alg/astar'], fieldnames=fieldnames),
        'alg/bddijkstra': csv.DictWriter(files['alg/bddijkstra'], fieldnames=fieldnames),
        'alg/bdastar': csv.DictWriter(files['alg/bdastar'], fieldnames=fieldnames),
        'opt/none': csv.DictWriter(files['opt/none'], fieldnames=fieldnames),
        'opt/tnr': csv.DictWriter(files['opt/tnr'], fieldnames=fieldnames),
        'opt/ch': csv.DictWriter(files['opt/ch'], fieldnames=fieldnames),
        'opt/chtnr': csv.DictWriter(files['opt/chtnr'], fieldnames=fieldnames)
    }

    for _, writer in writers.items():
        writer.writeheader()

    return files, writers


def get_random_nodes(min, max):
    source = random.randint(min, max)
    target = random.randint(min, max)
    while target == source:
        target = random.randint(min, max)

    return source, target


files, writers = get_files_and_writers()
for test in range(100):
    print(f"{test}/100")

    source, target = get_random_nodes(MIN, MAX)
    response = requests.get(URL+"opt/ch", params={
        'from': source,
        'to': target
    })

    while response.status_code != 200:
        source, target = get_random_nodes(MIN, MAX)
        response = requests.get(URL+"opt/ch", params={
            'from': source,
            'to': target
        })

    for alg, writer in writers.items():
        n = 3
        i = n

        while i > 0:
            response = requests.get(URL+alg, params={
                'from': source,
                'to': target
            })

            if response.status_code == 200:
                data = response.json()
                query_time_ms = data["query_time_ns"] / 1_000_000
                response_time_ms = data["response_time_ns"] / 1_000_000
                total_time_ms = query_time_ms + response_time_ms

                writer.writerow({
                    'iteration': n - i + 1,
                    'algorithm': alg,
                    'query_time_ms': query_time_ms,
                    'response_time_ms': response_time_ms,
                    'total_time_ms': total_time_ms
                })
                files[alg].flush()

                i -= 1
