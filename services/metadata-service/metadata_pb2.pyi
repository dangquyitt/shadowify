from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class MetadataRequest(_message.Message):
    __slots__ = ("video_code",)
    VIDEO_CODE_FIELD_NUMBER: _ClassVar[int]
    video_code: str
    def __init__(self, video_code: _Optional[str] = ...) -> None: ...

class MetadataResponse(_message.Message):
    __slots__ = ("title", "duration", "duration_string", "thumbnail", "view_count", "tags")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DURATION_FIELD_NUMBER: _ClassVar[int]
    DURATION_STRING_FIELD_NUMBER: _ClassVar[int]
    THUMBNAIL_FIELD_NUMBER: _ClassVar[int]
    VIEW_COUNT_FIELD_NUMBER: _ClassVar[int]
    TAGS_FIELD_NUMBER: _ClassVar[int]
    title: str
    duration: int
    duration_string: str
    thumbnail: str
    view_count: int
    tags: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, title: _Optional[str] = ..., duration: _Optional[int] = ..., duration_string: _Optional[str] = ..., thumbnail: _Optional[str] = ..., view_count: _Optional[int] = ..., tags: _Optional[_Iterable[str]] = ...) -> None: ...
