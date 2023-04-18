# XSwitch XCC Proto Buffer 协议参考文档

<a name="top"></a>
<a name="user-content-top"></a>

这是[XCC API文档](https://docs.xswitch.cn/xcc-api/)的协议参考，使用[Google Protocol Buffers](https://protobuf.dev/)描述。

本文档只是对具体协议数据格式及类型的参考说明，详细的字段说明和用法请参考[XCC API列表](https://docs.xswitch.cn/xcc-api/api/)，原始的`.proto`文件可以在[proto](../)相关目录中找到。

## 目录

- [base.proto](#base.proto)
  - [Debug](#base.Debug)
  - [Filter](#base.Filter)
  - [Filter.And](#base.Filter.And)
  - [Filter.Or](#base.Filter.Or)
  - [Filter.OrderBy](#base.Filter.OrderBy)

  - [Filter.Cond](#base.Filter.Cond)
  - [Filter.OrderByDerection](#base.Filter.OrderByDerection)




- [Scalar Value Types](#scalar-value-types)



<a name="base.proto"/>
<p align="right"><a href="#top">Top</a></p>

## base.proto



<a name="base.Debug"/>

### Debug



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| log_level | [string](#string) |  |  |
| show_sql | [bool](#bool) |  |  |






<a name="base.Filter"/>

### Filter
查询条件


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| and | [Filter.And](#base.Filter.And) |  |  |
| or | [Filter.Or](#base.Filter.Or) |  |  |
| order_by | [Filter.OrderBy](#base.Filter.OrderBy) |  |  |






<a name="base.Filter.And"/>

### Filter.And
and 条件


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cond | [string](#string) |  | 条件 enum Cond |
| column | [string](#string) |  | 字段 |
| args | [string](#string) | repeated | 判断参数 |






<a name="base.Filter.Or"/>

### Filter.Or
or 条件


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cond | [string](#string) |  | 条件 |
| column | [string](#string) |  | 字段 |
| args | [string](#string) | repeated | 参数 |






<a name="base.Filter.OrderBy"/>

### Filter.OrderBy
order by


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| order_by_direction | [Filter.OrderByDerection](#base.Filter.OrderByDerection) |  | 条件 |
| column | [string](#string) |  | 字段 |





 <!-- end messages -->


<a name="base.Filter.Cond"/>

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



<a name="base.Filter.OrderByDerection"/>

### Filter.OrderByDerection


| Name | Number | Description |
| ---- | ------ | ----------- |
| Asc | 0 |  |
| Desc | 1 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

