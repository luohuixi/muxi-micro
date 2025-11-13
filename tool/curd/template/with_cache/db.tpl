{{- define "db" -}}
func (u *{{.ModelName}}Exec) Create(ctx context.Context, data *{{.ModelName}}) error {
	err := u.db.WithContext(ctx).Create(data).Error
    if err != nil {
    	return err
    }
    go u.DelCache(ctx, data)
    return nil
}

func (u *{{.ModelName}}Exec) FindOne(ctx context.Context, {{.Pr}} int64) (*{{.ModelName}}, error) {
	cachestr := fmt.Sprintf("%s%v", cache{{.ModelName}}{{.Pr}}Prefix, {{.Pr}})
	datacache := u.Get(ctx, cachestr)
    if datacache != nil {
    	return datacache, nil
    }

    var data {{.ModelName}}
    err := u.db.WithContext(ctx).Where("{{.GPr}} = ?", {{.Pr}}).First(&data).Error
    if err != nil {
    	return nil, err
    }

    go u.Set(cachestr, &data)
    return &data, nil
}

{{- $outer := . -}}
{{- range $index, $notpr := .NotPrs}}
    {{- $gNotpr := index $.GNotPrs $index}}

func (u *{{$outer.ModelName}}Exec) FindBy{{$notpr.Name}}(ctx context.Context, {{$notpr.Name}} {{$notpr.Type}}) (*[]{{$outer.ModelName}}, error) {
	cachestr := fmt.Sprintf("%s%v", cache{{$outer.ModelName}}{{$notpr.Name}}Prefix, {{$notpr.Name}})
	datacache := u.GetMany(ctx, cachestr)
    if datacache != nil {
    	return datacache, nil
    }

    var datas []{{$outer.ModelName}}
    err := u.db.WithContext(ctx).Where("{{$gNotpr}} = ?", {{$notpr.Name}}).Find(&datas).Error
    if err != nil {
    	return nil, err
    }

    go u.SetMany(cachestr, &datas)
    return &datas, nil
}
{{- end}}

func (u *{{.ModelName}}Exec) Update(ctx context.Context, data *{{.ModelName}}) error {
    err := u.db.WithContext(ctx).Updates(data).Error
	if err != nil {
		return err
	}
	go u.DelCache(ctx, data)
	return nil
}

func (u *{{.ModelName}}Exec) Delete(ctx context.Context, {{.Pr}} int64) error {
	err := u.db.WithContext(ctx).Where("{{.GPr}} = ?", {{.Pr}}).Delete(&{{.ModelName}}{}).Error
    if err != nil {
    	return err
    }

    go func() {
    	cachestr := fmt.Sprintf("%s%v", cache{{.ModelName}}{{.Pr}}Prefix, {{.Pr}})
    	datacache := u.Get(ctx, cachestr)
    	if datacache != nil {
    		u.DelCache(ctx, datacache)
    	}
    }()
    return nil
}
{{- end -}}