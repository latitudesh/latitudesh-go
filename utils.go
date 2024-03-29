package latitude

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"regexp"
)

var timestampType = reflect.TypeOf(Timestamp{})

func Stringify(message interface{}) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)

	err := stringifyValue(&buf, v)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

// contains tells whether a contains x.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func stringifyValue(w io.Writer, val reflect.Value) error {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		_, err := w.Write([]byte("<nil>"))
		return err
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		if _, err := fmt.Fprintf(w, `"%s"`, v); err != nil {
			return err
		}
	case reflect.Slice:
		if _, err := w.Write([]byte{'['}); err != nil {
			return err
		}
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				if _, err := w.Write([]byte{' '}); err != nil {
					return err
				}
			}

			if err := stringifyValue(w, v.Index(i)); err != nil {
				return err
			}
		}

		if _, err := w.Write([]byte{']'}); err != nil {
			return err
		}
		return nil
	case reflect.Struct:
		if v.Type().Name() != "" {
			if _, err := w.Write([]byte(v.Type().String())); err != nil {
				return err
			}
		}

		// special handling of Timestamp values
		if v.Type() == timestampType {
			_, err := fmt.Fprintf(w, "{%s}", v.Interface())
			return err
		}

		if _, err := w.Write([]byte{'{'}); err != nil {
			return err
		}

		var sep bool
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Slice && fv.IsNil() {
				continue
			}

			if sep {
				if _, err := w.Write([]byte(", ")); err != nil {
					return err
				}
			} else {
				sep = true
			}

			if _, err := w.Write([]byte(v.Type().Field(i).Name)); err != nil {
				return err
			}
			if _, err := w.Write([]byte{':'}); err != nil {
				return err
			}

			if err := stringifyValue(w, fv); err != nil {
				return err
			}
		}

		if _, err := w.Write([]byte{'}'}); err != nil {
			return err
		}
	default:
		if v.CanInterface() {
			if _, err := fmt.Fprint(w, v.Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}

// validate UUID
func ValidateUUID(uuid string) error {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	if !r.MatchString(uuid) {
		return fmt.Errorf("%s is not a valid UUID", uuid)
	}
	return nil
}
