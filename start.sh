#!/bin/bash
# Get the LAN IP address
IP=$(hostname -I | awk '{print $1}')
echo "========================================================"
echo "Starting File Transfer Server..."
echo "On your iPhone, open Safari and go to: http://$IP:5000"
echo "Files uploaded will appear in the 'uploads' folder."
echo "Press Ctrl+C to stop."
echo "========================================================"

# Activate venv
source venv/bin/activate

# Run server with Gunicorn
# Timeout 3600 seconds (1 hour) to ensure 700MB+ files can upload even on slow WiFi
# Workers=4 for better concurrency
gunicorn -b 0.0.0.0:5000 --timeout 3600 --workers 2 --threads 4 server:app
