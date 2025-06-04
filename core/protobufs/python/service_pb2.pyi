from google.protobuf import struct_pb2 as _struct_pb2
from vendor import validate_pb2 as _validate_pb2
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

class DataGetter(_message.Message):
    __slots__ = ("name", "window_type", "ttl_seconds", "max_size_bytes")
    NAME_FIELD_NUMBER: _ClassVar[int]
    WINDOW_TYPE_FIELD_NUMBER: _ClassVar[int]
    TTL_SECONDS_FIELD_NUMBER: _ClassVar[int]
    MAX_SIZE_BYTES_FIELD_NUMBER: _ClassVar[int]
    name: str
    window_type: WindowType
    ttl_seconds: int
    max_size_bytes: int
    def __init__(self, name: _Optional[str] = ..., window_type: _Optional[_Union[WindowType, _Mapping]] = ..., ttl_seconds: _Optional[int] = ..., max_size_bytes: _Optional[int] = ...) -> None: ...

class CacheConnectionInfo(_message.Message):
    __slots__ = ("cache_type", "connection_string")
    CACHE_TYPE_FIELD_NUMBER: _ClassVar[int]
    CONNECTION_STRING_FIELD_NUMBER: _ClassVar[int]
    cache_type: str
    connection_string: str
    def __init__(self, cache_type: _Optional[str] = ..., connection_string: _Optional[str] = ...) -> None: ...

class DataGetterExecutionTask(_message.Message):
    __slots__ = ("execution_id", "window", "data_getter", "cache_key", "cache_info")
    EXECUTION_ID_FIELD_NUMBER: _ClassVar[int]
    WINDOW_FIELD_NUMBER: _ClassVar[int]
    DATA_GETTER_FIELD_NUMBER: _ClassVar[int]
    CACHE_KEY_FIELD_NUMBER: _ClassVar[int]
    CACHE_INFO_FIELD_NUMBER: _ClassVar[int]
    execution_id: str
    window: Window
    data_getter: DataGetter
    cache_key: str
    cache_info: CacheConnectionInfo
    def __init__(self, execution_id: _Optional[str] = ..., window: _Optional[_Union[Window, _Mapping]] = ..., data_getter: _Optional[_Union[DataGetter, _Mapping]] = ..., cache_key: _Optional[str] = ..., cache_info: _Optional[_Union[CacheConnectionInfo, _Mapping]] = ...) -> None: ...

