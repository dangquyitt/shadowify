#!/usr/bin/env python3
import os
import time
import logging
import argparse
import yaml
import grpc
from concurrent import futures

import extractor_pb2
import extractor_pb2_grpc
from youtube_extractor import YoutubeExtractor

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class ExtractorServicer(extractor_pb2_grpc.ExtractorServiceServicer):
    def __init__(self, config):
        self.extractor = YoutubeExtractor(config)
        logger.info("Initialized ExtractorServicer")

    def ExtractYoutubeMetadata(self, request, context):
        try:
            logger.info(
                f"Received metadata request for YouTube ID: {request.youtube_id}")
            metadata = self.extractor.extract_metadata(request.youtube_id)

            response = extractor_pb2.YoutubeMetadataResponse(
                id=metadata.get('youtube_id', ''),
                title=metadata.get('title', ''),
                full_title=metadata.get('full_title', ''),
                description=metadata.get('description', ''),
                duration=metadata.get('duration', 0),
                duration_string=metadata.get('duration_string', ''),
                thumbnail=metadata.get('thumbnail', ''),
                tags=metadata.get('tags', []),
                categories=metadata.get('categories', [])
            )

            logger.info(
                f"Successfully extracted metadata for {request.youtube_id}")
            return response

        except Exception as e:
            logger.error(f"Error extracting metadata: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error extracting metadata: {str(e)}")
            return extractor_pb2.YoutubeMetadataResponse()

    def ExtractYoutubeTranscript(self, request, context):
        try:
            logger.info(
                f"Received transcript request for YouTube ID: {request.youtube_id}")
            language_code, segments = self.extractor.extract_transcript(
                request.youtube_id)

            transcript_segments = []
            for segment in segments:
                transcript_segments.append(
                    extractor_pb2.TranscriptSegment(
                        start=segment['start'],
                        end=segment['end'],
                        text=segment['text']
                    )
                )

            response = extractor_pb2.YoutubeTranscriptResponse(
                language_code=language_code,
                segments=transcript_segments
            )

            logger.info(
                f"Successfully extracted transcript for {request.youtube_id}")
            return response

        except Exception as e:
            logger.error(f"Error extracting transcript: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error extracting transcript: {str(e)}")
            return extractor_pb2.YoutubeTranscriptResponse()


def load_config(config_path):
    with open(config_path, 'r') as f:
        return yaml.safe_load(f)


def serve(port, config):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    extractor_pb2_grpc.add_ExtractorServiceServicer_to_server(
        ExtractorServicer(config), server
    )
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    logger.info(f"Extractor service started on port {port}")

    try:
        while True:
            time.sleep(86400)  # One day in seconds
    except KeyboardInterrupt:
        server.stop(0)
        logger.info("Server stopped")


if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        description='Run the Extractor gRPC server')
    parser.add_argument('--port', type=int, default=50051,
                        help='The server port')
    parser.add_argument('--config', type=str,
                        default='config.yaml', help='Path to config file')
    args = parser.parse_args()

    config = load_config(args.config)
    serve(args.port, config)
