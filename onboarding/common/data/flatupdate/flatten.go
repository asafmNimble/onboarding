// Package flatbson provides a function for recursively flattening a Go struct using its BSON tags.
package flatupdate

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
	"time"
)

// Flatten returns a map with keys and values corresponding to the field name
// and values of struct v and its nested structs according to its custom flatupdate tags.
// It iterates over each field recursively and sets fields that are not nil.
//
// The BSON struct tags behave in line with the bsoncodec specification. See:
// https://godoc.org/go.mongodb.org/mongo-driver/bson/bsoncodec#StructTags
// for definitions. The supported tags are name, skip, omitempty, and inline.
//
// Flatten does not flatten structs with unexported fields, e.g. time.Time.
// It returns an error if v is not a struct or a pointer to a struct, or if
// the tags produce duplicate keys.
//
//     type A struct {
//       B *X `flatupdate:"b,omitempty"`
//       C X  `flatupdate:"c"`
//     }
//
//     type X struct { Y string `bson:"y"` }
//
//     Flatten(A{nil, X{"hello"}})
//     // Returns:
//     // map[string]interface{}{"c.y": "hello"}
func Flatten(v interface{}, timeGte bool) (map[string]interface{}, error) {
	val := reflect.ValueOf(v)

	val, ok := asStruct(val)
	if !ok {
		return nil, errors.New("v must be a struct or a pointer to a struct")
	}

	m := make(map[string]interface{})
	if err := flattenFields(val, m, "", timeGte); err != nil {
		return nil, err
	}

	return m, nil
}

// flattenFields recursively adds the values of v's fields to map m.
func flattenFields(v reflect.Value, m map[string]interface{}, p string, timeGte bool) error {
	for i := 0; i < v.NumField(); i++ {
		tags, _ := DefaultStructTagParser(v.Type().Field(i))

		if tags.Skip {
			continue
		}

		field := v.Field(i)
		// Default behaviour if tag doesn't exist is same as omit empty
		if (tags.Non || tags.OmitEmpty) && field.IsZero() {
			continue
		}

		// If the field can marshal itself into a BSON type, or it's a struct with
		// unexported fields, like time.Time, we shouldn't recurse into its fields.
		if _, ok := field.Interface().(bson.ValueMarshaler); !ok {
			if s, ok := asStruct(field); ok && !hasUnexportedField(s) {
				fp := p
				if !tags.Inline {
					fp = p + tags.Name + "."
				}
				if err := flattenFields(s, m, fp, false); err != nil {
					return err
				}
				continue
			}
		}

		key := p + tags.Name
		if _, ok := m[key]; ok {
			return fmt.Errorf("duplicated key %s", key)
		}

		// handle the special case of time.Time (having it be a $gte field)
		if field.Type().String() == "time.Time" && timeGte {
			m[key] = map[string]time.Time{"$gte": field.Interface().(time.Time)}
		} else {
			m[key] = field.Interface()
		}
	}

	return nil
}

// asStruct returns that value of v as a struct.
// 	- If v is already a struct, it is returned immediately.
// 	- If v is a pointer, it is dereferenced till a struct is found.
// 	- If a non-struct value is found, it returns false.
func asStruct(v reflect.Value) (reflect.Value, bool) {
	for {
		switch v.Kind() {
		case reflect.Struct:
			return v, true
		case reflect.Ptr:
			v = v.Elem()
		default:
			return reflect.Value{}, false
		}
	}
}

// hasUnexportedField returns true if struct s has a field
// that is not exported.
func hasUnexportedField(s reflect.Value) bool {
	for i := 0; i < s.NumField(); i++ {
		if !s.Field(i).CanInterface() {
			return true
		}
	}
	return false
}

type StructTags struct {
	Name      string
	OmitEmpty bool
	MinSize   bool
	Truncate  bool
	Inline    bool
	Skip      bool
	Non       bool
}

func DefaultStructTagParser(sf reflect.StructField) (StructTags, error) {
	key := strings.ToLower(sf.Name)
	tag, ok := sf.Tag.Lookup("flatupdate")
	if !ok && !strings.Contains(string(sf.Tag), ":") && len(sf.Tag) > 0 {
		tag = string(sf.Tag)
	}

	// Lookup at the bson tag, to extract field name
	bsontag, ok := sf.Tag.Lookup("bson")
	if !ok && !strings.Contains(string(sf.Tag), ":") && len(sf.Tag) > 0 {
		bsontag = string(sf.Tag)
	}
	fieldName := strings.Split(bsontag, ",")[0]
	if fieldName != "" {
		key = fieldName
	}
	var st StructTags
	st.Name = key

	if tag == "-" {
		st.Skip = true
		return st, nil
	}

	if tag == "" {
		st.Non = true
		return st, nil
	}

	for idx, str := range strings.Split(tag, ",") {
		if idx == 0 && str != "" {
			key = str
		}
		switch str {
		case "omitempty":
			st.OmitEmpty = true
		case "minsize":
			st.MinSize = true
		case "truncate":
			st.Truncate = true
		case "inline":
			st.Inline = true
		}
	}

	return st, nil
}
