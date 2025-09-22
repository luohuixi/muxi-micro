{{- define "test" -}}
{{range $func := .Func}}
func Test{{$func.FuncName}}(t *testing.T) {
    var MockInjection = func() {{$func.Receive}} {
        // TODO: 添加注入的逻辑，推荐 "go.uber.org/mock/gomock"
	}
	object := MockInjection()

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
