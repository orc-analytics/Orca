# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [interface.proto](#interface-proto)
    - [Algorithm](#-Algorithm)
    - [Origin](#-Origin)
    - [Payload](#-Payload)
    - [Type](#-Type)
    - [Version](#-Version)
    - [Window](#-Window)
    - [WindowRequest](#-WindowRequest)
    - [WindowResponse](#-WindowResponse)
  
    - [WindowService](#-WindowService)
  
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
| version | [Version](#Version) |  | The version of the algorithm. When versioning, the semver rules should be used in the context of the algorithms result. If the result is backwards compatible, then a minor change, etc. |
| window_type | [Type](#Type) | repeated | The Window type that triggers the algorithm. Many different window types, can trigger one algorithm. |






<a name="-Origin"></a>

### Origin
Defines an arbitrary location, for where an Window was generated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The location name. Best practice is use Snake, Pascal or Camel case consistently. |






<a name="-Payload"></a>

### Payload
Arbitrary information that can be carried along with the window.
It is often useful, when performing batch analysis, to include
&#39;expensive&#39; data that can be queried once, in this structure. This
data can then be access by all algorithms that are triggered by
this window.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The payload data. |






<a name="-Type"></a>

### Type
Defines the type associated with a window. Can be freeform but
must be used consistently across identical window types.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the type. |
| version | [Version](#Version) |  | The version of the window type. |






<a name="-Version"></a>

### Version
A generic versioning struct.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| minor | [int32](#int32) |  |  |
| major | [int32](#int32) |  |  |
| patch | [int32](#int32) |  |  |






<a name="-Window"></a>

### Window
The window definition. The Window is the Cardinal trigger for all
processing DAGs. It defines the complete set of information 
required to successfully run an algorithm, pipeline and/or 
complete DAG.

It should contain the minimal set of information required for the
algorithm to get the relevant data and complete processing.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| start | [string](#string) |  | The start of the window, it the same units as the basis (e.g. Time). |
| end | [string](#string) |  | The end time of the window, in the same units as the basis (e.g. Time). |
| origin | [Origin](#Origin) |  | Where the window was generated. E.g. by an automated service or locally. |
| type | [Type](#Type) |  | The type of the window. It is the Window Type that is the fundamentally distinguishing characteristic between Windows. E.g. Window A may define a region of time where a certain event happened, and Window B may define a sub-region within Window A. Both of these windows will have a unique `` |
| payload | [Payload](#Payload) |  | The additional arbitrary information that can be taken along with the Window. |
| key | [string](#string) |  | A globally unique hash identifying this epoch. |
| parent_key | [string](#string) |  | If this window has been derived from an invoked algorithm within the framework, then the `parent_key` is the key of that Algorithm. |






<a name="-WindowRequest"></a>

### WindowRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| window | [Window](#Window) |  |  |






<a name="-WindowResponse"></a>

### WindowResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [int32](#int32) |  |  |





 

 

 


<a name="-WindowService"></a>

### WindowService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RegisterWindow | [.WindowRequest](#WindowRequest) | [.WindowResponse](#WindowResponse) | Unary Requests - i.e. no open connection whilst processes happen. |
| DeleteWindow | [.WindowRequest](#WindowRequest) | [.WindowResponse](#WindowResponse) |  |
| ReprocessWindow | [.WindowRequest](#WindowRequest) | [.WindowResponse](#WindowResponse) |  |
| ModifyWindow | [.WindowRequest](#WindowRequest) | [.WindowResponse](#WindowResponse) |  |

 



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

