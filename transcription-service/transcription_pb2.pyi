from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class GetTranscriptRequest(_message.Message):
    __slots__ = ("video_url",)
    VIDEO_URL_FIELD_NUMBER: _ClassVar[int]
    video_url: str
    def __init__(self, video_url: _Optional[str] = ...) -> None: ...

class GetTranscriptResponse(_message.Message):
    __slots__ = ("transcript",)
    TRANSCRIPT_FIELD_NUMBER: _ClassVar[int]
    transcript: str
    def __init__(self, transcript: _Optional[str] = ...) -> None: ...
