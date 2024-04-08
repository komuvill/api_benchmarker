# dummy REST server for benchmarking GET/POST/PUT/DELETE

from flask import Flask, jsonify, request

app = Flask(__name__)

@app.route('/posts', methods=['GET', 'POST'])
def posts():
    if request.method == 'GET':
        return jsonify([{"id": 1, "title": "Hello World!"}])
    elif request.method == 'POST':
        data = request.json
        return jsonify(data), 201

@app.route('/posts/<int:post_id>', methods=['GET', 'PUT', 'DELETE'])
def post(post_id):
    if request.method == 'GET':
        return jsonify({"id": post_id, "title": "Hello World!"})
    elif request.method == 'PUT':
        data = request.json
        return jsonify(data)
    elif request.method == 'DELETE':
        return '', 204

if __name__ == '__main__':
    app.run(debug=True)
