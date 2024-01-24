import requests


class PrometheusInterface:
    endpoint = "http://127.0.0.1:9090"

    def __init__(self, endpoint):
        self.endpoint = endpoint

    def get_prometheus_instant(self, metric_name):
        # Construct the Prometheus query for instant vector query
        query = f'{metric_name}'  # Adjust the query as needed

        # Make an HTTP request to Prometheus
        try:
            response = requests.get(f'{self.endpoint}/api/v1/query', params={'query': query})
            response.raise_for_status()  # Raise an exception for 4xx and 5xx status codes
            result = response.json()

            if result['status'] == 'success' and result['data']['result'][0]['value']:
                return result['data']['result'][0]['value'][1]
            else:
                return None
        except requests.RequestException as e:
            print(f"Error querying Prometheus: {e}")
            return None

    def get_prometheus_vector(self, metric_name, interval):
        # Construct the Prometheus query for instant vector query
        query = f'{metric_name}[{interval}h]'  # Adjust the query as needed

        # Make an HTTP request to Prometheus
        try:
            response = requests.get(f'{self.endpoint}/api/v1/query', params={'query': query})
            response.raise_for_status()  # Raise an exception for 4xx and 5xx status codes
            result = response.json()

            if result['status'] == 'success' and len(result['data']['result'][0]['values']):
                return [x[1] for x in result['data']['result'][0]['values']]
            else:
                return None
        except requests.RequestException as e:
            print(f"Error querying Prometheus: {e}")
            return None
