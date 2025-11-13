{{- define "db2" -}}
func (u *{{.ModelName}}Exec) Create(ctx context.Context, data *{{.ModelName}}) error {
	return u.db.WithContext(ctx).Create(data).Error
}

func (u *{{.ModelName}}Exec) FindOne(ctx context.Context, {{.Pr}} int64) (*{{.ModelName}}, error) {
	var data {{.ModelName}}
    err := u.db.WithContext(ctx).Where("{{.GPr}} = ?", {{.Pr}}).First(&data).Error
    if err != nil {
    	return nil, err
    }
    return &data, nil
}

{{- $outer := . -}}
{{- range $index, $notpr := .NotPrs}}
    {{- $gNotpr := index $.GNotPrs $index}}

func (u *{{$outer.ModelName}}Exec) FindBy{{$notpr.Name}}(ctx context.Context, {{$notpr.Name}} {{$notpr.Type}}) (*[]{{$outer.ModelName}}, error) {
	var datas []{{$outer.ModelName}}
    err := u.db.WithContext(ctx).Where("{{$gNotpr}} = ?", {{$notpr.Name}}).Find(&datas).Error
    if err != nil {
    	return nil, err
    }
    return &datas, nil
}
{{- end}}

func (u *{{.ModelName}}Exec) Update(ctx context.Context, data *{{.ModelName}}) error {
	return u.db.WithContext(ctx).Updates(data).Error
}

func (u *{{.ModelName}}Exec) Delete(ctx context.Context, {{.Pr}} int64) error {
	return u.db.WithContext(ctx).Where("{{.GPr}} = ?", {{.Pr}}).Delete(&{{.ModelName}}{}).Error
}
{{- end -}}