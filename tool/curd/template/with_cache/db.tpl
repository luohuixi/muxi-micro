{{- define "db" -}}
func (u *{{.ModelName}}Exec) Create(ctx context.Context, data *{{.ModelName}}) error {
	err := u.exec.Create(ctx, data, {{.ModelName}}{})
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
	u.exec.Model.Where("{{.Pr}} = ?", {{.Pr}})
	err := u.exec.Find(ctx, &data, {{.ModelName}}{})
	if err != nil {
		return nil, err
	}
	if data.{{.Pr}} == 0 {
       	return nil, DBNotFound
    }
	go u.Set(cachestr, &data)
	return &data, nil
}

{{- $outer := . -}}
{{- range $notpr := .NotPrs}}

func (u *{{$outer.ModelName}}Exec) FindBy{{$notpr.Name}}(ctx context.Context, {{$notpr.Name}} {{$notpr.Type}}) (*[]{{$outer.ModelName}}, error) {
	cachestr := fmt.Sprintf("%s%v", cache{{$outer.ModelName}}{{$notpr.Name}}Prefix, {{$notpr.Name}})
	cacheval, err := u.cacheExec.GetCache(cachestr, ctx)
	datascache := u.GetMany(ctx, cacheval)
	if datascache != nil {
		return datascache, nil
	}
	var datas []{{$outer.ModelName}}
	u.exec.Model.Where("{{$notpr.Name}} = ?", {{$notpr.Name}})
	err = u.exec.Find(ctx, &datas, {{$outer.ModelName}}{})
	if err != nil {
		return nil, err
	}
	if len(datas) == 0 {
        return nil, DBNotFound
    }
	go u.SetMany(cachestr, &datas)
	return &datas, nil
}
{{- end}}

func (u *{{.ModelName}}Exec) Update(ctx context.Context, data *{{.ModelName}}) error {
	u.exec.Model.Where("id = ?", data.{{.Pr}})
	err := u.exec.Update(ctx, data, {{.ModelName}}{})
	if err != nil {
		return err
	}
	go u.DelCache(ctx, data)
	return nil
}

func (u *{{.ModelName}}Exec) Delete(ctx context.Context, id int64) error {
	var data {{.ModelName}}
	d, err := u.FindOne(ctx, id)
	if err != nil {
		return err
	}
	data = *d
	err = u.exec.Delete(ctx, &data, {{.ModelName}}{})
	if err != nil {
		return err
	}
	go u.DelCache(ctx, &data)
	return nil
}
{{- end -}}