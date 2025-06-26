package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetAll(t *testing.T) {
	// ...
}

func TestGetAll_WithMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erreur lors de la création du mock: %v", err)
	}
	defer db.Close()

	repo := NewBreedRepository(db)

	rows := sqlmock.NewRows([]string{"id", "species", "pet_size", "name", "average_male_adult_weight", "average_female_adult_weight"}).
		AddRow(1, "dog", "medium", "Border Collie", 20, 18)

	mock.ExpectQuery("SELECT id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight FROM breeds").
		WillReturnRows(rows)

	breeds, err := repo.GetAll("", nil, nil, "", 10, 0)
	if err != nil {
		t.Fatalf("Erreur lors de GetAll avec mock: %v", err)
	}
	if len(breeds) != 1 {
		t.Errorf("Attendu 1 résultat, obtenu %d", len(breeds))
	}
	if breeds[0].Name != "Border Collie" {
		t.Errorf("Nom attendu 'Border Collie', obtenu '%s'", breeds[0].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Toutes les attentes du mock n'ont pas été satisfaites: %v", err)
	}
} 