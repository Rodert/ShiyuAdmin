"""Verify only an isolated CI SQLite deployment; never point this at user data."""
import json
import os
import pathlib
import subprocess
import sys
import tempfile
import time
import urllib.error
import urllib.request

BASE = "http://127.0.0.1:18000/api/v1/system"
STATE = pathlib.Path(tempfile.gettempdir()) / "shiyu-container-smoke.json"


def request(path, method="GET", token=None, data=None, content_type="application/json", expected=200):
    headers = {"Content-Type": content_type}
    if token:
        headers["Authorization"] = "Bearer " + token
    body = json.dumps(data).encode() if isinstance(data, dict) else data
    req = urllib.request.Request(BASE + path, data=body, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=10) as res:
            status, result = res.status, res.read()
    except urllib.error.HTTPError as error:
        status, result = error.code, error.read()
    assert status == expected, f"{method} {path}: expected {expected}, got {status}"
    return result


def ready():
    for _ in range(60):
        try:
            request("/health")
            return
        except (OSError, AssertionError):
            time.sleep(1)
    raise RuntimeError("Application did not become healthy")


ready()
if sys.argv[1] == "prepare":
    status = json.loads(request("/installation"))["data"]
    assert not status["installed"], "Smoke test requires a fresh disposable deployment"
    assert all(status["checks"].values())
    setup = subprocess.check_output([
        "docker", "compose", "-f", "docker-compose.sqlite.yml", "exec", "-T",
        "shiyu-app", "cat", "/app/data/setup-token",
    ], text=True).strip()
    credentials = {"username": "ci-owner", "password": "Shiyu!Container2026"}
    request("/installation", "POST", data={**credentials, "token": setup, "name": "Container CI"})
    token = json.loads(request("/auth/login", "POST", data=credentials))["data"]["token"]
    body = b'--shiyu-boundary\r\nContent-Disposition: form-data; name="file"; filename="proof.txt"\r\nContent-Type: text/plain\r\n\r\ncontainer-persistent-proof\r\n--shiyu-boundary--\r\n'
    uploaded = json.loads(request("/files/upload", "POST", token, body, "multipart/form-data; boundary=shiyu-boundary"))["data"]
    fd = os.open(STATE, os.O_WRONLY | os.O_CREAT | os.O_TRUNC, 0o600)
    with os.fdopen(fd, "w") as out:
        json.dump({"token": token, "code": uploaded["file_code"]}, out)
    print("Installation and upload passed")
elif sys.argv[1] == "verify":
    state = json.loads(STATE.read_text())
    token, code = state["token"], state["code"]
    request("/profile", token=token)
    assert json.loads(request("/installation"))["data"]["installed"]
    assert request(f"/files/{code}/download", token=token) == b"container-persistent-proof"
    request("/auth/logout", "POST", token)
    request("/profile", token=token, expected=401)
    STATE.unlink()
    print("Container recreation preserved installation, signing key and file; logout passed")
else:
    raise ValueError("Expected prepare or verify")
