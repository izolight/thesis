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
	req, err := decodeHashRequest(w, r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	}

	session, err := store.Get(r, "SESSIONID")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	session.Values["hashes"] = req.Hashes
	session.Save(r, w)

	respondJSON(w, http.StatusCreated, createResponse{Msg: "accepted hashes"})
}

func decodeHashRequest(w http.ResponseWriter, r *http.Request) (hashRequest, error) {
	req := hashRequest{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&req)
	return req, err
}

func updateHashHandler(w http.ResponseWriter, r *http.Request) {
	req, err := decodeHashRequest(w, r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	session, err := store.Get(r, "SESSIONID")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	existingHashes, ok := session.Values["hashes"].([]string)
	if !ok {
		session.Values["hashes"] = req.Hashes
	} else {
		session.Values["hashes"] = append(existingHashes, req.Hashes...)
	}
	session.Save(r, w)

	respondJSON(w, http.StatusCreated, createResponse{Msg: "accepted hashes"})

}

func dumpHashes(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "SESSIONID")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashes, ok := session.Values["hashes"].([]string)
	if !ok {
		respondError(w, http.StatusBadRequest, "no hashes found")
		return
	}

	req := hashRequest{
		Hashes: hashes,
	}

	respondJSON(w, http.StatusOK, req)
}
