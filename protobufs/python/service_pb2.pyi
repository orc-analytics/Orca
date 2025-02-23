from google.protobuf import struct_pb2 as _struct_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ResultStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    RESULT_STATUS_HANDLED_FAILED: _ClassVar[ResultStatus]
    RESULT_STATUS_UNHANDLED_FAILED: _ClassVar[ResultStatus]
    RESULT_STATUS_SUCEEDED: _ClassVar[ResultStatus]
RESULT_STATUS_HANDLED_FAILED: ResultStatus
RESULT_STATUS_UNHANDLED_FAILED: ResultStatus
RESULT_STATUS_SUCEEDED: ResultStatus

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

class WindowEmitStatus(_message.Message):
    __slots__ = ("status",)
    class StatusEnum(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        NO_TRIGGERED_ALGORITHMS: _ClassVar[WindowEmitStatus.StatusEnum]
        PROCESSING_TRIGGERED: _ClassVar[WindowEmitStatus.StatusEnum]
    NO_TRIGGERED_ALGORITHMS: WindowEmitStatus.StatusEnum
    PROCESSING_TRIGGERED: WindowEmitStatus.StatusEnum
    STATUS_FIELD_NUMBER: _ClassVar[int]
    status: WindowEmitStatus.StatusEnum
    def __init__(self, status: _Optional[_Union[WindowEmitStatus.StatusEnum, str]] = ...) -> None: ...

class AlgorithmDependency(_message.Message):
    __slots__ = ("name", "version")
    NAME_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    name: str
    version: str
    def __init__(self, name: _Optional[str] = ..., version: _Optional[str] = ...) -> None: ...

class Algorithm(_message.Message):
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

class Result(_message.Message):
    __slots__ = ("algorithm_name", "version", "status", "float_values", "struct_value", "timestamp")
    ALGORITHM_NAME_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    FLOAT_VALUES_FIELD_NUMBER: _ClassVar[int]
    STRUCT_VALUE_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    algorithm_name: str
    version: str
    status: ResultStatus
    float_values: FloatArray
    struct_value: _struct_pb2.Struct
    timestamp: int
    def __init__(self, algorithm_name: _Optional[str] = ..., version: _Optional[str] = ..., status: _Optional[_Union[ResultStatus, str]] = ..., float_values: _Optional[_Union[FloatArray, _Mapping]] = ..., struct_value: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., timestamp: _Optional[int] = ...) -> None: ...

class FloatArray(_message.Message):
    __slots__ = ("values",)
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, values: _Optional[_Iterable[float]] = ...) -> None: ...

class ProcessorRegistration(_message.Message):
    __slots__ = ("runtime", "supported_algorithms")
    RUNTIME_FIELD_NUMBER: _ClassVar[int]
    SUPPORTED_ALGORITHMS_FIELD_NUMBER: _ClassVar[int]
    runtime: str
    supported_algorithms: _containers.RepeatedCompositeFieldContainer[Algorithm]
    def __init__(self, runtime: _Optional[str] = ..., supported_algorithms: _Optional[_Iterable[_Union[Algorithm, _Mapping]]] = ...) -> None: ...

class ProcessingTask(_message.Message):
    __slots__ = ("task_id", "algorithm", "window", "dependency_results")
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    ALGORITHM_FIELD_NUMBER: _ClassVar[int]
    WINDOW_FIELD_NUMBER: _ClassVar[int]
    DEPENDENCY_RESULTS_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    algorithm: Algorithm
    window: Window
    dependency_results: _containers.RepeatedCompositeFieldContainer[Result]
    def __init__(self, task_id: _Optional[str] = ..., algorithm: _Optional[_Union[Algorithm, _Mapping]] = ..., window: _Optional[_Union[Window, _Mapping]] = ..., dependency_results: _Optional[_Iterable[_Union[Result, _Mapping]]] = ...) -> None: ...

class ExecutionRequest(_message.Message):
    __slots__ = ("task_id", "algorithm", "inputs")
    class InputsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: bytes
        def __init__(self, key: _Optional[str] = ..., value: _Optional[bytes] = ...) -> None: ...
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    ALGORITHM_FIELD_NUMBER: _ClassVar[int]
    INPUTS_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    algorithm: Algorithm
    inputs: _containers.ScalarMap[str, bytes]
    def __init__(self, task_id: _Optional[str] = ..., algorithm: _Optional[_Union[Algorithm, _Mapping]] = ..., inputs: _Optional[_Mapping[str, bytes]] = ...) -> None: ...

