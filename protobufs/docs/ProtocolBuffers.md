# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [interface.proto](#interface-proto)
    - [Algorithm](#AnalyticsFrameworkInterface-Algorithm)
    - [Epoch](#AnalyticsFrameworkInterface-Epoch)
    - [Origin](#AnalyticsFrameworkInterface-Origin)
    - [Payload](#AnalyticsFrameworkInterface-Payload)
    - [Trigger](#AnalyticsFrameworkInterface-Trigger)
    - [Type](#AnalyticsFrameworkInterface-Type)
  
    - [TriggerType](#AnalyticsFrameworkInterface-TriggerType)
  
- [Scalar Value Types](#scalar-value-types)



<a name="interface-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## interface.proto



<a name="AnalyticsFrameworkInterface-Algorithm"></a>

### Algorithm



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| version | [string](#string) |  |  |






<a name="AnalyticsFrameworkInterface-Epoch"></a>

### Epoch



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| start | [string](#string) |  |  |
| end | [string](#string) |  |  |
| origin | [Origin](#AnalyticsFrameworkInterface-Origin) |  |  |
| type | [Type](#AnalyticsFrameworkInterface-Type) |  |  |
| payload | [Payload](#AnalyticsFrameworkInterface-Payload) |  |  |






<a name="AnalyticsFrameworkInterface-Origin"></a>

### Origin
Here is an example comment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="AnalyticsFrameworkInterface-Payload"></a>

### Payload



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| payload | [bytes](#bytes) |  |  |






<a name="AnalyticsFrameworkInterface-Trigger"></a>

### Trigger



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| triggerType | [TriggerType](#AnalyticsFrameworkInterface-TriggerType) |  |  |






<a name="AnalyticsFrameworkInterface-Type"></a>

### Type



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [string](#string) |  |  |





 


<a name="AnalyticsFrameworkInterface-TriggerType"></a>

### TriggerType


| Name | Number | Description |
| ---- | ------ | ----------- |
| CRON | 0 |  |
| EPOCH | 1 |  |


 

 

 



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

