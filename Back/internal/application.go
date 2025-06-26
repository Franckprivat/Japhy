package internal

import (
	"database/sql"
	"fmt"
	"net/http"

	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/internal/handlers"
	"github.com/japhy-tech/backend-test/internal/repository"
	"github.com/japhy-tech/backend-test/internal/service"
)

type App struct {
	logger      *charmLog.Logger
	db          *sql.DB
	breedRepo   *repository.BreedRepository
	breedHandler *handlers.BreedHandler
	csvService  *service.CSVService
}

func NewApp(logger *charmLog.Logger, db *sql.DB) *App {
	breedRepo := repository.NewBreedRepository(db)
	
	csvService := service.NewCSVService()
	
	breedHandler := handlers.NewBreedHandler(breedRepo, logger)
	
	return &App{
		logger:       logger,
		db:           db,
		breedRepo:    breedRepo,
		breedHandler: breedHandler,
		csvService:   csvService,
	}
}

func (a *App) RegisterRoutes(r *mux.Router) {
	// Routes pour les races
	r.HandleFunc("/breeds", a.breedHandler.GetAllBreeds).Methods(http.MethodGet)
	r.HandleFunc("/breeds", a.breedHandler.CreateBreed).Methods(http.MethodPost)
	r.HandleFunc("/breeds/{id:[0-9]+}", a.breedHandler.GetBreedByID).Methods(http.MethodGet)
	r.HandleFunc("/breeds/{id:[0-9]+}", a.breedHandler.UpdateBreed).Methods(http.MethodPut)
	r.HandleFunc("/breeds/{id:[0-9]+}", a.breedHandler.DeleteBreed).Methods(http.MethodDelete)
	r.HandleFunc("/import-breeds", a.ImportBreedsFromCSV).Methods(http.MethodPost)
}

func (a *App) ImportBreedsFromCSV(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("Début de l'import des races depuis le CSV")
	
	// Pour lire les races depuis le fichier CSV
	breeds, err := a.csvService.ReadBreedsFromCSV("./breeds.csv")
	if err != nil {
		a.logger.Error("Erreur lors de la lecture du CSV", "error", err)
		http.Error(w, "Erreur lors de la lecture du fichier CSV: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	a.logger.Info("Races lues depuis le CSV", "count", len(breeds))
	
	// Afin d'importer les races dans la base de données
	err = a.breedRepo.ImportFromCSV(breeds)
	if err != nil {
		a.logger.Error("Erreur lors de l'import en base", "error", err)
		http.Error(w, "Erreur lors de l'import en base de données: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	a.logger.Info("Import des races terminé avec succès", "count", len(breeds))
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf(`{"message": "Import des races terminé avec succès", "count": %d}`, len(breeds))
	w.Write([]byte(response))
}