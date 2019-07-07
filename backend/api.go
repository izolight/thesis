package backend

import (
	"encoding/json"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

type hashRequest struct {
	Hashes []string `json:"hashes"`
}

type createResponse struct {
	Msg string `json:"msg"`
}

func uploadHashHandler(w http.ResponseWriter, r *http.Request) {
	req := hashRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	session, err := store.Get(r, "SESSIONID")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	session.Values["hashes"] = req.Hashes
	session.Save(r, w)

	respondJSON(w, http.StatusCreated, createResponse{Msg: "accepted hashes"})
}
