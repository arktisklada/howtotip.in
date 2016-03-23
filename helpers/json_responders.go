package helpers

import (
  "encoding/json"
  "fmt"
  "net/http"
  "reflect"
)

func jsonResponder(w http.ResponseWriter, r *http.Request, data interface{}, err error) {
  w.Header().Set("Content-Type", "application/json")
  enc := json.NewEncoder(w)
  if err != nil {
    http.Error(w, "", http.StatusBadRequest)
    enc.Encode(map[string]string{"error": err.Error()})
  } else {
    kind := reflect.TypeOf(data).Kind()
    if (kind == reflect.Slice || kind == reflect.Map || kind == reflect.Ptr) && reflect.ValueOf(data).IsNil() {
      if kind == reflect.Slice {
        fmt.Fprintf(w, "[]")
      } else {
        fmt.Fprintf(w, "{}")
      }
    } else {
      enc.Encode(data)
    }
  }
}

func successResponder(w http.ResponseWriter, r *http.Request, err error) {
  jsonResponder(w, r, map[string]bool{"success": err == nil}, err)
}
