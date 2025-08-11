#!/bin/bash

if command -v minio >/dev/null 2>&1; then
    echo "MinIO is already installed at $(command -v minio)"
else
    if [[ "$1" == "--linux" ]]; then
        curl -O https://dl.min.io/server/minio/release/linux-amd64/minio
    elif [[ "$1" == "--macos" ]]; then
        curl -O https://dl.min.io/server/minio/release/darwin-amd64/minio
    else
        echo "Usage: $0 --linux|--macos"
        exit 1
    fi

    chmod +x minio
    sudo mv minio /usr/local/bin/
    echo "MinIO installed successfully. You can run it using the command 'minio server /path/to/data'."
fi

echo "Make sure to replace '/path/to/data' with your desired data directory."
echo "For more information, visit https://docs.min.io/docs/minio-quickstart-guide.html"
echo "To start MinIO, use the command: minio server /path/to/data"
echo "You can also set up MinIO as a service for easier management."
echo "For service setup, refer to the MinIO documentation."
echo "To configure MinIO, you can set environment variables like MINIO_ROOT_USER and MINIO_ROOT_PASSWORD."
echo "Example: export MINIO_ROOT_USER=yourusername"
echo "Example: export MINIO_ROOT_PASSWORD=yourpassword"
echo "You can also create a systemd service file for MinIO to run it as a service."
echo "Example systemd service file:"
echo "[Unit]
Description=MinIO
After=network.target   
[Service]
User=minio-user
Group=minio-user
ExecStart=/usr/local/bin/minio server /path/to/data
Restart=always
EnvironmentFile=-/etc/default/minio
[Install]
WantedBy=multi-user.target"
echo "=========================="

if [[ "$2" == "--port" && -n "$3" && -n "$4" ]]; then
    if lsof -i ":$3" >/dev/null 2>&1; then
        echo "Error: Port $3 is already in use. Please choose a different port."
        exit 1
    fi
    sudo mkdir -p /var/minio/data
    sudo chmod 777 /var/minio/data
    minio server --console-address ":$3" /var/minio/data --address ":$4"
    echo "MinIO server started on port $3. Access the web console at http://localhost:$3/minio"
else
    exit 0
fi
