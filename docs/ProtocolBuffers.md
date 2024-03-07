# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [interface.proto](#interface-proto)
    - [Algorithm](#-Algorithm)
    - [Epoch](#-Epoch)
    - [EpochRequest](#-EpochRequest)
    - [EpochResponse](#-EpochResponse)
    - [Origin](#-Origin)
    - [Payload](#-Payload)
    - [Pipeline](#-Pipeline)
    - [Pipeline.AlgorithmDependency](#-Pipeline-AlgorithmDependency)
    - [Type](#-Type)
    - [Version](#-Version)
  
    - [EpochService](#-EpochService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="interface-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## interface.proto



<a name="-Algorithm"></a>

### Algorithm
The definition of an algorithm.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The name of the algorithm. |
| version | [string](#string) |  | The version of the algorithm. Follow [SemVer](https://semver.org/) convention |
| EpochType | [Type](#Type) |  | The Epoch type that triggers the algorithm |






<a name="-Epoch"></a>

### Epoch
The epoch definition. The Epoch is the Cardinal trigger for all
processing DAGs. It defines the complete set of information 
required to successfully run an algorithm, pipeline and/or 
complete DAG.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| start | [string](#string) |  | The start of the epoch, it the same units as the basis. |
| end | [string](#string) |  | The end time of the epoch, in the same units as the basis. |
| origin | [Origin](#Origin) |  | Where the epoch was generated. E.g. by service A or locally. |
| type | [Type](#Type) |  | The type of the epoch. It is the Epoch Type that is the fundamentally distinguishing characteristic between Epochs. E.g. Epoch A may define a region of time where a certain event happened and Epoch B may define a sub-region within Epoch A. |
| payload | [Payload](#Payload) |  | Additional arbitrary information that can be taken along with the Epoch. |
| key | [string](#string) |  | A globally unique hash identifying this epoch |
| parent_key | [string](#string) |  | If this epoch has been derived from an invoked algorithm within the Analytical Framework, then the `parent_key` is the key of that Algorithm. |






<a name="-EpochRequest"></a>

### EpochRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| epoch | [Epoch](#Epoch) |  |  |






<a name="-EpochResponse"></a>

### EpochResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [int32](#int32) |  |  |






<a name="-Origin"></a>

### Origin
Defines an arbitrary location, for where an Epoch was generated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The location name. Best practice is use Snake, Pascal or Camel case consistently. |






<a name="-Payload"></a>

### Payload
Arbitrary information that can be carried along with the epoch.
It is often useful, when performing batch analysis, to include
&#39;expensive&#39; data that can be queried once, in this field. Such
data would be common across algorithms.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The payload data. |






<a name="-Pipeline"></a>

### Pipeline
An explicit declaration of a proessing DAG, defining algorithms
that should be triggered, and in what order, from a single epoch.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The name of the Pipeline. |
| algorithms | [Algorithm](#Algorithm) | repeated | Algorithms to execute as part of the pipeline. |
| dependencies | [Pipeline.AlgorithmDependency](#Pipeline-AlgorithmDependency) | repeated | Algorithm result dependencies |






<a name="-Pipeline-AlgorithmDependency"></a>

### Pipeline.AlgorithmDependency
Message struct for defnining the depencies between algorithms,
in the context of a pipeline


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| parent_algorithm | [Algorithm](#Algorithm) |  | The parent algorithm, that creates the dependent result. |
| dependent_algorithm | [Algorithm](#Algorithm) |  | The dependent algorithm that inherits the result of the parent algorithm. |






<a name="-Type"></a>

### Type
Defines the type associated with an epoch. Can be freeform but
must be used consistently across identical epochs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the type. |
| version | [Version](#Version) |  |  |






<a name="-Version"></a>

### Version
A generic versioning struct.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| minor | [int32](#int32) |  |  |
| major | [int32](#int32) |  |  |
| patch | [int32](#int32) |  |  |





 

 

 


<a name="-EpochService"></a>

### EpochService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RegisterEpoch | [.EpochRequest](#EpochRequest) | [.EpochResponse](#EpochResponse) | Unary Requests |
| DeleteEpoch | [.EpochRequest](#EpochRequest) | [.EpochResponse](#EpochResponse) |  |
| ReprocessEpoch | [.EpochRequest](#EpochRequest) | [.EpochResponse](#EpochResponse) |  |
| ModifyEpoch | [.EpochRequest](#EpochRequest) | [.EpochResponse](#EpochResponse) |  |

 



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

