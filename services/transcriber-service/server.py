import grpc
from concurrent import futures
import time
import os

from transcriber_pb2 import Segment, TranscribeResponse, DESCRIPTOR
from transcriber_pb2_grpc import TranscriberServiceServicer, add_TranscriberServiceServicer_to_server

from audio_dowloader import AudioDownloader
from transcriber import Transcriber
from grpc_reflection.v1alpha import reflection


class TranscriberService(TranscriberServiceServicer):
    def __init__(self, downloader: AudioDownloader, transcriber: Transcriber):
        self.downloader = downloader
        self.transcriber = transcriber

    def Transcribe(self, request, context):
        video_code = request.video_code
        print(f"📥 Received video code: {video_code}")

        audio_path = ""
        try:
            audio_path = self.downloader.download(video_code)
            print(f"🎵 Audio downloaded at: {audio_path}")

            segments = self.transcriber.extract_segments(audio_path)

            response = TranscribeResponse(
                segments=[
                    Segment(
                        text=seg["text"],
                        start_time=seg["start"],
                        end_time=seg["end"]
                    ) for seg in segments
                ]
            )
            return response

        except Exception as e:
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return TranscribeResponse()

        finally:
            if os.path.exists(audio_path):
                os.remove(audio_path)
                print(f"🗑️ Removed audio file: {audio_path}")


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))

    # Inject class instances here
    downloader = AudioDownloader(output_dir="/tmp", format="wav")
    # Or "base", "small", etc.
    transcriber = Transcriber(model_size="distil-large-v3")

    transcriber_service = TranscriberService(downloader, transcriber)
    add_TranscriberServiceServicer_to_server(transcriber_service, server)

    SERVICE_NAMES = (
        DESCRIPTOR.services_by_name['TranscriberService'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)

    server.add_insecure_port('[::]:50051')
    server.start()
    print("🚀 gRPC server running at http://localhost:50051")
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
