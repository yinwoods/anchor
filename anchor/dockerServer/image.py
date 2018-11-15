import random
import requests
from flask import Blueprint, jsonify, request
from config import servers


image_bp = Blueprint('image', __name__)


def all_images():
    server_images_map = dict()
    for server in servers:
        resp = requests.get(f"{server}/images/json?all=1")
        server_images_map[server] = resp.json()
    return server_images_map


def search(mid):
    server_images_map = all_images()
    for server, images in server_images_map.items():
        for image in images:
            if len(image["Id"]) <= 7:
                continue
            image_id = image["Id"][7:]
            if image_id == mid or image_id.startswith(mid):
                return image, server
    return f"Error, {mid} Not Found", random.choice(servers)


@image_bp.route('/<mid>', methods=["GET"])
def get(mid):
    image, server = search(mid)
    if not isinstance(image, dict):
        return jsonify({"message": f"{mid} not found"})
    return jsonify(image)


@image_bp.route("/", methods=["GET"])
def list():
    result = []
    images = all_images()
    for value in images.values():
        result += value
    return jsonify(result)


@image_bp.route("/<cid>/json", methods=["GET"])
def inspect(cid):
    image, server = search(cid)
    resp = requests.get(f"{server}/images/{cid}/json")
    return jsonify(resp.json())

# TODO
# no need
# @image_bp.route("", methods=["POST"])
# def post():
#     fromImage = request.args.get("fromImage")
#     tag = request.args.get("tag")
#     url = f"{random.choice(servers)}/images/create?fromImage={fromImage}&tag={tag}"
#     resp = requests.post(url)
#     if resp.ok:
#         return jsonify({"message": "success"})
#     return jsonify(resp.json())


@image_bp.route("/<mid>", methods=["DELETE"])
def delete(mid):
    image, server = search(mid)
    resp = requests.delete(f"{server}/images/{mid}?force=true")
    if resp.ok:
        return jsonify({"message": "success"})
    return jsonify(resp.json())
