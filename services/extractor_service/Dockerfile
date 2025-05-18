FROM python:3.10-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    ffmpeg \
    git \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy proto file and generate gRPC code
COPY ./proto/extractor.proto /app/proto/
RUN python -m grpc_tools.protoc -I/app/proto --python_out=/app --grpc_python_out=/app /app/proto/extractor.proto

# Copy application code
COPY server.py youtube_extractor.py config.yaml ./

# Default port
EXPOSE 50051

CMD ["python", "server.py"]
