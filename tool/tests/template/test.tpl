{{- define "test" -}}
{{range $func := .Func}}
func Test{{$func.FuncName}}(t *testing.T) {
    type testCase struct {
        description string
        {{range $input := $func.Params -}}
        {{$input.Param}} {{$input.ParamType}}
        {{end -}}
        {{range $return := $func.Returns -}}
        {{$return.Param}} {{$return.ParamType}}
        {{end}}
    }
    tests := []testCase{
        // TODO: 添加测试用例
        
    }

    for _, test := range tests {
        t.Run(test.description, func(t *testing.T) {
            // TODO: 添加测试逻辑
        })
    }
}
{{end -}}
{{- end -}}