class DataGetterResult(_message.Message):
    __slots__ = ("execution_id", "status", "cache_key", "data_size_bytes", "error_message", "timestamp")
    class Status(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        SUCCESS: _ClassVar[DataGetterResult.Status]
        FAILED: _ClassVar[DataGetterResult.Status]
        CACHE_STORE_FAILED: _ClassVar[DataGetterResult.Status]
    SUCCESS: DataGetterResult.Status
    FAILED: DataGetterResult.Status
    CACHE_STORE_FAILED: DataGetterResult.Status
    EXECUTION_ID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    CACHE_KEY_FIELD_NUMBER: _ClassVar[int]
    DATA_SIZE_BYTES_FIELD_NUMBER: _ClassVar[int]
    ERROR_MESSAGE_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    execution_id: str
    status: DataGetterResult.Status
    cache_key: str
    data_size_bytes: int
    error_message: str
    timestamp: int
    def __init__(self, execution_id: _Optional[str] = ..., status: _Optional[_Union[DataGetterResult.Status, str]] = ..., cache_key: _Optional[str] = ..., data_size_bytes: _Optional[int] = ..., error_message: _Optional[str] = ..., timestamp: _Optional[int] = ...) -> None: ...

class CachedDataReference(_message.Message):
    __slots__ = ("data_getter_name", "cache_key", "cache_info", "cached_timestamp", "data_size_bytes")
    DATA_GETTER_NAME_FIELD_NUMBER: _ClassVar[int]
    CACHE_KEY_FIELD_NUMBER: _ClassVar[int]
    CACHE_INFO_FIELD_NUMBER: _ClassVar[int]
    CACHED_TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    DATA_SIZE_BYTES_FIELD_NUMBER: _ClassVar[int]
    data_getter_name: str
    cache_key: str
    cache_info: CacheConnectionInfo
    cached_timestamp: int
    data_size_bytes: int
    def __init__(self, data_getter_name: _Optional[str] = ..., cache_key: _Optional[str] = ..., cache_info: _Optional[_Union[CacheConnectionInfo, _Mapping]] = ..., cached_timestamp: _Optional[int] = ..., data_size_bytes: _Optional[int] = ...) -> None: ...

class Window(_message.Message):
    __slots__ = ("time_from", "time_to", "window_type_name", "window_type_version", "origin", "metadata")
    TIME_FROM_FIELD_NUMBER: _ClassVar[int]
    TIME_TO_FIELD_NUMBER: _ClassVar[int]
    WINDOW_TYPE_NAME_FIELD_NUMBER: _ClassVar[int]
    WINDOW_TYPE_VERSION_FIELD_NUMBER: _ClassVar[int]
    ORIGIN_FIELD_NUMBER: _ClassVar[int]
    METADATA_FIELD_NUMBER: _ClassVar[int]
    time_from: int
    time_to: int
    window_type_name: str
    window_type_version: str
    origin: str
    metadata: _struct_pb2.Struct
    def __init__(self, time_from: _Optional[int] = ..., time_to: _Optional[int] = ..., window_type_name: _Optional[str] = ..., window_type_version: _Optional[str] = ..., origin: _Optional[str] = ..., metadata: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class WindowType(_message.Message):
    __slots__ = ("name", "version")
    NAME_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    name: str
    version: str
    def __init__(self, name: _Optional[str] = ..., version: _Optional[str] = ...) -> None: ...

class WindowEmitStatus(_message.Message):
    __slots__ = ("status",)
    class StatusEnum(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        NO_TRIGGERED_ALGORITHMS: _ClassVar[WindowEmitStatus.StatusEnum]
        PROCESSING_TRIGGERED: _ClassVar[WindowEmitStatus.StatusEnum]
        TRIGGERING_FAILED: _ClassVar[WindowEmitStatus.StatusEnum]
    NO_TRIGGERED_ALGORITHMS: WindowEmitStatus.StatusEnum
    PROCESSING_TRIGGERED: WindowEmitStatus.StatusEnum
    TRIGGERING_FAILED: WindowEmitStatus.StatusEnum
    STATUS_FIELD_NUMBER: _ClassVar[int]
    status: WindowEmitStatus.StatusEnum
    def __init__(self, status: _Optional[_Union[WindowEmitStatus.StatusEnum, str]] = ...) -> None: ...

class AlgorithmDependency(_message.Message):
    __slots__ = ("name", "version", "processor_name", "processor_runtime")
    NAME_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    PROCESSOR_NAME_FIELD_NUMBER: _ClassVar[int]
    PROCESSOR_RUNTIME_FIELD_NUMBER: _ClassVar[int]
    name: str
    version: str
    processor_name: str
    processor_runtime: str
    def __init__(self, name: _Optional[str] = ..., version: _Optional[str] = ..., processor_name: _Optional[str] = ..., processor_runtime: _Optional[str] = ...) -> None: ...

class Algorithm(_message.Message):
    __slots__ = ("name", "version", "window_type", "dependencies", "required_data_getters")
    NAME_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    WINDOW_TYPE_FIELD_NUMBER: _ClassVar[int]
    DEPENDENCIES_FIELD_NUMBER: _ClassVar[int]
    REQUIRED_DATA_GETTERS_FIELD_NUMBER: _ClassVar[int]
    name: str
    version: str
    window_type: WindowType
    dependencies: _containers.RepeatedCompositeFieldContainer[AlgorithmDependency]
    required_data_getters: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, name: _Optional[str] = ..., version: _Optional[str] = ..., window_type: _Optional[_Union[WindowType, _Mapping]] = ..., dependencies: _Optional[_Iterable[_Union[AlgorithmDependency, _Mapping]]] = ..., required_data_getters: _Optional[_Iterable[str]] = ...) -> None: ...

class FloatArray(_message.Message):
    __slots__ = ("values",)
    VALUES_FIELD_NUMBER: _ClassVar[int]
    values: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, values: _Optional[_Iterable[float]] = ...) -> None: ...

class Result(_message.Message):
    __slots__ = ("status", "single_value", "float_values", "struct_value", "timestamp")
    STATUS_FIELD_NUMBER: _ClassVar[int]
    SINGLE_VALUE_FIELD_NUMBER: _ClassVar[int]
    FLOAT_VALUES_FIELD_NUMBER: _ClassVar[int]
    STRUCT_VALUE_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    status: ResultStatus
    single_value: float
    float_values: FloatArray
    struct_value: _struct_pb2.Struct
    timestamp: int
    def __init__(self, status: _Optional[_Union[ResultStatus, str]] = ..., single_value: _Optional[float] = ..., float_values: _Optional[_Union[FloatArray, _Mapping]] = ..., struct_value: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., timestamp: _Optional[int] = ...) -> None: ...

class ProcessorRegistration(_message.Message):
    __slots__ = ("name", "runtime", "connection_str", "supported_algorithms", "data_getters")
    NAME_FIELD_NUMBER: _ClassVar[int]
    RUNTIME_FIELD_NUMBER: _ClassVar[int]
    CONNECTION_STR_FIELD_NUMBER: _ClassVar[int]
    SUPPORTED_ALGORITHMS_FIELD_NUMBER: _ClassVar[int]
    DATA_GETTERS_FIELD_NUMBER: _ClassVar[int]
    name: str
    runtime: str
    connection_str: str
    supported_algorithms: _containers.RepeatedCompositeFieldContainer[Algorithm]
    data_getters: _containers.RepeatedCompositeFieldContainer[DataGetter]
    def __init__(self, name: _Optional[str] = ..., runtime: _Optional[str] = ..., connection_str: _Optional[str] = ..., supported_algorithms: _Optional[_Iterable[_Union[Algorithm, _Mapping]]] = ..., data_getters: _Optional[_Iterable[_Union[DataGetter, _Mapping]]] = ...) -> None: ...

class ExecutionRequest(_message.Message):
    __slots__ = ("exec_id", "window", "algorithm_results", "algorithms", "cached_data")
    EXEC_ID_FIELD_NUMBER: _ClassVar[int]
    WINDOW_FIELD_NUMBER: _ClassVar[int]
    ALGORITHM_RESULTS_FIELD_NUMBER: _ClassVar[int]
    ALGORITHMS_FIELD_NUMBER: _ClassVar[int]
    CACHED_DATA_FIELD_NUMBER: _ClassVar[int]
    exec_id: str
    window: Window
    algorithm_results: _containers.RepeatedCompositeFieldContainer[AlgorithmResult]
    algorithms: _containers.RepeatedCompositeFieldContainer[Algorithm]
    cached_data: _containers.RepeatedCompositeFieldContainer[CachedDataReference]
    def __init__(self, exec_id: _Optional[str] = ..., window: _Optional[_Union[Window, _Mapping]] = ..., algorithm_results: _Optional[_Iterable[_Union[AlgorithmResult, _Mapping]]] = ..., algorithms: _Optional[_Iterable[_Union[Algorithm, _Mapping]]] = ..., cached_data: _Optional[_Iterable[_Union[CachedDataReference, _Mapping]]] = ...) -> None: ...

class ExecutionResult(_message.Message):
    __slots__ = ("exec_id", "algorithm_result")
    EXEC_ID_FIELD_NUMBER: _ClassVar[int]
    ALGORITHM_RESULT_FIELD_NUMBER: _ClassVar[int]
    exec_id: str
    algorithm_result: AlgorithmResult
    def __init__(self, exec_id: _Optional[str] = ..., algorithm_result: _Optional[_Union[AlgorithmResult, _Mapping]] = ...) -> None: ...

class AlgorithmResult(_message.Message):
    __slots__ = ("algorithm", "result")
    ALGORITHM_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    algorithm: Algorithm
    result: Result
    def __init__(self, algorithm: _Optional[_Union[Algorithm, _Mapping]] = ..., result: _Optional[_Union[Result, _Mapping]] = ...) -> None: ...

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
