import random
import json
import requests
from flask import Blueprint, jsonify, request
from config import config


container_bp = Blueprint('container', __name__)


servers = config["test"]["servers"].split(",")


def search(cid):
    for server in servers:
        resp = requests.get(f"{server}/containers/json?all=1")
        containers = resp.json()
        for container in containers:
            container.pop("Labels")
            container.pop("Mounts")
            if container["Id"] == cid or container["Id"].startswith(cid):
                return container, server
    return f"Error, {cid} Not Found", servers[random.randint(len(servers))]


@container_bp.route('/<cid>', methods=["GET"])
def get(cid):
    container, server = search(cid)
    if not isinstance(container, dict):
        return jsonify({"message": f"{cid} not found"})
    return jsonify(container)


@container_bp.route("/", methods=["POST"])
def post():
    data = request.get_json()
    resp = requests.post(f"{servers[random.randint(len(servers))]}/containers/create", json=data)
    return jsonify(resp.json())


@container_bp.route("/<cid>", methods=["PUT"])
def update(cid):
    data = request.get_json()
    container, server = search(cid)
    resp = requests.post(f"{server}/containers/{cid}/update", json=data)
    return jsonify(resp.json())


@container_bp.route("/<cid>", methods=["DELETE"])
def delete(cid):
    container, server = search(cid)
    resp = requests.delete(f"{server}/containers/{cid}?force=true")
    if resp.ok:
        return jsonify({"message": "success"})
    return jsonify(resp.json())
