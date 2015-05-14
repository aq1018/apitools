package errormapper

import (
	"encoding/json"
	"net/http"
)

func (mapper *Mapper) Register(fn ErrorMappingFunc) {
	mapper.mappingFuncs = append(mapper.mappingFuncs, fn)
}

func (mapper *Mapper) MapError(err error) ErrorMessage {
	for _, fn := range mapper.mappingFuncs {
		if message := fn(err); message != nil {
			return message
		}
	}
	return nil
}

func (mapper *Mapper) Middleware(handler HTTPHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// run the handler
		err := handler(w, r)

		// no error, return early
		if err == nil {
			return
		}

		message := mapper.MapError(err)

		// cannot map this error, writing the error out as text
		if message == nil {
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// write header to http socket with content type as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(message.HTTPStatus())

		// write encoded error message as http body
		err = json.NewEncoder(w).Encode(message)

		// something went wrong with JSON Encoding / writing to http socket
		// we should panic in this case
		if err != nil {
			panic(err)
		}
	})
}

type (
	ErrorMappingFunc func(err error) ErrorMessage

	HTTPHandler func(w http.ResponseWriter, r *http.Request) error

	ErrorMessage interface {
		HTTPStatus() int
	}

	Mapper struct {
		mappingFuncs []ErrorMappingFunc
	}
)
