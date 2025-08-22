{{- define "cache" -}}
// cache
func (u *{{.ModelName}}Exec) DelCache(ctx context.Context, model *{{.ModelName}}) {
	err := u.cacheExec.DeleteCache(fmt.Sprintf("%s%v", cache{{.ModelName}}{{.Pr}}Prefix, model.{{.Pr}}), ctx)
	if err != nil {
		u.logger.Error("Primary key cache delete failure", err)
	}
	{{- $outer := . -}}
    {{- range $notpr := .NotPrs}}
    err = u.cacheExec.DeleteCache(fmt.Sprintf("%s%v", cache{{$outer.ModelName}}{{$notpr}}Prefix, model.{{$notpr}}), ctx)
    if err != nil {
        u.logger.Warn("Not-primary key cache delete failure", err)
    }
    {{- end}}
}

func (u *{{.ModelName}}Exec) Get(ctx context.Context, cachestr string) *{{.ModelName}} {
	var data {{.ModelName}}
	cacheval, err := u.cacheExec.GetCache(cachestr, ctx)
	if err == nil {
		err := UnMarshalJSON(cacheval, &data)
		if err != nil {
			u.logger.Warn("UnMarshal failure: ", err)
			return nil
		}
		return &data
	}
	if !errors.Is(err, CacheNotFound) {
		u.logger.Warn("Primary key cache get failure: ", err)
		return nil
	}
	return nil
}

func (u *{{.ModelName}}Exec) GetMany(ctx context.Context, cachestr string) *[]{{.ModelName}} {
	var datas []{{.ModelName}}
	cacheval, err := u.cacheExec.GetCache(cachestr, ctx)
	if err == nil {
		var key []int64
		err := UnMarshalString(cacheval, &key)
		if err != nil {
			u.logger.Warn("UnMarshal failure: ", err)
			return nil
		}
		for _, c := range key {
			data, err := u.FindOne(ctx, c)
			if err != nil {
				return nil
			}
			datas = append(datas, *data)
		}
		return &datas
	}
	if !errors.Is(err, CacheNotFound) {
		u.logger.Warn("Not-primary key cache get failure: ", err)
		return nil
	}
	return nil
}

func (u *{{.ModelName}}Exec) Set(cachestr string, data *{{.ModelName}}) {
	ctx, cancel := context.WithTimeout(context.Background(), u.cacheExec.SetTTl)
	err := u.cacheExec.SetCache(cachestr, ctx, data)
	if err != nil {
		u.logger.Warn("Primary key cache set failure: ", err)
	}
	cancel()
}

func (u *{{.ModelName}}Exec) SetMany(cachestr string, data *[]{{.ModelName}}) {
	ctx, cancel := context.WithTimeout(context.Background(), u.cacheExec.SetTTl)
	var key []int64
	for _, v := range *data {
		key = append(key, v.{{.Pr}})
		cachestr := fmt.Sprintf("%s%v", cache{{.ModelName}}IdPrefix, v.{{.Pr}})
		u.Set(cachestr, &v)
	}
	err := u.cacheExec.SetCache(cachestr, ctx, &key)
	if err != nil {
		u.logger.Warn("Not-primary key cache set failure: ", err)
	}
	cancel()
}

{{- end -}}