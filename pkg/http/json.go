package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// JSONB should return json key value as string
func JSONB(key, val string) []byte {
	b, _ := json.Marshal(map[string]string{key: val})
	return b
}

// JSON should return json key value as string
func EncodeJSONKV(key, val string) string {
	b, _ := json.Marshal(map[string]string{key: val})
	return string(b)
}

// EncodeJSON should marshal object to string
func EncodeJSON(o interface{}) []byte {
	b, _ := json.Marshal(o)
	return b
}

func JSON(w http.ResponseWriter, o interface{}) {
	_, _ = w.Write(EncodeJSON(o))
}

func OK(w http.ResponseWriter) {
	_, _ = w.Write(JSONB("response", "OK"))
}

func NotImplemented(w http.ResponseWriter) {
	_, _ = w.Write(JSONB("response", "not_implemented"))
}

func Read(r *http.Request, out interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, out)
	if err != nil {
		return err
	}

	return nil
}

// Error with bad request code and message
func Error(w http.ResponseWriter, err error) {
	http.Error(w, EncodeJSONKV("error", err.Error()), http.StatusBadRequest)
}
