{{- define "cache" -}}
// 序列化
func UnMarshalJSON_{{.ModelName}} (s string, model *{{.ModelName}}) error {
    return json.Unmarshal([]byte(s), model)
}

func UnMarshalString_{{.ModelName}} (s string, model *[]int64) error {
    return json.Unmarshal([]byte(s), model)
}

// cache
func (u *{{.ModelName}}Exec) DelCache(ctx context.Context, model *{{.ModelName}}) {
	err := u.cacheExec.DeleteCache(fmt.Sprintf("%s%v", cache{{.ModelName}}{{.Pr}}Prefix, model.{{.Pr}}), ctx)
	if err != nil {
	    u.logger.Error("主键缓存删除失败", logger.Field{"{{.Pr}}": model.{{.Pr}}, "error": err})
	}
	{{- $outer := . -}}
    {{- range $notpr := .NotPrs}}
    err = u.cacheExec.DeleteCache(fmt.Sprintf("%s%v", cache{{$outer.ModelName}}{{$notpr.Name}}Prefix, model.{{$notpr.Name}}), ctx)
    if err != nil {
        u.logger.Warn("非主键缓存删除失败", logger.Field{"{{$notpr.Name}}": model.{{$notpr.Name}}, "error": err})
    }
    {{- end}}
}

func (u *{{.ModelName}}Exec) Get(ctx context.Context, cachestr string) *{{.ModelName}} {
	var data {{.ModelName}}
	cacheval, err := u.cacheExec.GetCache(cachestr, ctx)
	if err == nil {
		err := UnMarshalJSON_{{.ModelName}}(cacheval, &data)
		if err != nil {
			u.logger.Warn("Json 序列化出错", logger.Field{"error": err})
			return nil
		}
		return &data
	}
	if !errors.Is(err, CacheNotFound) {
		u.logger.Warn("主键缓存获取失败", logger.Field{"error": err})
		return nil
	}
	return nil
}

func (u *{{.ModelName}}Exec) GetMany(ctx context.Context, cachestr string) *[]{{.ModelName}} {
	var datas []{{.ModelName}}
	cacheval, err := u.cacheExec.GetCache(cachestr, ctx)
	if err == nil {
		var key []int64
		err := UnMarshalString_{{.ModelName}}(cacheval, &key)
		if err != nil {
			u.logger.Warn("Json 序列化出错", logger.Field{"error": err})
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
		u.logger.Warn("非主键缓存获取失败", logger.Field{"error": err})
		return nil
	}
	return nil
}

func (u *{{.ModelName}}Exec) Set(cachestr string, data *{{.ModelName}}) {
	ctx, cancel := context.WithTimeout(context.Background(), u.cacheExec.SetTTl)
	err := u.cacheExec.SetCache(cachestr, ctx, data)
	if err != nil {
		u.logger.Warn("主键缓存设置失败", logger.Field{"error": err})
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
		u.logger.Warn("非主键缓存设置失败", logger.Field{"error": err})
	}
	cancel()
}

{{- end -}}