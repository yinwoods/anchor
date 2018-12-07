import random
import requests
from flask import Blueprint, jsonify, request
from config import servers


query_bp = Blueprint("query", __name__)


@query_bp.route("/", methods=["GET"])
def get():
    args = request.args
    url=f"http://192.168.5.89:30108/api/datasources/proxy/1/query?db=k8s&q={args['q']}&epoch={args['epoch']}"
    response = requests.get(url)
    return jsonify(response.json())
