# XSwitch XCC Proto Buffer 协议参考文档

<a name="top"></a>
<a name="user-content-top"></a>

这是[XCC API文档](https://docs.xswitch.cn/xcc-api/)的协议参考，使用[Google Protocol Buffers](https://protobuf.dev/)描述。

本文档只是对具体协议的参考说明，详细的字段说明和用法请参考[XCC API列表](https://docs.xswitch.cn/xcc-api/api/)，原始的`.proto`文件可以在[../xctrl](../xctrl)中找到。

## 目录

- [core/proto/base/base.proto](#core_proto_base_base-proto)
    - [Debug](#base-Debug)
    - [Filter](#base-Filter)
    - [Filter.And](#base-Filter-And)
    - [Filter.Or](#base-Filter-Or)
    - [Filter.OrderBy](#base-Filter-OrderBy)
  
    - [Filter.Cond](#base-Filter-Cond)
    - [Filter.OrderByDerection](#base-Filter-OrderByDerection)
  
- [Scalar Value Types](#scalar-value-types)



<a name="core_proto_base_base-proto"></a>
<a name="user-content-core_proto_base_base-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## core/proto/base/base.proto



<a name="base-Debug"></a>
<a name="user-content-base-Debug"></a>

### Debug



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| log_level | [string](#string) |  |  |
| show_sql | [bool](#bool) |  |  |






<a name="base-Filter"></a>
<a name="user-content-base-Filter"></a>

### Filter
查询条件


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| and | [Filter.And](#base-Filter-And) |  |  |
| or | [Filter.Or](#base-Filter-Or) |  |  |
| order_by | [Filter.OrderBy](#base-Filter-OrderBy) |  |  |






<a name="base-Filter-And"></a>
<a name="user-content-base-Filter-And"></a>

### Filter.And
and 条件


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cond | [string](#string) |  | 条件 enum Cond |
| column | [string](#string) |  | 字段 |
| args | [string](#string) | repeated | 判断参数 |






<a name="base-Filter-Or"></a>
<a name="user-content-base-Filter-Or"></a>

### Filter.Or
or 条件


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cond | [string](#string) |  | 条件 |
| column | [string](#string) |  | 字段 |
| args | [string](#string) | repeated | 参数 |






<a name="base-Filter-OrderBy"></a>
<a name="user-content-base-Filter-OrderBy"></a>

### Filter.OrderBy
order by


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| order_by_direction | [Filter.OrderByDerection](#base-Filter-OrderByDerection) |  | 条件 |
| column | [string](#string) |  | 字段 |





 <!-- end messages -->


<a name="base-Filter-Cond"></a>
<a name="user-content-base-Filter-Cond"></a>

### Filter.Cond


| Name | Number | Description |
| ---- | ------ | ----------- |
| None | 0 |  |
| Eq | 1 | a = ? |
| Neq | 2 | a <> ? |
| Gt | 3 | a > ? |
| Gte | 4 | a >= ? |
| Lt | 5 | a < ? |
| Lte | 6 | a <= ? |
| Like | 7 | a Like ? |
| In | 8 | a IN (?,?,?) |
| NotIn | 9 | a NOT IN (?, ?, ?) |
| IsNull | 10 | a IS NULL [] |
| NotNull | 11 | a IS NOT NULL [] |
| Between | 12 | a BETWEEN 1 AND 2 |



<a name="base-Filter-OrderByDerection"></a>
<a name="user-content-base-Filter-OrderByDerection"></a>

### Filter.OrderByDerection


| Name | Number | Description |
| ---- | ------ | ----------- |
| Asc | 0 |  |
| Desc | 1 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /><a name="user-content-double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /><a name="user-content-float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /><a name="user-content-int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /><a name="user-content-int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /><a name="user-content-uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /><a name="user-content-uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /><a name="user-content-sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /><a name="user-content-sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /><a name="user-content-fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /><a name="user-content-fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /><a name="user-content-sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /><a name="user-content-sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /><a name="user-content-bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /><a name="user-content-string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /><a name="user-content-bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

