#!/bin/bash

# Create Python virtual environment
python3 -m venv .venv

# Activate virtual environment
source venv/bin/activate

# Install dependencies
pip install --upgrade pip
pip install -r requirements.txt

# Generate gRPC code from proto file
python -m grpc_tools.protoc -I../../proto --python_out=. --grpc_python_out=. ../../proto/extractor.proto

echo "Setup completed successfully!"
