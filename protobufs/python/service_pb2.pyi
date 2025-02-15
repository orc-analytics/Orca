from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Window(_message.Message):
    __slots__ = ("to", "name")
    FROM_FIELD_NUMBER: _ClassVar[int]
    TO_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    to: int
    name: str
    def __init__(self, to: _Optional[int] = ..., name: _Optional[str] = ..., **kwargs) -> None: ...

class Status(_message.Message):
    __slots__ = ("recieved",)
    RECIEVED_FIELD_NUMBER: _ClassVar[int]
    recieved: bool
    def __init__(self, recieved: bool = ...) -> None: ...
