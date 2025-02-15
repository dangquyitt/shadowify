import grpc
from concurrent import futures
import transcription_pb2
import transcription_pb2_grpc
import transcription


class TranscriptionServiceServicer(transcription_pb2_grpc.TranscriptionServiceServicer):
    def GetTranscript(self, request, context):
        print(f"Received request for video URL: {request.video_url}")
        transcript = transcription.transcribe_audio(request.video_url)
        return transcription_pb2.GetTranscriptResponse(transcript=transcript)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    transcription_pb2_grpc.add_TranscriptionServiceServicer_to_server(
        TranscriptionServiceServicer(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    print("Server is running on port 50051...")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
