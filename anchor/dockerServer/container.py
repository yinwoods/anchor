import random
import requests
from flask import Blueprint, jsonify, request
from config import servers


container_bp = Blueprint("container", __name__)


def all_containers():
    server_containers_map = dict()
    for server in servers:
        resp = requests.get(f"{server}/containers/json?all=1")
        server_containers_map[server] = resp.json()
    return server_containers_map


def search(cid):
    server_containers_map = all_containers()
    for server, containers in server_containers_map.items():
        for container in containers:
            container.pop("Labels")
            container.pop("Mounts")
            if container["Id"] == cid or container["Id"].startswith(cid):
                return container, server
    return f"Error, {cid} Not Found", random.choice(servers)


@container_bp.route("/<cid>", methods=["GET"])
def get(cid):
    container, server = search(cid)
    if not isinstance(container, dict):
        return jsonify({"message": f"{cid} not found"})
    return jsonify(container)


@container_bp.route("/", methods=["GET"])
def list():
    result = []
    containers = all_containers()
    for value in containers.values():
        result += value
    return jsonify(result)


@container_bp.route("/<cid>/json", methods=["GET"])
def inspect(cid):
    container, server = search(cid)
    resp = requests.get(f"{server}/containers/{cid}/json")
    return jsonify(resp.json())

# TODO
# no need
# @container_bp.route("/", methods=["POST"])
# def post():
#     data = request.get_json()
#     resp = requests.post(f"{random.choice(servers)}/containers/create", json=data)
#     return jsonify(resp.json())


@container_bp.route("/<cid>", methods=["POST"])
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
