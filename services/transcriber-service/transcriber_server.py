import grpc
from concurrent import futures
import time
import os

from transcriber_pb2 import Segment, TranscribeResponse, DESCRIPTOR
from transcriber_pb2_grpc import TranscriberServiceServicer, add_TranscriberServiceServicer_to_server

from audio_dowloader import download_audio_from_youtube
from extract_segments import extract_segments
from grpc_reflection.v1alpha import reflection


class TranscriberService(TranscriberServiceServicer):
    def Transcribe(self, request, context):
        video_code = request.video_code
        print(f"📥 Received video code: {video_code}")

        try:
            audio_path = download_audio_from_youtube(video_code)
            print(f"🎵 Audio downloaded at: {audio_path}")

            segments = extract_segments(audio_path)

            response = TranscribeResponse(
                segments=[
                    Segment(
                        text=seg["text"],
                        start_time=seg["start"],
                        end_time=seg["end"]
                    ) for seg in segments
                ]
            )

            os.remove(audio_path)

            return response

        except Exception as e:
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return TranscribeResponse()  # empty

        finally:
            if os.path.exists(audio_path):
                os.remove(audio_path)
                print(f"🗑️ Removed audio file: {audio_path}")


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    add_TranscriberServiceServicer_to_server(TranscriberService(), server)
    SERVICE_NAMES = (
        DESCRIPTOR.services_by_name['TranscriberService'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    server.add_insecure_port('[::]:50051')
    print("🚀 gRPC server running at http://localhost:50051")
    server.start()

    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
