# XSwitch XCC Proto Buffer 协议参考文档

<a name="top"></a>
<a name="user-content-top"></a>

这是[XCC API文档](https://docs.xswitch.cn/xcc-api/)的协议参考，使用[Google Protocol Buffers](https://protobuf.dev/)描述。

本文档只是对具体协议数据格式及类型的参考说明，详细的字段说明和用法请参考[XCC API列表](https://docs.xswitch.cn/xcc-api/api/)，原始的`.proto`文件可以在[proto](../)相关目录中找到。

## 目录
{{range .Files}}
{{$file_name := .Name}}- [{{.Name}}](#{{.Name}})
{{range .Messages}}  - [{{.LongName}}](#{{.FullName}})
{{end}}
{{range .Enums}}  - [{{.LongName}}](#{{.FullName}})
{{end}}
{{range .Extensions}}  - [File-level Extensions](#{{$file_name}}-extensions)
{{end}}
{{range .Services}}  - [{{.Name}}](#{{.FullName}})
{{end}}
{{end}}
- [Scalar Value Types](#scalar-value-types)

{{range .Files}}
{{$file_name := .Name}}
<a name="user-content-{{.Name}}"/>
<a name="{{.Name}}"/>
<p align="right"><a href="#top">Top</a></p>

## {{.Name}}
{{.Description}}

{{range .Messages}}
<a name="user-content-{{.FullName}}"/>
<a name="{{.FullName}}"/>

### {{.LongName}}
{{.Description}}

{{if .HasFields}}
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
{{range .Fields -}}
| {{.Name}} | [{{.LongType}}](#{{.FullType}}) | {{.Label}} | {{nobr .Description}}{{if .DefaultValue}} Default: {{.DefaultValue}}{{end}} |
{{end}}
{{end}}

{{if .HasExtensions}}
| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
{{range .Extensions -}}
| {{.Name}} | {{.LongType}} | {{.ContainingLongType}} | {{.Number}} | {{nobr .Description}}{{if .DefaultValue}} Default: {{.DefaultValue}}{{end}} |
{{end}}
{{end}}

{{end}} <!-- end messages -->

{{range .Enums}}
<a name="user-content-{{.FullName}}"/>
<a name="{{.FullName}}"/>

### {{.LongName}}
{{.Description}}

| Name | Number | Description |
| ---- | ------ | ----------- |
{{range .Values -}}
| {{.Name}} | {{.Number}} | {{nobr .Description}} |
{{end}}

{{end}} <!-- end enums -->

{{if .HasExtensions}}
<a name="{{$file_name}}-extensions"/>

### File-level Extensions
| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
{{range .Extensions -}}
| {{.Name}} | {{.LongType}} | {{.ContainingLongType}} | {{.Number}} | {{nobr .Description}}{{if .DefaultValue}} Default: `{{.DefaultValue}}`{{end}} |
{{end}}
{{end}} <!-- end HasExtensions -->

{{range .Services}}
<a name="user-content-{{.FullName}}"/>
<a name="{{.FullName}}"/>

### {{.Name}}
{{.Description}}

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
{{range .Methods -}}
| {{.Name}} | [{{.RequestLongType}}](#{{.RequestFullType}}) | [{{.ResponseLongType}}](#{{.RequestFullType}}) | {{nobr .Description}} |
{{end}}
{{end}} <!-- end services -->

{{end}}

## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
{{range .Scalars -}}
| <a name="user-content-{{.ProtoType}}" /><a name="{{.ProtoType}}" /> {{.ProtoType}} | {{.Notes}} | {{.CppType}} | {{.JavaType}} | {{.PythonType}} |
{{end}}

<a name="user-content-map-string-string" />

map&lt;[string](#string), [string](#string)&gt;
