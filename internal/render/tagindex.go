package render

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jheddings/ccglow/internal/types"
)

type fieldAccessor struct {
	Provider      string
	FieldIndex    int
	DefaultFormat string
}

type TagIndex map[string]fieldAccessor

func BuildTagIndex(providers map[string]types.DataProvider) (TagIndex, error) {
	idx := make(TagIndex)

	for name, p := range providers {
		fp, ok := p.(types.FieldProvider)
		if !ok {
			continue
		}

		fields := fp.Fields()
		t := reflect.TypeOf(fields)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			continue
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("segment")
			if tag == "" {
				continue
			}

			segName, defaultFmt := parseSegmentTag(tag)

			if _, exists := idx[segName]; exists {
				return nil, fmt.Errorf("duplicate segment name %q", segName)
			}

			idx[segName] = fieldAccessor{
				Provider:      name,
				FieldIndex:    i,
				DefaultFormat: defaultFmt,
			}
		}
	}

	return idx, nil
}

func parseSegmentTag(tag string) (name, defaultFormat string) {
	parts := strings.SplitN(tag, ",", 2)
	name = parts[0]
	if len(parts) > 1 && strings.HasPrefix(parts[1], "format:") {
		defaultFormat = strings.TrimPrefix(parts[1], "format:")
	}
	return
}

func ResolveSegmentValues(idx TagIndex, providerData map[string]any) map[string]any {
	values := make(map[string]any)

	for segName, accessor := range idx {
		data, ok := providerData[accessor.Provider]
		if !ok || data == nil {
			continue
		}

		v := reflect.ValueOf(data)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				continue
			}
			v = v.Elem()
		}

		field := v.Field(accessor.FieldIndex)

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				values[segName] = nil
				continue
			}
			field = field.Elem()
		}

		values[segName] = field.Interface()
	}

	return values
}
