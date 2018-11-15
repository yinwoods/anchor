import random
import requests
from flask import Blueprint, jsonify, request
from config import servers


image_bp = Blueprint('image', __name__)


def search(mid):
    for server in servers:
        resp = requests.get(f"{server}/images/json")
        images = resp.json()
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
