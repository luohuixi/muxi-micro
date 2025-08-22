{{- define "db" -}}
func (u *{{.ModelName}}Exec) Create(ctx context.Context, data *{{.ModelName}}) error {
	err := u.exec.Create(ctx, data)
	if err != nil {
		return err
	}
	go u.DelCache(ctx, data)
	return nil
}

func (u *{{.ModelName}}Exec) FindOne(ctx context.Context, id int64) (*{{.ModelName}}, error) {
	cachestr := fmt.Sprintf("%s%v", cache{{.ModelName}}{{.Pr}}Prefix, id)
	result, err, _ := group.Do(cachestr, func() (interface{}, error) {
		datacache := u.Get(ctx, cachestr)
		if datacache != nil {
			return datacache, nil
		}
		var data {{.ModelName}}
		u.exec.AddWhere("id = ?", id)
		err := u.exec.Find(ctx, &data)
		if err != nil {
			return nil, err
		}
		if data.{{.Pr}} == 0 {
        	return nil, DBNotFound
        }
		go u.Set(cachestr, &data)
		return &data, nil
	})
	if result == nil {
        return nil, err
    }
	return result.(*{{.ModelName}}), err
}

{{- $outer := . -}}
{{- range $notpr := .NotPrs}}

func (u *{{$outer.ModelName}}Exec) FindBy{{$notpr}}(ctx context.Context, {{$notpr}} string) (*[]{{$outer.ModelName}}, error) {
	cachestr := fmt.Sprintf("%s%v", cache{{$outer.ModelName}}{{$notpr}}Prefix, {{$notpr}})
	result, err, _ := group.Do(cachestr, func() (interface{}, error) {
		cacheval, err := u.cacheExec.GetCache(cachestr, ctx)
		datascache := u.GetMany(ctx, cacheval)
		if datascache != nil {
			return datascache, nil
		}
		var datas []{{$outer.ModelName}}
		u.exec.AddWhere("{{$notpr}} = ?", {{$notpr}})
		err = u.exec.Find(ctx, &datas)
		if err != nil {
			return nil, err
		}
		if len(datas) == 0 {
        	return nil, DBNotFound
        }
		go u.SetMany(cachestr, &datas)
		return &datas, nil
	})
	if result == nil {
    	return nil, err
    }
	return result.(*[]{{$outer.ModelName}}), err
}
{{- end}}

func (u *{{.ModelName}}Exec) Update(ctx context.Context, data *{{.ModelName}}) error {
	u.exec.AddWhere("id = ?", data.{{.Pr}})
	err := u.exec.Update(ctx, data)
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
	err = u.exec.Delete(ctx, &data)
	if err != nil {
		return err
	}
	go u.DelCache(ctx, &data)
	return nil
}
{{- end -}}