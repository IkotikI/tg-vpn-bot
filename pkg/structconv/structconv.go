package structconv

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// Make a map[string]interface{} from a struct.
// It affects only 1st layer of a struct. So, underling interface{}
// can be a struct, which won't transform to a map.
func MakeMap(in interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	inJson, _ := json.Marshal(in)
	json.Unmarshal(inJson, &inInterface)
	return inInterface
}

// Compare structs fields by given fields names.
// If all fields are equal, returns nil error.
// Otherwise returns an error, displaying where are difference.
func CompareStructs(want interface{}, got interface{}, fields []string) error {
	wantMap := MakeMap(want)
	gotMap := MakeMap(got)

	for _, field := range fields {
		wantField, ok1 := wantMap[field]
		gotField, ok2 := gotMap[field]
		if ok1 == ok2 && wantField == gotField {
			continue
		} else {
			return fmt.Errorf(`field "%s": want "%v", got "%v"`, field, wantField, gotField)
		}
	}

	return nil
}

/* ----- Reflect Transformations ----- */

// Convert a Go type to a SQL null type.
func MapToSQLNullType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Int64:
		return reflect.TypeOf(sql.NullInt64{})
	case reflect.String:
		return reflect.TypeOf(sql.NullString{})
	case reflect.Bool:
		return reflect.TypeOf(sql.NullBool{})
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			return reflect.TypeOf(sql.NullTime{})
		}
		return CreateSQLNullStructType(reflect.New(t).Interface())
	case reflect.Interface, reflect.Ptr:
		return MapToSQLNullType(t.Elem())
	}

	return t
}

// Returns reflect.Type with converted fields to sql.Null_ of original struct.
func CreateSQLNullStructType(orig interface{}) reflect.Type {
	origVal := reflect.ValueOf(orig)
	origType := origVal.Type()

	if origType.Kind() == reflect.Ptr {
		origType = origType.Elem()
	}

	var fields []reflect.StructField

	for i := 0; i < origType.NumField(); i++ {
		field := origType.Field(i)
		nullType := MapToSQLNullType(field.Type)

		fields = append(fields, reflect.StructField{
			Name: field.Name,
			Type: nullType,
			Tag:  field.Tag,
		})
	}

	return reflect.StructOf(fields)
}

// Returns new struct with converted fields to sql.Null_.
func CreateSQLNullStruct(orig interface{}) interface{} {
	return reflect.New(CreateSQLNullStructType(orig)).Interface()
}

// Convert struct with sql.Null_ fields to destination struct.
func ConvertSQLNullStructToBasic(nullStruct interface{}, targetStruct interface{}) {
	// We need exact struct to work with
	// if reflect.TypeOf(nullStruct).Kind() == reflect.Ptr {
	// 	nullStruct = reflect.ValueOf(nullStruct).Elem()
	// }

	nullVal := reflect.ValueOf(nullStruct).Elem()
	targetVal := reflect.ValueOf(targetStruct).Elem()

	for i := 0; i < nullVal.NumField(); i++ {
		nullField := nullVal.Field(i)
		targetField := targetVal.Field(i)

		switch nullField.Type() {
		case reflect.TypeOf(sql.NullBool{}):
			if nullField.FieldByName("Valid").Bool() {
				targetField.SetBool(nullField.FieldByName("Bool").Bool())
			}
		case reflect.TypeOf(sql.NullInt64{}):
			if nullField.FieldByName("Valid").Bool() {
				targetField.SetInt(nullField.FieldByName("Int64").Int())
			}
		case reflect.TypeOf(sql.NullString{}):
			if nullField.FieldByName("Valid").Bool() {
				targetField.SetString(nullField.FieldByName("String").String())
			}
		case reflect.TypeOf(sql.NullTime{}):
			if nullField.FieldByName("Valid").Bool() {
				targetField.Set(reflect.ValueOf(nullField.FieldByName("Time").Interface()))
			}

		default:
			if targetField.Type() == nullField.Type() {
				targetField.Set(nullField)
			}
		}

	}

}

func PrintStructFields(t interface{}) {
	val := reflect.ValueOf(t).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fmt.Printf("Field: %s, Type: %s\n", field.Name, field.Type)
	}
}

/* ---- Parse map ---- */

// func ParseDefaults(dest interface{}, defaults interface{}) error {
// 	destVal := reflect.ValueOf(dest).Elem()
// 	defaultsVal := reflect.ValueOf(defaults).Elem()

// 	// destType := destVal.Type()
// 	defaultsType := defaultsVal.Type()

// 	for i := 0; i < defaultsVal.NumField(); i++ {
// 		defaultsField := defaultsVal.Field(i)
// 		defaultsFieldType := defaultsType.Field(i)

// 		destField := destVal.FieldByName(defaultsFieldType.Name)
// 		if !destField.IsValid() {
// 			return fmt.Errorf("dest has not field \"%+v\", which provided in defaults", defaultsFieldType.Name)
// 		}

// 		if destField.Type() != defaultsField.Type() {
// 			return fmt.Errorf("dest and defaults have incompatible field types: want %s, got %s", defaultsField.Type(), destField.Type())
// 		}

// 		if isZeroValue(destField) {
// 			if destField.CanSet() {
// 				destField.Set(defaultsField)
// 			}
// 			// } else {
// 			// 	return fmt.Errorf()
// 			// }
// 		}

// 		destField.Set(defaultsVal.Field(i))

// 	}
// 	return nil
// }

func ParseDefaultsStrict[T any](dest, defaults T) (err error) {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return fmt.Errorf("given values is not a pointer of type, but %T", dest)
	}

	destVal := reflect.ValueOf(dest).Elem()
	defaultsVal := reflect.ValueOf(defaults).Elem()

	if destVal.Type().Kind() != reflect.Struct {
		return fmt.Errorf("given values is not a struct, but %T", destVal.Interface())
	}

	// fmt.Printf("iterate over fields of %T\n", dest)

	for i := 0; i < defaultsVal.NumField(); i++ {
		defaultsField := defaultsVal.Field(i)
		destField := destVal.Field(i)

		if !destField.CanSet() {
			continue
		}

		if destField.IsZero() {
			destField.Set(defaultsField)
		}

		switch destField.Type().Kind() {
		case reflect.Struct:
			destFieldAddr := destField.Addr().Interface()
			defaultsFieldAddr := defaultsField.Addr().Interface()
			return ParseDefaultsStrict(destFieldAddr, defaultsFieldAddr)

		case reflect.Interface:
			// Possible issues, because interface can point to aliases of simple types,
			// like Stringer { String() string } -> UserID.String().
			destFieldAddr := reflect.ValueOf(destField).Addr().Interface()
			defaultsFieldAddr := reflect.ValueOf(destField).Addr().Interface()
			return ParseDefaultsStrict(destFieldAddr, defaultsFieldAddr)

		case reflect.Ptr:
			return ParseDefaultsStrict(destField.Interface(), defaultsField.Interface())
		}
	}

	return nil
}

// func isZeroValue(v reflect.Value) bool {
// 	zero := reflect.Zero(v.Type()).Interface()
// 	current := v.Interface()
// 	reflect.ValueOf(v.Interface()).IsZero()
// 	return reflect.DeepEqual(current, zero)
// }

func StructToMap(in interface{}, tag string) (out map[string]interface{}, err error) {
	out = make(map[string]interface{})

	// If pointer, dereference
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Function accepts only structs, got %T.", v)
	}

	vType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := vType.Field(i)
		if vTag := field.Tag.Get(tag); vTag != "" {
			out[vTag] = v.Field(i).Interface()
		}
	}
	return out, nil
}
