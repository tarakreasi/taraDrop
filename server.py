from flask import Flask, render_template, request, redirect, url_for
import os

app = Flask(__name__)

# Configure upload folder
UPLOAD_FOLDER = 'uploads'
if not os.path.exists(UPLOAD_FOLDER):
    os.makedirs(UPLOAD_FOLDER)

app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER
# Set max file size to 16GB to support large videos
app.config['MAX_CONTENT_LENGTH'] = 16 * 1024 * 1024 * 1024

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/upload_chunk', methods=['POST'])
def upload_chunk():
    try:
        if 'file' not in request.files:
            return 'No file part', 400
        
        file = request.files['file']
        filename = request.form['filename']
        chunk_index = int(request.form['chunkIndex'])
        
        # Secure the filename to prevent directory traversal
        from werkzeug.utils import secure_filename
        filename = secure_filename(filename)
        
        save_path = os.path.join(app.config['UPLOAD_FOLDER'], filename)
        
        # If it's the first chunk, overwrite/create new file. Otherwise append.
        mode = 'wb' if chunk_index == 0 else 'ab'
        
        with open(save_path, mode) as f:
            f.write(file.read())
            
        print(f"Received chunk {chunk_index} for {filename}")
        return 'Chunk received', 200
        
    except Exception as e:
        print(f"Error: {e}")
        return str(e), 500

if __name__ == '__main__':
    # Listen on all interfaces
    # debug=False to prevent reloader from restarting server when a file is saved in the directory
    app.run(host='0.0.0.0', port=5000, debug=False, threaded=True)
