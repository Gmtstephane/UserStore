package userStoreAPi

import (
	"encoding/json"
	"net/http"
)

//internalErr encode an 500 error
func internalErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(err.Error())
}
