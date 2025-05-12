from concurrent import futures
import grpc
import yt_dlp
import metadata_pb2
import metadata_pb2_grpc
from metadata_pb2 import DESCRIPTOR
from metadata_pb2_grpc import MetadataServiceServicer, add_MetadataServiceServicer_to_server
from grpc_reflection.v1alpha import reflection
from metadata_extractor import MetadataExtractor


class MetadataServiceServicer(metadata_pb2_grpc.MetadataService):
    def __init__(self, extractor: MetadataExtractor):
        self.extractor = extractor

    def GetMetadata(self, request, context) -> metadata_pb2.MetadataResponse:
        try:
            info = self.extractor.extract(request.video_code)
            return metadata_pb2.MetadataResponse(
                title=info.get("title", ""),
                duration=info.get("duration", 0),
                duration_string=info.get("duration_string", ""),
                thumbnail=info.get("thumbnail", ""),
                view_count=info.get("view_count", 0),
                tags=info.get("tags", [])
            )
        except Exception as e:
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return metadata_pb2.MetadataResponse()


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    metadata_extractor = MetadataExtractor()

    add_MetadataServiceServicer_to_server(
        MetadataServiceServicer(extractor=metadata_extractor), server)
    SERVICE_NAMES = (
        DESCRIPTOR.services_by_name['MetadataService'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(
        SERVICE_NAMES, server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("gRPC server listening on port 50051 🚀")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
