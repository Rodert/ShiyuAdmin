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


def data(path, method="GET", token=None, payload=None):
    return json.loads(request(path, method, token, payload))["data"]


def upload(path, token):
    body = b'--shiyu-boundary\r\nContent-Disposition: form-data; name="file"; filename="proof.txt"\r\nContent-Type: text/plain\r\n\r\ncontainer-persistent-proof\r\n--shiyu-boundary--\r\n'
    return json.loads(request(path, "POST", token, body, "multipart/form-data; boundary=shiyu-boundary"))["data"]["file_code"]


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
    code = upload("/files/upload", token)
    settings = data("/settings", token=token)
    settings.update(name="Persistent team", max_upload_mb=10, file_retention_days=14)
    data("/settings", "PUT", token, settings)
    reviewer_credentials = {"username": "ci-reviewer", "password": "Shiyu!Reviewer2026"}
    reviewer = data("/users", "POST", token, {**reviewer_credentials, "status": 1})
    reviewer_token = data("/auth/login", "POST", payload=reviewer_credentials)["token"]
    attachment = upload("/work/attachments/upload", token)
    recycled = upload("/work/attachments/upload", token)
    data(f"/work/attachments/{recycled}", "DELETE", token)
    original = data("/work/requests", "POST", token, {
        "title": "Container persistence", "content": "Keep this approval history",
        "approver_code": reviewer["user_code"], "file_codes": [attachment],
    })["id"]
    data(f"/work/requests/{original}/actions", "POST", reviewer_token, {
        "action": "rejected", "comment": "Please revise",
    })
    resubmitted = data("/work/requests", "POST", token, {
        "title": "Container persistence revised", "content": "Revision after rejection",
        "source_request_id": original, "approver_code": reviewer["user_code"], "file_codes": [attachment],
    })["id"]
    fd = os.open(STATE, os.O_WRONLY | os.O_CREAT | os.O_TRUNC, 0o600)
    with os.fdopen(fd, "w") as out:
        json.dump({"token": token, "code": code, "reviewer_token": reviewer_token,
                   "attachment": attachment, "recycled": recycled,
                   "original": original, "resubmitted": resubmitted}, out)
    print("Installation, settings, uploads, rejection and resubmission passed")
elif sys.argv[1] == "verify":
    state = json.loads(STATE.read_text())
    token, code = state["token"], state["code"]
    request("/profile", token=token)
    assert json.loads(request("/installation"))["data"]["installed"]
    assert request(f"/files/{code}/download", token=token) == b"container-persistent-proof"
    settings = data("/settings", token=token)
    assert settings["name"] == "Persistent team"
    assert settings["file_retention_days"] == 14 and settings["max_upload_mb"] == 10
    assert data("/settings/changes", token=token)["total"] == 1
    original, resubmitted = state["original"], state["resubmitted"]
    old = data(f"/work/requests/{original}", token=token)
    assert old["status"] == "rejected" and len(old["events"]) == 2
    current = data(f"/work/requests/{resubmitted}", token=state["reviewer_token"])
    assert current["source_request_id"] == original and current["status"] == "pending"
    assert current["attachments"][0]["file_code"] == state["attachment"]
    assert request(f"/work/requests/{resubmitted}/attachments/{state['attachment']}",
                   token=state["reviewer_token"]) == b"container-persistent-proof"
    filtered = data("/work/requests?view=mine&status=pending&keyword=revised", token=token)
    assert filtered["total"] == 1 and filtered["items"][0]["id"] == resubmitted
    trash = data("/work/attachments?recycled=true", token=token)
    assert trash["total"] == 1 and trash["items"][0]["file_code"] == state["recycled"]
    assert data("/files/cleanup", token=token)["total"] == 0
    data(f"/work/attachments/{state['recycled']}/restore", "POST", token)
    assert data("/work/attachments/usage", token=token)["unused"]["count"] == 1
    data(f"/work/requests/{resubmitted}/actions", "POST", state["reviewer_token"], {"action": "approved"})
    assert data("/work/notifications", token=token)["total"] == 2
    request("/auth/logout", "POST", token)
    request("/profile", token=token, expected=401)
    STATE.unlink()
    print("Container recreation preserved installation, signing key, files, settings, approval history and trash; approval and logout passed")
else:
    raise ValueError("Expected prepare or verify")
