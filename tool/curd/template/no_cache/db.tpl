{{- define "db2" -}}
func (u *{{.ModelName}}Exec) Create(ctx context.Context, data *{{.ModelName}}) error {
	err := u.exec.Create(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (u *{{.ModelName}}Exec) FindOne(ctx context.Context, {{.Pr}} int64) (*{{.ModelName}}, error) {
	result, err, _ := group.Do(fmt.Sprintf("%s%v", "{{.Pr}}", {{.Pr}}), func() (interface{}, error) {
		var data {{.ModelName}}
		u.exec.AddWhere("{{.Pr}} = ?", {{.Pr}})
		err := u.exec.Find(ctx, &data)
		if err != nil {
			return nil, err
		}
		if data.{{.Pr}} == 0 {
        	return nil, DBNotFound
        }
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
	result, err, _ := group.Do(fmt.Sprintf("%s%v", "{{$notpr}}", {{$notpr}}), func() (interface{}, error) {
		var datas []{{$outer.ModelName}}
		u.exec.AddWhere("{{$notpr}} = ?", {{$notpr}})
		err := u.exec.Find(ctx, &datas)
		if err != nil {
			return nil, err
		}
		if len(datas) == 0 {
        	return nil, DBNotFound
        }
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
	return nil
}
{{- end -}}