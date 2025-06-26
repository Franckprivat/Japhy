package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/japhy-tech/backend-test/internal/repository"
	"github.com/charmbracelet/log"
)

type MockBreedRepo struct{}

func (m *MockBreedRepo) GetAll(species string, weightMin, weightMax *int, petSize string, limit, offset int) ([]repository.Breed, error) {
	return []repository.Breed{
		{ID: 1, Species: "dog", PetSize: "medium", Name: "Border Collie", AverageMaleAdultWeight: 20, AverageFemaleAdultWeight: 18},
	}, nil
}

func (m *MockBreedRepo) Create(breed *repository.Breed) (*repository.Breed, error) { return nil, nil }
func (m *MockBreedRepo) GetByID(id int) (*repository.Breed, error) { return nil, nil }
func (m *MockBreedRepo) Update(id int, breed *repository.Breed) (*repository.Breed, error) { return nil, nil }
func (m *MockBreedRepo) Delete(id int) error { return nil }
func (m *MockBreedRepo) ImportFromCSV(breeds []repository.Breed) error { return nil }

func TestGetAllBreeds(t *testing.T) {
	mockRepo := &MockBreedRepo{}
	logger := log.NewWithOptions(nil, log.Options{})
	handler := NewBreedHandler(mockRepo, logger)

	req := httptest.NewRequest("GET", "/breeds", nil)
	w := httptest.NewRecorder()

	handler.GetAllBreeds(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("attendu 200, obtenu %d", resp.StatusCode)
	}
}
