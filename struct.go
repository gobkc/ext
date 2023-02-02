package gext

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func BindDefault(dest interface{}) error {
	t := reflect.TypeOf(dest)
	if dt := t.Kind(); dt != reflect.Ptr {
		return errors.New("dest must be a struct pointer")
	}
	if dt := t.Elem().Kind(); dt != reflect.Struct {
		return errors.New("dest must be a struct pointer")
	}
	v := reflect.ValueOf(dest).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag
		df := tag.Get("default")
		if fmt.Sprintf("%v", v.Field(i).Interface()) != "" {
			continue
		}
		switch fk := field.Type.Kind(); fk {
		case reflect.String:
			v.Field(i).SetString(df)
		case reflect.Int:
			val, err := strconv.Atoi(df)
			if err != nil {
				return err
			}
			v.Field(i).Set(reflect.ValueOf(val))
		case reflect.Int64:
			val, err := strconv.ParseInt(df, 10, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetInt(val)
		case reflect.Int32:
			val, err := strconv.ParseInt(df, 10, 32)
			if err != nil {
				return err
			}
			newV := int32(val)
			v.Field(i).Set(reflect.ValueOf(newV))
		case reflect.Float32:
			val, err := strconv.ParseFloat(df, 32)
			if err != nil {
				return err
			}
			newVal := float32(val)
			v.Field(i).Set(reflect.ValueOf(newVal))
		case reflect.Float64:
			val, err := strconv.ParseFloat(df, 64)
			if err != nil {
				return err
			}
			v.Field(i).SetFloat(val)
		case reflect.Bool:
			var val bool
			if df = strings.ToUpper(df); df == "TRUE" {
				val = true
			}
			v.Field(i).SetBool(val)
		default:
			return errors.New("unsupported type")
		}
	}
	return nil
}

func TrimAll(data any) (err error) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Ptr:
		switch reflect.ValueOf(data).Elem().Kind() {
		case reflect.String:
			old := data.(*string)
			reflect.ValueOf(data).Elem().SetString(strings.TrimSpace(*old))
			return
		case reflect.Struct:
			for idx := 0; idx < reflect.ValueOf(data).Elem().NumField(); idx++ {
				if fKind := reflect.ValueOf(data).Elem().Field(idx).Kind(); fKind == reflect.String {
					oldStr := reflect.ValueOf(data).Elem().Field(idx).String()
					newStr := strings.TrimSpace(oldStr)
					reflect.ValueOf(data).Elem().Field(idx).SetString(newStr)
				}
			}
			return
		}
	}
	return errors.New(`dest must be a string/struct pointer`)
}
