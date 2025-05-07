from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class TranscribeRequest(_message.Message):
    __slots__ = ("video_code",)
    VIDEO_CODE_FIELD_NUMBER: _ClassVar[int]
    video_code: str
    def __init__(self, video_code: _Optional[str] = ...) -> None: ...

class Segment(_message.Message):
    __slots__ = ("text", "start_time", "end_time")
    TEXT_FIELD_NUMBER: _ClassVar[int]
    START_TIME_FIELD_NUMBER: _ClassVar[int]
    END_TIME_FIELD_NUMBER: _ClassVar[int]
    text: str
    start_time: float
    end_time: float
    def __init__(self, text: _Optional[str] = ..., start_time: _Optional[float] = ..., end_time: _Optional[float] = ...) -> None: ...

class TranscribeResponse(_message.Message):
    __slots__ = ("segments",)
    SEGMENTS_FIELD_NUMBER: _ClassVar[int]
    segments: _containers.RepeatedCompositeFieldContainer[Segment]
    def __init__(self, segments: _Optional[_Iterable[_Union[Segment, _Mapping]]] = ...) -> None: ...
