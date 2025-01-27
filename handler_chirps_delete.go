package main

import (
	"net/http"

	"github.com/Bephoros/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	if userID != dbChirp.UserID {
		respondWithError(w, http.StatusForbidden, "User is not authorized to delete the chirp", err)
		return
	}

	deleted := cfg.db.DeleteChirp(r.Context(), chirpID)
	if deleted != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