class ExecutionResult(_message.Message):
    __slots__ = ("task_id", "status", "outputs")
    class OutputsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: bytes
        def __init__(self, key: _Optional[str] = ..., value: _Optional[bytes] = ...) -> None: ...
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    OUTPUTS_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    status: ResultStatus
    outputs: _containers.ScalarMap[str, bytes]
    def __init__(self, task_id: _Optional[str] = ..., status: _Optional[_Union[ResultStatus, str]] = ..., outputs: _Optional[_Mapping[str, bytes]] = ...) -> None: ...

class DagStateRequest(_message.Message):
    __slots__ = ("window_id",)
    WINDOW_ID_FIELD_NUMBER: _ClassVar[int]
    window_id: str
    def __init__(self, window_id: _Optional[str] = ...) -> None: ...

class DagState(_message.Message):
    __slots__ = ("window", "algorithm_states")
    WINDOW_FIELD_NUMBER: _ClassVar[int]
    ALGORITHM_STATES_FIELD_NUMBER: _ClassVar[int]
    window: Window
    algorithm_states: _containers.RepeatedCompositeFieldContainer[AlgorithmState]
    def __init__(self, window: _Optional[_Union[Window, _Mapping]] = ..., algorithm_states: _Optional[_Iterable[_Union[AlgorithmState, _Mapping]]] = ...) -> None: ...

class AlgorithmState(_message.Message):
    __slots__ = ("algorithm", "status", "result")
    ALGORITHM_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    algorithm: Algorithm
    status: str
    result: Result
    def __init__(self, algorithm: _Optional[_Union[Algorithm, _Mapping]] = ..., status: _Optional[str] = ..., result: _Optional[_Union[Result, _Mapping]] = ...) -> None: ...

class Status(_message.Message):
    __slots__ = ("received", "message")
    RECEIVED_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    received: bool
    message: str
    def __init__(self, received: bool = ..., message: _Optional[str] = ...) -> None: ...

class HealthCheckRequest(_message.Message):
    __slots__ = ("timestamp",)
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    timestamp: int
    def __init__(self, timestamp: _Optional[int] = ...) -> None: ...

class HealthCheckResponse(_message.Message):
    __slots__ = ("status", "message", "metrics")
    class Status(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        STATUS_UNKNOWN: _ClassVar[HealthCheckResponse.Status]
        STATUS_SERVING: _ClassVar[HealthCheckResponse.Status]
        STATUS_TRANSITIONING: _ClassVar[HealthCheckResponse.Status]
        STATUS_NOT_SERVING: _ClassVar[HealthCheckResponse.Status]
    STATUS_UNKNOWN: HealthCheckResponse.Status
    STATUS_SERVING: HealthCheckResponse.Status
    STATUS_TRANSITIONING: HealthCheckResponse.Status
    STATUS_NOT_SERVING: HealthCheckResponse.Status
    STATUS_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    METRICS_FIELD_NUMBER: _ClassVar[int]
    status: HealthCheckResponse.Status
    message: str
    metrics: ProcessorMetrics
    def __init__(self, status: _Optional[_Union[HealthCheckResponse.Status, str]] = ..., message: _Optional[str] = ..., metrics: _Optional[_Union[ProcessorMetrics, _Mapping]] = ...) -> None: ...

class ProcessorMetrics(_message.Message):
    __slots__ = ("active_tasks", "memory_bytes", "cpu_percent", "uptime_seconds")
    ACTIVE_TASKS_FIELD_NUMBER: _ClassVar[int]
    MEMORY_BYTES_FIELD_NUMBER: _ClassVar[int]
    CPU_PERCENT_FIELD_NUMBER: _ClassVar[int]
    UPTIME_SECONDS_FIELD_NUMBER: _ClassVar[int]
    active_tasks: int
    memory_bytes: int
    cpu_percent: float
    uptime_seconds: int
    def __init__(self, active_tasks: _Optional[int] = ..., memory_bytes: _Optional[int] = ..., cpu_percent: _Optional[float] = ..., uptime_seconds: _Optional[int] = ...) -> None: ...
