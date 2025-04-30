# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [service.proto](#service-proto)
    - [Algorithm](#-Algorithm)
    - [AlgorithmDependency](#-AlgorithmDependency)
    - [AlgorithmResult](#-AlgorithmResult)
    - [ExecuteDAG](#-ExecuteDAG)
    - [ExecuteDAG.InputsEntry](#-ExecuteDAG-InputsEntry)
    - [ExecutionRequest](#-ExecutionRequest)
    - [ExecutionResult](#-ExecutionResult)
    - [FloatArray](#-FloatArray)
    - [HealthCheckRequest](#-HealthCheckRequest)
    - [HealthCheckResponse](#-HealthCheckResponse)
    - [ProcessingTask](#-ProcessingTask)
    - [ProcessorMetrics](#-ProcessorMetrics)
    - [ProcessorRegistration](#-ProcessorRegistration)
    - [Result](#-Result)
    - [Status](#-Status)
    - [Window](#-Window)
    - [WindowEmitStatus](#-WindowEmitStatus)
    - [WindowType](#-WindowType)
  
    - [HealthCheckResponse.Status](#-HealthCheckResponse-Status)
    - [ResultStatus](#-ResultStatus)
    - [WindowEmitStatus.StatusEnum](#-WindowEmitStatus-StatusEnum)
  
    - [OrcaCore](#-OrcaCore)
    - [OrcaProcessor](#-OrcaProcessor)
  
- [Scalar Value Types](#scalar-value-types)



<a name="service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## service.proto



<a name="-Algorithm"></a>

### Algorithm
Algorithm defines a processing unit that can be executed by processors.
Algorithms form the nodes in the processing DAG and are triggered by specific window types.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the algorithm - must be globally unique This identifies the algorithm across the system |
| version | [string](#string) |  | Version of the algorithm - must follow semantic versioning Allows for algorithm evolution while maintaining compatibility |
| window_type | [WindowType](#WindowType) |  | Type of window that triggers this algorithm References a WindowType that will cause this algorithm to execute |
| dependencies | [AlgorithmDependency](#AlgorithmDependency) | repeated | Other algorithms that this algorithm depends on The algorithm won&#39;t execute until all dependencies have completed Dependencies must not form cycles - this is statically checked on processor registration |






<a name="-AlgorithmDependency"></a>

### AlgorithmDependency
AlgorithmDependency defines a requirement that one algorithm has on another&#39;s results.
These dependencies form the edges in the processing DAG.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the required algorithm Must reference an existing algorithm name in the system |
| version | [string](#string) |  | Version of the required algorithm Must follow semantic versioning (e.g., &#34;1.0.0&#34;) |
| processor_name | [string](#string) |  | Name of the processor that the algorithm is associated with |
| processor_runtime | [string](#string) |  | Runtime of the processor that the algorithm is associated with |






<a name="-AlgorithmResult"></a>

### AlgorithmResult
AlgorithmWindowResult Packaged algorithm and result to a window


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| algorithm | [Algorithm](#Algorithm) |  |  |
| result | [Result](#Result) |  |  |






<a name="-ExecuteDAG"></a>

### ExecuteDAG
ExecuteDAG contains all information needed for a processor to execute
a specific algorithm instance. Sent by the orchestrator to processors.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| task_id | [string](#string) |  | Task ID from the original ProcessingTask Used to correlate results back to the task |
| algorithm | [Algorithm](#Algorithm) |  | Algorithm to execute with its full specification |
| inputs | [ExecuteDAG.InputsEntry](#ExecuteDAG-InputsEntry) | repeated | Input data/parameters for the algorithm Keys are parameter names, values are serialised data Format of data is specific to the algorithm implementation |






<a name="-ExecuteDAG-InputsEntry"></a>

### ExecuteDAG.InputsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [bytes](#bytes) |  |  |






<a name="-ExecutionRequest"></a>

### ExecutionRequest
ExecutionRequest provides a complete view of a processing DAG&#39;s execution
status for a specific window. Used for monitoring and debugging.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| window | [Window](#Window) |  | The window that triggered the algorithm |
| algorithm_results | [AlgorithmResult](#AlgorithmResult) | repeated | Results from dependant algorithms |
| algorithms | [Algorithm](#Algorithm) | repeated | The algorithms to execute |






<a name="-ExecutionResult"></a>

### ExecutionResult



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| task_id | [string](#string) |  | Task ID |
| status | [ResultStatus](#ResultStatus) |  | Execution status |
| outputs | [google.protobuf.Struct](#google-protobuf-Struct) |  | Output data |






<a name="-FloatArray"></a>

### FloatArray
Container for array of float values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| values | [float](#float) | repeated |  |






<a name="-HealthCheckRequest"></a>

### HealthCheckRequest
HealthCheckRequest is sent to processors to verify they are functioning


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| timestamp | [int64](#int64) |  | Timestamp of the request in unix epoch milliseconds Used to measure response latency |






<a name="-HealthCheckResponse"></a>

### HealthCheckResponse
HealthCheckResponse indicates the health status of a processor


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [HealthCheckResponse.Status](#HealthCheckResponse-Status) |  | Current health status |
| message | [string](#string) |  | Optional message providing more detail about the health status |
| metrics | [ProcessorMetrics](#ProcessorMetrics) |  | System metrics about the processor |






<a name="-ProcessingTask"></a>

### ProcessingTask
ProcessingTask represents a single algorithm execution request sent to a processor.
Tasks are streamed to processors as their dependencies are satisfied.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| task_id | [string](#string) |  | Unique ID for this specific task execution Used to correlate results and track execution state |
| algorithm | [Algorithm](#Algorithm) |  | Algorithm to execute Must be one of the algorithms the processor registered support for |
| window | [Window](#Window) |  | Window that triggered this task Provides the time context for the algorithm execution |
| dependency_results | [Result](#Result) | repeated | Results from dependent algorithms Contains all results that this algorithm declared dependencies on All dependencies will be present when task is sent |






<a name="-ProcessorMetrics"></a>

### ProcessorMetrics
ProcessorMetrics provides runtime information about a processor


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| active_tasks | [int32](#int32) |  | Number of algorithms currently being executed |
| memory_bytes | [int64](#int64) |  | Memory usage in bytes |
| cpu_percent | [float](#float) |  | CPU usage percentage (0-100) |
| uptime_seconds | [int64](#int64) |  | Time since processor started in seconds |






<a name="-ProcessorRegistration"></a>

### ProcessorRegistration
ProcessorRegistration is sent by processors when they start up to announce their capabilities
to the orchestrator. This establishes a long-lived connection for receiving tasks.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Unique name of the runtime |
| runtime | [string](#string) |  | Language/runtime of the processor Examples: &#34;python3.9&#34;, &#34;go1.19&#34;, &#34;R4.1&#34; |
| connection_str | [string](#string) |  | The connection string of the orca slave server e.g. grpc://localhost:5433 |
| supported_algorithms | [Algorithm](#Algorithm) | repeated | Algorithms this processor can execute The processor must implement all listed algorithms |






<a name="-Result"></a>

### Result
Result of an algorithm execution


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| algorithm_name | [string](#string) |  | Name of the algorithm that produced the result |
| version | [string](#string) |  | Version of the algorithm that produced the result |
| status | [ResultStatus](#ResultStatus) |  | Status of the result execution |
| single_value | [float](#float) |  | for single number results |
| float_values | [FloatArray](#FloatArray) |  | For numeric array results |
| struct_value | [google.protobuf.Struct](#google-protobuf-Struct) |  | For structured data results (JSON-like) Must follow a map&lt;string, value&gt; schema where value corresponds to https://protobuf.dev/reference/protobuf/google.protobuf/#value |
| timestamp | [int64](#int64) |  | Timestamp when the result was produced |






<a name="-Status"></a>

### Status



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| received | [bool](#bool) |  |  |
| message | [string](#string) |  |  |






<a name="-Window"></a>

### Window
Window represents a time-bounded processing context that triggers algorithm execution. Windows are the primary input that start DAG processing flows.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| from | [uint64](#uint64) |  | Time that the window starts - nanoseconds since epoch Required: Must be &gt; 0 and &lt; to |
| to | [uint64](#uint64) |  | Time that the window ends - nanoseconds since epoch Required: Must be &gt; from |
| window_type_name | [string](#string) |  | The canonical name of the window that uniquely identifies it This allows tracking of window state and results across the system Required: Must be unique within the system, and refer directly to window type |
| window_type_version | [string](#string) |  | The version of the window type, as defined by WindoType |
| origin | [string](#string) |  | A unique identifier that defines where the window came from |






<a name="-WindowEmitStatus"></a>

### WindowEmitStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [WindowEmitStatus.StatusEnum](#WindowEmitStatus-StatusEnum) |  |  |






<a name="-WindowType"></a>

### WindowType
WindowType defines a category of window that can trigger algorithms.
Algorithms subscribe to window types to indicate when they should be executed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the window type - must be globally unique Examples: &#34;daily&#34;, &#34;hourly&#34;, &#34;market_close&#34;, &#34;event_triggered&#34; |
| version | [string](#string) |  | Version of the algorithm. Follows basic semver and allows window types to be changed over time, with traceability |





 


<a name="-HealthCheckResponse-Status"></a>

### HealthCheckResponse.Status
Overall health status of the processor

| Name | Number | Description |
| ---- | ------ | ----------- |
| STATUS_UNKNOWN | 0 | Unknown status - should never be used |
| STATUS_SERVING | 1 | Processor is healthy and ready to accept tasks |
| STATUS_TRANSITIONING | 2 | Processor is starting up or shutting down |
| STATUS_NOT_SERVING | 3 | Processor is not healthy and cannot accept tasks |



<a name="-ResultStatus"></a>

### ResultStatus
ResultStatus indicates the outcome of algorithm execution

| Name | Number | Description |
| ---- | ------ | ----------- |
| RESULT_STATUS_HANDLED_FAILED | 0 | Algorithm failed but the error was handled gracefully The system may retry or skip depending on configuration |
| RESULT_STATUS_UNHANDLED_FAILED | 1 | Algorithm failed with an unexpected error Requires investigation and may halt dependent processing |
| RESULT_STATUS_SUCEEDED | 2 | Algorithm completed successfully Results are valid and can be used by dependent algorithms |



<a name="-WindowEmitStatus-StatusEnum"></a>

### WindowEmitStatus.StatusEnum
A status enum that captures scenarios regarding a window being emmited

| Name | Number | Description |
| ---- | ------ | ----------- |
| NO_TRIGGERED_ALGORITHMS | 0 | When no algorithms could be found that are triggered by this window |
| PROCESSING_TRIGGERED | 1 | When processing has successfully been triggered |


 

 


<a name="-OrcaCore"></a>

### OrcaCore
OrcaCore is the central orchestration service that:
- Manages the lifecycle of processing windows
- Coordinates algorithm execution across distributed processors
- Tracks DAG dependencies and execution state
- Routes results between dependent algorithms

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RegisterProcessor | [.ProcessorRegistration](#ProcessorRegistration) | [.Status](#Status) | Register a processor node and its supported algorithms |
| EmitWindow | [.Window](#Window) | [.WindowEmitStatus](#WindowEmitStatus) | Submit a window for processing |


<a name="-OrcaProcessor"></a>

### OrcaProcessor
OrcaProcessor defines the interface that each processing node must implement.
Processors are language-agnostic services that:
- Execute individual algorithms
- Handle their own internal state
- Report results back to the orchestrator
Orca will schedule processors asynchronously as per the DAG

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ExecuteDagPart | [.ExecutionRequest](#ExecutionRequest) | [.ExecutionResult](#ExecutionResult) stream | Execute part of a DAG with streaming results Server streams back execution results as they become available |
| HealthCheck | [.HealthCheckRequest](#HealthCheckRequest) | [.HealthCheckResponse](#HealthCheckResponse) | Check health/status of processor. i.e. a heartbeat |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

