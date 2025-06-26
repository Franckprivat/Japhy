package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"github.com/japhy-tech/backend-test/internal/repository"
)

type CSVService struct{}

func NewCSVService() *CSVService {
	return &CSVService{}
}

func (s *CSVService) ReadBreedsFromCSV(filename string) ([]repository.Breed, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'ouverture du fichier CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du fichier CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("le fichier CSV est vide")
	}

	var breeds []repository.Breed
	for i, record := range records[1:] {
		if len(record) != 6 {
			return nil, fmt.Errorf("ligne %d: nombre de colonnes incorrect (attendu: 6, reçu: %d)", i+2, len(record))
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("ligne %d: ID invalide '%s': %w", i+2, record[0], err)
		}

		maleWeight, err := strconv.Atoi(record[4])
		if err != nil {
			return nil, fmt.Errorf("ligne %d: poids mâle invalide '%s': %w", i+2, record[4], err)
		}

		femaleWeight, err := strconv.Atoi(record[5])
		if err != nil {
			return nil, fmt.Errorf("ligne %d: poids femelle invalide '%s': %w", i+2, record[5], err)
		}

		breed := repository.Breed{
			ID:                        id,
			Species:                   record[1],
			PetSize:                   record[2],
			Name:                      record[3],
			AverageMaleAdultWeight:    maleWeight,
			AverageFemaleAdultWeight:  femaleWeight,
		}

		breeds = append(breeds, breed)
	}

	return breeds, nil
}