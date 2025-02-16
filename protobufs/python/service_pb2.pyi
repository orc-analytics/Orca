from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Window(_message.Message):
    __slots__ = ("to", "name")
    FROM_FIELD_NUMBER: _ClassVar[int]
    TO_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    to: int
    name: str
    def __init__(self, to: _Optional[int] = ..., name: _Optional[str] = ..., **kwargs) -> None: ...

class WindowType(_message.Message):
    __slots__ = ("name",)
    NAME_FIELD_NUMBER: _ClassVar[int]
    name: str
    def __init__(self, name: _Optional[str] = ...) -> None: ...

class AlgorithmDependency(_message.Message):
    __slots__ = ("name", "version")
    NAME_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    name: str
    version: str
    def __init__(self, name: _Optional[str] = ..., version: _Optional[str] = ...) -> None: ...

class AlgorithmType(_message.Message):
    __slots__ = ("name", "version", "window_type_name", "dependencies")
    NAME_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    WINDOW_TYPE_NAME_FIELD_NUMBER: _ClassVar[int]
    DEPENDENCIES_FIELD_NUMBER: _ClassVar[int]
    name: str
    version: str
    window_type_name: str
    dependencies: _containers.RepeatedCompositeFieldContainer[AlgorithmDependency]
    def __init__(self, name: _Optional[str] = ..., version: _Optional[str] = ..., window_type_name: _Optional[str] = ..., dependencies: _Optional[_Iterable[_Union[AlgorithmDependency, _Mapping]]] = ...) -> None: ...

class Status(_message.Message):
    __slots__ = ("recieved",)
    RECIEVED_FIELD_NUMBER: _ClassVar[int]
    recieved: bool
    def __init__(self, recieved: bool = ...) -> None: ...
