package gher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// Gher is a generic handler that takes a function and returns a http.Handler.
// The function must take an input and a http.Request and return an output and an error.
// The input and output can be any type, as long as they can be marshaled and unmarshalled
// to and from JSON, are strings, or that the output implements the io.Reader interface.
func Gher[I, O any](next func(I, *http.Request) (O, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := new(I)

		if _, ok := reflect.ValueOf(*i).Interface().(string); ok {
			buffer := new(bytes.Buffer)
			_, err := buffer.ReadFrom(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "failed to read request"}`))
				return
			}
			reflect.ValueOf(i).Elem().SetString(buffer.String())
		} else {
			err := json.NewDecoder(r.Body).Decode(i)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "failed to parse request"}`))
				return
			}
		}

		o, err := next(*i, r)
		if err != nil {
			e := fmt.Sprintf(`{"error": %q}`, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(e))
			return
		}

		if oreader, ok := reflect.ValueOf(o).Interface().(io.Reader); ok {
			w.WriteHeader(http.StatusOK)
			io.Copy(w, oreader)
			return
		} else if ostring, ok := reflect.ValueOf(o).Interface().(string); ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(ostring))
			return
		} else {
			err = json.NewEncoder(w).Encode(o)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "failed to encode response"}`))
				return
			}
		}
	})
}
