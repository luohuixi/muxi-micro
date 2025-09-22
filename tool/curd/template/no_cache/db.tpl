{{- define "db2" -}}
func (u *{{.ModelName}}Exec) Create(ctx context.Context, data *{{.ModelName}}) error {
	err := u.Exec.Create(ctx, data, {{.ModelName}}{})
	if err != nil {
		return err
	}
	return nil
}

func (u *{{.ModelName}}Exec) FindOne(ctx context.Context, {{.Pr}} int64) (*{{.ModelName}}, error) {
	var data {{.ModelName}}
	u.Exec.Model.Where("{{.Pr}} = ?", {{.Pr}})
	err := u.Exec.Find(ctx, &data, {{.ModelName}}{})
	if err != nil {
		return nil, err
	}
	if data.{{.Pr}} == 0 {
        return nil, DBNotFound
    }
	return &data, nil
}

{{- $outer := . -}}
{{- range $notpr := .NotPrs}}

func (u *{{$outer.ModelName}}Exec) FindBy{{$notpr.Name}}(ctx context.Context, {{$notpr.Name}} {{$notpr.Type}}) (*[]{{$outer.ModelName}}, error) {
	var datas []{{$outer.ModelName}}
	u.Exec.Model.Where("{{$notpr.Name}} = ?", {{$notpr.Name}})
	err := u.Exec.Find(ctx, &datas, {{$outer.ModelName}}{})
	if err != nil {
		return nil, err
	}
	if len(datas) == 0 {
        return nil, DBNotFound
       }
	return &datas, nil
}
{{- end}}

func (u *{{.ModelName}}Exec) Update(ctx context.Context, data *{{.ModelName}}) error {
	u.Exec.Model.Where("id = ?", data.{{.Pr}})
	err := u.Exec.Update(ctx, data, {{.ModelName}}{})
	if err != nil {
		return err
	}
	return nil
}

func (u *{{.ModelName}}Exec) Delete(ctx context.Context, id int64) error {
	var data {{.ModelName}}
	d, err := u.FindOne(ctx, id)
	if err != nil {
		return err
	}
	data = *d
	err = u.Exec.Delete(ctx, &data, {{.ModelName}}{})
	if err != nil {
		return err
	}
	return nil
}
{{- end -}}