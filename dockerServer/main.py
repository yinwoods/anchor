from flask import Flask, Blueprint
from container import container_bp
from image import image_bp
from query import query_bp


app = Flask(__name__)

app.register_blueprint(container_bp, url_prefix="/api/container")
app.register_blueprint(image_bp, url_prefix="/api/image")
app.register_blueprint(query_bp, url_prefix="/api/query")
app.run(host="0.0.0.0", port="8089", debug=True)
