from flask import Flask, Blueprint
from config import config
from container import container_bp


app = Flask(__name__)

app.register_blueprint(container_bp, url_prefix="/api/container")
app.run(host="0.0.0.0", port="8089", debug=True)
