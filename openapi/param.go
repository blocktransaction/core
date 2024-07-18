package openapi

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// json转成url.values
func json2UrlValues(data interface{}) url.Values {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	values := url.Values{}
	for i := 0; i < v.NumField(); i++ {
		//-和sign字段不参与签名
		if key := t.Field(i).Tag.Get("json"); key != "-" && key != "sign" {
			switch v.Field(i).Kind() {
			case reflect.String:
				if val := v.Field(i).String(); val != "" {
					values.Set(key, val)
				}
			case reflect.Float64:
				if val := v.Field(i).Float(); val > 0 {
					values.Set(key, strconv.FormatFloat(val, 'f', -1, 64))
				}
			case reflect.Int:
				if val := v.Field(i).Int(); val > 0 {
					values.Set(key, strconv.Itoa(int(val)))
				}
			case reflect.Int64:
				if val := v.Field(i).Int(); val > 0 {
					//json为string，但类型为int64的
					if index := strings.Index(key, ",string"); index > -1 {
						key = key[:index]
					}
					values.Set(key, strconv.FormatInt(val, 10))
				}
			}
		}
	}
	return values
}

// map to urlvalues
func Map2UrlValues(data map[string]string) url.Values {
	values := url.Values{}
	for k, v := range data {
		values.Set(k, v)
	}
	return values
}

// json to map
func Json2Map(data interface{}) map[string]string {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	m := make(map[string]string, 0)
	for i := 0; i < v.NumField(); i++ {
		//不包含含有'-'字段
		if t.Field(i).Tag.Get("json") != "-" {
			switch v.Field(i).Kind() {
			case reflect.String:
				if vv := v.Field(i).String(); vv != "" {
					if t.Field(i).Tag.Get("json") != "sign" {
						m[t.Field(i).Tag.Get("json")] = vv
					}
				}
			case reflect.Int:
				if vv := v.Field(i).Int(); vv >= 0 {
					m[t.Field(i).Tag.Get("json")] = strconv.Itoa(int(vv))
				}
			case reflect.Float64:
				if vv := v.Field(i).Float(); vv >= 0 {
					m[t.Field(i).Tag.Get("json")] = strconv.FormatFloat(vv, 'f', -1, 64)
				}
			}
		}
	}
	return m
}
