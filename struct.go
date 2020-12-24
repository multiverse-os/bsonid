package bsonid

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"
)

func Hash(c interface{}) string {
	return fmt.Sprintf("%x", Sha3(c))
}

func Sha3(c interface{}) string {
	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	h := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(
		h,
		[]byte(fmt.Sprintf("%v", c)),
	)
	return fmt.Sprintf("%x", h)
}

// Dump takes a data structure and returns its byte representation. This can be
// useful if you need to use your own hashing function or formatter.
func Dump(c interface{}) []byte {
	return serialize(c)
}

type item struct {
	name  string
	value reflect.Value
}

type itemSorter []item
type tagError string

type structFieldFilter func(reflect.StructField, *item) (bool, error)

func (self itemSorter) Len() int           { return len(self) }
func (self itemSorter) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self itemSorter) Less(i, j int) bool { return self[i].name < self[j].name }
func (self tagError) Error() string        { return "incorrect tag " + string(self) }

func writeValue(buf *bytes.Buffer, val reflect.Value, fltr structFieldFilter) {
	switch val.Kind() {
	case reflect.String:
		buf.WriteByte('"')
		buf.WriteString(val.String())
		buf.WriteByte('"')
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(strconv.FormatInt(val.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(strconv.FormatUint(val.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(strconv.FormatFloat(val.Float(), 'E', -1, 64))
	case reflect.Bool:
		if val.Bool() {
			buf.WriteByte('t')
		} else {
			buf.WriteByte('f')
		}
	case reflect.Ptr:
		if !val.IsNil() || val.Type().Elem().Kind() == reflect.Struct {
			writeValue(buf, reflect.Indirect(val), fltr)
		} else {
			writeValue(buf, reflect.Zero(val.Type().Elem()), fltr)
		}
	case reflect.Array, reflect.Slice:
		buf.WriteByte('[')
		len := val.Len()
		for i := 0; i < len; i++ {
			if i != 0 {
				buf.WriteByte(',')
			}
			writeValue(buf, val.Index(i), fltr)
		}
		buf.WriteByte(']')
	case reflect.Map:
		mk := val.MapKeys()
		items := make([]item, len(mk), len(mk))
		// Get all values
		for i, _ := range items {
			items[i].name = formatValue(mk[i], fltr)
			items[i].value = val.MapIndex(mk[i])
		}

		// Sort values by key
		sort.Sort(itemSorter(items))

		buf.WriteByte('[')
		for i, _ := range items {
			if i != 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(items[i].name)
			buf.WriteByte(':')
			writeValue(buf, items[i].value, fltr)
		}
		buf.WriteByte(']')
	case reflect.Struct:
		vtype := val.Type()
		flen := vtype.NumField()
		items := make([]item, 0, flen)
		// Get all fields
		for i := 0; i < flen; i++ {
			field := vtype.Field(i)
			it := item{field.Name, val.Field(i)}
			if fltr != nil {
				ok, err := fltr(field, &it)
				if err != nil && strings.Contains(err.Error(), "method:") {
					panic(err)
				}
				if !ok {
					continue
				}
			}
			items = append(items, it)
		}
		// Sort fields by name
		sort.Sort(itemSorter(items))

		buf.WriteByte('{')
		for i, _ := range items {
			if i != 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(items[i].name)
			buf.WriteByte(':')
			writeValue(buf, items[i].value, fltr)
		}
		buf.WriteByte('}')
	case reflect.Interface:
		if !val.CanInterface() {
			return
		}
		writeValue(buf, reflect.ValueOf(val.Interface()), fltr)
	default:
		buf.WriteString(val.String())
	}
}

func formatValue(val reflect.Value, fltr structFieldFilter) string {
	if val.Kind() == reflect.String {
		return "\"" + val.String() + "\""
	}

	var buf bytes.Buffer
	writeValue(&buf, val, fltr)

	return string(buf.Bytes())
}

func filterField(f reflect.StructField, i *item) (bool, error) {
	if str := f.Tag.Get("hash"); str != "" {
		if str == "-" {
			return false, nil
		}
		for _, tag := range strings.Split(str, " ") {
			args := strings.Split(strings.TrimSpace(tag), ":")
			if len(args) != 2 {
				return false, tagError(tag)
			}
			switch args[0] {
			case "name":
				i.name = args[1]
			case "method":
				property, found := f.Type.MethodByName(strings.TrimSpace(args[1]))
				if !found || property.Type.NumOut() != 1 {
					return false, tagError(tag)
				}
				i.value = property.Func.Call([]reflect.Value{i.value})[0]
			}
		}
	}
	return true, nil
}

func serialize(object interface{}) []byte {
	var buf bytes.Buffer

	writeValue(&buf, reflect.ValueOf(object),
		func(f reflect.StructField, i *item) (bool, error) {
			return filterField(f, i)
		})

	return buf.Bytes()
}
