# XSwitch XCC Proto Buffer 协议参考文档

<a name="top"></a>
<a name="user-content-top"></a>

这是[XCC API文档](https://docs.xswitch.cn/xcc-api/)的协议参考，使用[Google Protocol Buffers](https://protobuf.dev/)描述。

本文档只是对具体协议数据格式及类型的参考说明，详细的字段说明和用法请参考[XCC API列表](https://docs.xswitch.cn/xcc-api/api/)，原始的`.proto`文件可以在[proto](../)相关目录中找到。

## 目录
{{range .Files}}
{{$file_name := .Name}}- [{{.Name}}](#{{.Name | anchor}})
  {{- if .Messages }}
  {{range .Messages}}  - [{{.LongName}}](#{{.FullName | anchor}})
  {{end}}
  {{- end -}}
  {{- if .Enums }}
  {{range .Enums}}  - [{{.LongName}}](#{{.FullName | anchor}})
  {{end}}
  {{- end -}}
  {{- if .Extensions }}
  {{range .Extensions}}  - [File-level Extensions](#{{$file_name | anchor}}-extensions)
  {{end}}
  {{- end -}}
  {{- if .Services }}
  {{range .Services}}  - [{{.Name}}](#{{.FullName | anchor}})
  {{end}}
  {{- end -}}
{{end}}
- [Scalar Value Types](#scalar-value-types)

{{range .Files}}
{{$file_name := .Name}}
<a name="{{.Name | anchor}}"></a>
<a name="user-content-{{.Name | anchor}}"></a>
<p align="right"><a href="#top">Top</a></p>

## {{.Name}}
{{.Description}}

{{range .Messages}}
<a name="{{.FullName | anchor}}"></a>
<a name="user-content-{{.FullName | anchor}}"></a>

### {{.LongName}}
{{.Description}}

{{if .HasFields}}
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
{{range .Fields -}}
  | {{.Name}} | [{{.LongType}}](#{{.FullType | anchor}}) | {{.Label}} | {{if (index .Options "deprecated"|default false)}}**Deprecated.** {{end}}{{nobr .Description}}{{if .DefaultValue}} Default: {{.DefaultValue}}{{end}} |
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
<a name="{{.FullName | anchor}}"></a>
<a name="user-content-{{.FullName | anchor}}"></a>

### {{.LongName}}
{{.Description}}

| Name | Number | Description |
| ---- | ------ | ----------- |
{{range .Values -}}
  | {{.Name}} | {{.Number}} | {{nobr .Description}} |
{{end}}

{{end}} <!-- end enums -->

{{if .HasExtensions}}
<a name="{{$file_name | anchor}}-extensions"></a>
<a name="user-content-{{$file_name | anchor}}-extensions"></a>

### File-level Extensions
| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
{{range .Extensions -}}
  | {{.Name}} | {{.LongType}} | {{.ContainingLongType}} | {{.Number}} | {{nobr .Description}}{{if .DefaultValue}} Default: `{{.DefaultValue}}`{{end}} |
{{end}}
{{end}} <!-- end HasExtensions -->

{{range .Services}}
<a name="{{.FullName | anchor}}"></a>
<a name="user-content-{{.FullName | anchor}}"></a>

### {{.Name}}
{{.Description}}

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
{{range .Methods -}}
  | {{.Name}} | [{{.RequestLongType}}](#{{.RequestFullType | anchor}}){{if .RequestStreaming}} stream{{end}} | [{{.ResponseLongType}}](#{{.ResponseFullType | anchor}}){{if .ResponseStreaming}} stream{{end}} | {{nobr .Description}} |
{{end}}
{{end}} <!-- end services -->

{{end}}

## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
{{range .Scalars -}}
  | <a name="{{.ProtoType | anchor}}" /><a name="user-content-{{.ProtoType | anchor}}" /> {{.ProtoType}} | {{.Notes}} | {{.CppType}} | {{.JavaType}} | {{.PythonType}} | {{.GoType}} | {{.CSharp}} | {{.PhpType}} | {{.RubyType}} |
{{end}}
