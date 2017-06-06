package monascaclient

import (
	"time"
	"net/url"
	"reflect"
	"fmt"
	"sort"
	"strings"
)

func convertStructToQueryParameters(inputStruct interface{}) (url.Values){
	urlValues := url.Values{}
	values := reflect.ValueOf(inputStruct)
	if values.IsNil() {
		return urlValues
	}
	values = values.Elem()
	typ := values.Type()
	// Loop through the struct
	for i := 0; i < typ.NumField(); i++ {
		currentValue := values.Field(i)
		currentType := typ.Field(i)
		// Get Query Parameter Name
		queryParameterKey := currentType.Tag.Get("queryParameter")
		if currentValue.Kind() == reflect.Ptr {
			if currentValue.IsNil() {
				continue
			}
			currentValue = currentValue.Elem()
		}
		addQueryParameter(currentValue, queryParameterKey, &urlValues)
	}
	return urlValues
}

func addQueryParameter(value reflect.Value, key string, values *url.Values) {
	if value.Kind() == reflect.Bool || value.Kind() == reflect.String || value.Kind() == reflect.Int {
		(*values).Add(key, fmt.Sprint(value.Interface()))
	} else if value.Type() == reflect.TypeOf(time.Time{}) {
		timeValue := value.Interface().(time.Time)
		(*values).Add(key, timeValue.UTC().Format(timeFormat))
	} else if value.Kind() == reflect.Map {
		mapValue := value.Interface().(map[string]string)
		if len(mapValue) > 0 {
			dimensionsSlice := make([]string, 0, len(mapValue))
			for key := range mapValue {
				dimensionsSlice = append(dimensionsSlice, key + ":" + mapValue[key])
			}
			// Make sure dimensions are always in correct order to ensure tests pass
			sort.Strings(dimensionsSlice)
			(*values).Add("dimensions", strings.Join(dimensionsSlice, ","))
		}
	}
}
