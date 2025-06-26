// internal/handlers/breed_handlers.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/internal/repository"
)

// BreedHandler gère les requêtes HTTP pour les races
type BreedHandler struct {
	repo   repository.BreedRepositoryInterface
	logger *charmLog.Logger
}

// NewBreedHandler crée un nouveau handler
func NewBreedHandler(repo repository.BreedRepositoryInterface, logger *charmLog.Logger) *BreedHandler {
	return &BreedHandler{
		repo:   repo,
		logger: logger,
	}
}

// CreateBreedRequest représente la requête pour créer une race
type CreateBreedRequest struct {
	Species                   string `json:"species"`
	PetSize                   string `json:"pet_size"`
	Name                      string `json:"name"`
	AverageMaleAdultWeight    int    `json:"average_male_adult_weight"`
	AverageFemaleAdultWeight  int    `json:"average_female_adult_weight"`
}

// UpdateBreedRequest représente la requête pour modifier une race
type UpdateBreedRequest struct {
	Species                   string `json:"species,omitempty"`
	PetSize                   string `json:"pet_size,omitempty"`
	Name                      string `json:"name,omitempty"`
	AverageMaleAdultWeight    int    `json:"average_male_adult_weight,omitempty"`
	AverageFemaleAdultWeight  int    `json:"average_female_adult_weight,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// GetAllBreeds récupère toutes les races avec filtres optionnels
// GET /breeds?species=dog&weight_min=5000&weight_max=10000&pet_size=small&limit=10&offset=0
func (h *BreedHandler) GetAllBreeds(w http.ResponseWriter, r *http.Request) {
	// Récupérer les paramètres de requête
	species := r.URL.Query().Get("species")
	petSize := r.URL.Query().Get("pet_size")
	
	var weightMin, weightMax *int
	if weightMinStr := r.URL.Query().Get("weight_min"); weightMinStr != "" {
		if val, err := strconv.Atoi(weightMinStr); err == nil {
			weightMin = &val
		}
	}
	if weightMaxStr := r.URL.Query().Get("weight_max"); weightMaxStr != "" {
		if val, err := strconv.Atoi(weightMaxStr); err == nil {
			weightMax = &val
		}
	}
	
	limit := 50 // valeur par défaut
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}
	
	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			offset = val
		}
	}
	
	breeds, err := h.repo.GetAll(species, weightMin, weightMax, petSize, limit, offset)
	if err != nil {
		h.logger.Error("Erreur lors de la récupération des races", "error", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Erreur serveur", err.Error())
		return
	}
	
	h.sendSuccessResponse(w, http.StatusOK, breeds, "")
}

// GetBreedByID récupère une race par son ID
// GET /breeds/{id}
func (h *BreedHandler) GetBreedByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "ID invalide", "L'ID doit être un nombre entier")
		return
	}
	
	breed, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.Error("Erreur lors de la récupération de la race", "id", id, "error", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Erreur serveur", err.Error())
		return
	}
	
	if breed == nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Race non trouvée", "")
		return
	}
	
	h.sendSuccessResponse(w, http.StatusOK, breed, "")
}

// CreateBreed crée une nouvelle race
// POST /breeds
func (h *BreedHandler) CreateBreed(w http.ResponseWriter, r *http.Request) {
	var req CreateBreedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Corps de requête invalide", err.Error())
		return
	}
	
	// Validation basique
	if req.Species == "" || req.Name == "" || req.PetSize == "" {
		h.sendErrorResponse(w, http.StatusBadRequest, "Champs obligatoires manquants", "species, name et pet_size sont requis")
		return
	}
	
	if req.AverageMaleAdultWeight <= 0 || req.AverageFemaleAdultWeight <= 0 {
		h.sendErrorResponse(w, http.StatusBadRequest, "Poids invalides", "Les poids doivent être supérieurs à 0")
		return
	}
	
	breed := &repository.Breed{
		Species:                   req.Species,
		PetSize:                   req.PetSize,
		Name:                      req.Name,
		AverageMaleAdultWeight:    req.AverageMaleAdultWeight,
		AverageFemaleAdultWeight:  req.AverageFemaleAdultWeight,
	}
	
	createdBreed, err := h.repo.Create(breed)
	if err != nil {
		h.logger.Error("Erreur lors de la création de la race", "breed", req, "error", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la création", err.Error())
		return
	}
	
	h.sendSuccessResponse(w, http.StatusCreated, createdBreed, "Race créée avec succès")
}

// UpdateBreed met à jour une race existante
// PUT /breeds/{id}
func (h *BreedHandler) UpdateBreed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "ID invalide", "L'ID doit être un nombre entier")
		return
	}
	
	var req UpdateBreedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Corps de requête invalide", err.Error())
		return
	}
	
	// Vérifier que la race existe
	existingBreed, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.Error("Erreur lors de la vérification de la race", "id", id, "error", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Erreur serveur", err.Error())
		return
	}
	
	if existingBreed == nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Race non trouvée", "")
		return
	}
	
	// Créer l'objet breed pour la mise à jour
	breed := &repository.Breed{
		Species:                   req.Species,
		PetSize:                   req.PetSize,
		Name:                      req.Name,
		AverageMaleAdultWeight:    req.AverageMaleAdultWeight,
		AverageFemaleAdultWeight:  req.AverageFemaleAdultWeight,
	}
	
	updatedBreed, err := h.repo.Update(id, breed)
	if err != nil {
		h.logger.Error("Erreur lors de la mise à jour de la race", "id", id, "error", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la mise à jour", err.Error())
		return
	}
	
	h.sendSuccessResponse(w, http.StatusOK, updatedBreed, "Race mise à jour avec succès")
}

// DeleteBreed supprime une race
// DELETE /breeds/{id}
func (h *BreedHandler) DeleteBreed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "ID invalide", "L'ID doit être un nombre entier")
		return
	}
	
	// Vérifier que la race existe
	existingBreed, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.Error("Erreur lors de la vérification de la race", "id", id, "error", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Erreur serveur", err.Error())
		return
	}
	
	if existingBreed == nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Race non trouvée", "")
		return
	}
	
	err = h.repo.Delete(id)
	if err != nil {
		h.logger.Error("Erreur lors de la suppression de la race", "id", id, "error", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Erreur lors de la suppression", err.Error())
		return
	}
	
	h.sendSuccessResponse(w, http.StatusOK, nil, "Race supprimée avec succès")
}

// Fonctions utilitaires pour les réponses HTTP
func (h *BreedHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   error,
		Message: message,
	})
}

func (h *BreedHandler) sendSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{
		Data:    data,
		Message: message,
	})
}