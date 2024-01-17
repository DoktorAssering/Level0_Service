#ВНИМАНИЕ НЕ РАБОТАЕТ, ЖДИТЕ ОБНОВЛЕНИЯ

"""
import pytest
import requests
import json
from faker import Faker

BASE_URL = "http://127.0.0.1:8080"

@pytest.fixture
def server_is_running():
    response = requests.get(BASE_URL)
    print(response.text)
    assert response.status_code == 200

def test_add_data_with_fake_data(server_is_running):
    fake = Faker()
    fake_json_data = {
        "name": fake.name(),
        "email": fake.email(),
        "phone_number": fake.phone_number(),
        "address": fake.address(),
    }

    response = requests.post(BASE_URL + "/add-json", data={"jsonData": json.dumps(fake_json_data)})
    assert response.status_code == 200
    assert "New data added" in response.text
    assert "ID" in response.text

def test_get_info(server_is_running):
    response = requests.get(BASE_URL + "/get-json?number=1")
    assert response.status_code == 200
    assert "Info" in response.text
    assert "Number" in response.text

def test_get_info_invalid_number(server_is_running):
    response = requests.get(BASE_URL + "/get-json?number=invalid")
    assert response.status_code == 400
    assert "Invalid number" in response.text

def test_add_data(server_is_running):
    json_data = '{"key": "value"}'
    response = requests.post(BASE_URL + "/add-json", data={"jsonData": json_data})
    assert response.status_code == 200
    assert "New data added" in response.text
    assert "ID" in response.text

def test_add_data_missing_json_data(server_is_running):
    response = requests.post(BASE_URL + "/add-json")
    assert response.status_code == 500
    assert "Error adding data" in response.text

def test_update_data(server_is_running):
    json_data = '{"key": "updated_value"}'
    response = requests.put(BASE_URL + "/update-json?number=1", data={"jsonData": json_data})
    assert response.status_code == 200
    assert "Data updated successfully" in response.text

def test_update_data_invalid_number(server_is_running):
    json_data = '{"key": "updated_value"}'
    response = requests.put(BASE_URL + "/update-json?number=invalid", data={"jsonData": json_data})
    assert response.status_code == 400
    assert "Invalid number" in response.text

def test_delete_data(server_is_running):
    response = requests.delete(BASE_URL + "/delete-json?number=1")
    assert response.status_code == 200
    assert "Data deleted successfully" in response.text

def test_delete_data_invalid_number(server_is_running):
    response = requests.delete(BASE_URL + "/delete-json?number=invalid")
    assert response.status_code == 400
    assert "Invalid number" in response.text

def test_server_shutdown(server_is_running):
    response = requests.post(BASE_URL + "/shutdown")
    assert response.status_code == 200
    assert "Shutting down the server" in response.text

    response = requests.get(BASE_URL)
    assert response.status_code == 200
    assert "Server is offline" in response.text

def test_backup_and_restore(server_is_running):
    response = requests.post(BASE_URL + "/backup")
    assert response.status_code == 200
    assert "Successful creation of a backup" in response.text

    response = requests.post(BASE_URL + "/restore")
    assert response.status_code == 200
    assert "Restored cache successfully" in response.text

def test_restore_invalid_backup(server_is_running):
    response = requests.post(BASE_URL + "/restore?file=invalid_backup.json")
    assert response.status_code == 500
    assert "Error restoring cache from the database" in response.text

def test_server_status(server_is_running):
    response = requests.get(BASE_URL + "/status")
    assert response.status_code == 200
    assert "Server online" in response.text

def test_invalid_endpoint(server_is_running):
    response = requests.get(BASE_URL + "/invalid-endpoint")
    assert response.status_code == 404
    assert "404 Not Found" in response.text

def test_long_running_process(server_is_running):
    response = requests.post(BASE_URL + "/long-running-process")
    assert response.status_code == 200
    assert "Process completed successfully" in response.text

def test_long_running_process_timeout(server_is_running):
    response = requests.post(BASE_URL + "/long-running-process?timeout=1")
    assert response.status_code == 500
    assert "Process timed out" in response.text
"""