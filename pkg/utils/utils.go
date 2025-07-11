package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Err  bool `json:"err"`
	Data any  `json:"data"`
}

func WriteJSON(w http.ResponseWriter, data any, status int, headers ...http.Header) error {
	resp := Response{
		Err:  false,
		Data: data,
	}

	out, err := json.Marshal(resp)

	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func WriteError(w http.ResponseWriter, err error, status ...int) {
	code := http.StatusBadRequest

	if len(status) > 0 {
		code = status[0]
	}

	resp := Response{
		Err:  true,
		Data: err.Error(),
	}

	out, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(out)
}

func ReadJSON(r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}
