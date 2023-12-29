import requests


class AuthClient:
    def __init__(self, server_url):
        self.server_url = server_url
        self.token = None

    def login(self, username, password) -> bool:
        # Prepare the payload for the login request
        payload = {'username': username, 'password': password}

        # Make the login request
        response = requests.post(self.server_url, data=payload)

        # Check if the request was successful (status code 202)
        approved = response.status_code == 200
        if not approved:
            print(f"Login failed. Status code: {response.status_code}, Error: {response.text}")

        return approved
