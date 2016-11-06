// Package orm is just a simple orm.
// Created by simplejia [9/2016]
package orm

import (
	"reflect"
	"strings"

	"database/sql"
)

var TagName = "orm"

func getFieldInfo(typ reflect.Type) map[string][]int {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	finfo := make(map[string][]int)

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		tag := f.Tag.Get(TagName)

		// Skip unexported fields or fields marked with "-"
		if f.PkgPath != "" || tag == "-" {
			continue
		}

		// Handle embedded structs
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			for k, v := range getFieldInfo(f.Type) {
				finfo[k] = append(f.Index, v...)
			}
			continue
		}

		// Use field name for untagged fields
		if tag == "" {
			tag = f.Name
		}

		tag = strings.ToLower(tag)

		finfo[tag] = f.Index
	}

	return finfo
}

func Rows2Strus(rows *sql.Rows, strus interface{}) (err error) {
	columns, err := rows.Columns()
	if err != nil {
		return
	}

	strusRV := reflect.Indirect(reflect.ValueOf(strus))
	elemRT := strusRV.Type().Elem()

	fieldInfo := getFieldInfo(elemRT)

	for rows.Next() {
		var struRV reflect.Value
		var struField reflect.Value
		if elemRT.Kind() == reflect.Ptr {
			struRV = reflect.New(elemRT.Elem())
			struField = reflect.Indirect(struRV)
		} else {
			struRV = reflect.Indirect(reflect.New(elemRT))
			struField = struRV
		}
		var values []interface{}
		for _, column := range columns {
			idx, ok := fieldInfo[strings.ToLower(column)]
			var v interface{}
			if !ok {
				var i interface{}
				v = &i
			} else {
				v = struField.FieldByIndex(idx).Addr().Interface()
			}
			values = append(values, v)
		}
		err = rows.Scan(values...)
		if err != nil {
			return
		}
		strusRV = reflect.Append(strusRV, struRV)
	}
	if err = rows.Err(); err != nil {
		return
	}
	reflect.Indirect(reflect.ValueOf(strus)).Set(strusRV)

	return
}

func Rows2Stru(rows *sql.Rows, stru interface{}) (err error) {
	struRT := reflect.TypeOf(stru).Elem()

	strusPtrRV := reflect.New(reflect.SliceOf(struRT))
	err = Rows2Strus(rows, strusPtrRV.Interface())
	if err != nil {
		return
	}
	strusRV := reflect.Indirect(strusPtrRV)
	if strusRV.Len() == 0 {
		err = sql.ErrNoRows
		return
	}
	reflect.Indirect(reflect.ValueOf(stru)).Set(strusRV.Index(0))
	return
}

func Rows2Cnts(rows *sql.Rows, cnts interface{}) (err error) {
	cntsRV := reflect.Indirect(reflect.ValueOf(cnts))
	elemRT := cntsRV.Type().Elem()

	for rows.Next() {
		var values []interface{}
		var cntRV reflect.Value
		if elemRT.Kind() == reflect.Ptr {
			cntRV = reflect.New(elemRT.Elem())
			values = append(values, cntRV.Interface())
		} else {
			cntRV = reflect.Indirect(reflect.New(elemRT))
			values = append(values, cntRV.Addr().Interface())
		}
		err = rows.Scan(values...)
		if err != nil {
			return
		}
		cntsRV = reflect.Append(cntsRV, cntRV)
	}
	if err = rows.Err(); err != nil {
		return
	}
	reflect.Indirect(reflect.ValueOf(cnts)).Set(cntsRV)

	return
}

func Rows2Cnt(rows *sql.Rows, cnt interface{}) (err error) {
	cntRT := reflect.TypeOf(cnt).Elem()

	cntsPtrRV := reflect.New(reflect.SliceOf(cntRT))
	err = Rows2Cnts(rows, cntsPtrRV.Interface())
	if err != nil {
		return
	}
	cntsRV := reflect.Indirect(cntsPtrRV)
	if cntsRV.Len() == 0 {
		err = sql.ErrNoRows
		return
	}
	reflect.Indirect(reflect.ValueOf(cnt)).Set(cntsRV.Index(0))
	return
}
