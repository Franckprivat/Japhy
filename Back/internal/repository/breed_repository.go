package repository

import (
	"database/sql"
	"fmt"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)

type Breed struct {
	ID                       int    `json:"id" db:"id"`
	Species                  string `json:"species" db:"species"`
	PetSize                  string `json:"pet_size" db:"pet_size"`
	Name                     string `json:"name" db:"name"`
	AverageMaleAdultWeight   int    `json:"average_male_adult_weight" db:"average_male_adult_weight"`
	AverageFemaleAdultWeight int    `json:"average_female_adult_weight" db:"average_female_adult_weight"`
}

type BreedRepositoryInterface interface {
	GetAll(species string, weightMin, weightMax *int, petSize string, limit, offset int) ([]Breed, error)
	GetByID(id int) (*Breed, error)
	Create(breed *Breed) (*Breed, error)
	Update(id int, breed *Breed) (*Breed, error)
	Delete(id int) error
	ImportFromCSV(breeds []Breed) error
}

type BreedRepository struct {
	db *sql.DB
}

func NewBreedRepository(db *sql.DB) *BreedRepository {
	return &BreedRepository{db: db}
}

func (r *BreedRepository) GetAll(species string, weightMin, weightMax *int, petSize string, limit, offset int) ([]Breed, error) {
	query := "SELECT id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight FROM breeds WHERE 1=1"
	args := []interface{}{}

	if species != "" {
		query += " AND species = ?"
		args = append(args, species)
	}

	if petSize != "" {
		query += " AND pet_size = ?"
		args = append(args, petSize)
	}

	if weightMin != nil {
		query += " AND (average_male_adult_weight >= ? OR average_female_adult_weight >= ?)"
		args = append(args, *weightMin, *weightMin)
	}

	if weightMax != nil {
		query += " AND (average_male_adult_weight <= ? OR average_female_adult_weight <= ?)"
		args = append(args, *weightMax, *weightMax)
	}

	query += " ORDER BY name"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)

		if offset > 0 {
			query += " OFFSET ?"
			args = append(args, offset)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des races: %w", err)
	}
	defer rows.Close()

	var breeds []Breed
	for rows.Next() {
		var breed Breed
		err := rows.Scan(
			&breed.ID,
			&breed.Species,
			&breed.PetSize,
			&breed.Name,
			&breed.AverageMaleAdultWeight,
			&breed.AverageFemaleAdultWeight,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du scan de la race: %w", err)
		}
		breeds = append(breeds, breed)
	}

	return breeds, nil
}

func (r *BreedRepository) GetByID(id int) (*Breed, error) {
	query := "SELECT id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight FROM breeds WHERE id = ?"

	var breed Breed
	err := r.db.QueryRow(query, id).Scan(
		&breed.ID,
		&breed.Species,
		&breed.PetSize,
		&breed.Name,
		&breed.AverageMaleAdultWeight,
		&breed.AverageFemaleAdultWeight,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de la race: %w", err)
	}

	return &breed, nil
}

func (r *BreedRepository) Create(breed *Breed) (*Breed, error) {
	query := `INSERT INTO breeds (species, pet_size, name, average_male_adult_weight, average_female_adult_weight) 
			  VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, breed.Species, breed.PetSize, breed.Name, breed.AverageMaleAdultWeight, breed.AverageFemaleAdultWeight)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la race: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de l'ID: %w", err)
	}

	breed.ID = int(id)
	return breed, nil
}

func (r *BreedRepository) Update(id int, breed *Breed) (*Breed, error) {
	setParts := []string{}
	args := []interface{}{}

	if breed.Species != "" {
		setParts = append(setParts, "species = ?")
		args = append(args, breed.Species)
	}
	if breed.PetSize != "" {
		setParts = append(setParts, "pet_size = ?")
		args = append(args, breed.PetSize)
	}
	if breed.Name != "" {
		setParts = append(setParts, "name = ?")
		args = append(args, breed.Name)
	}
	if breed.AverageMaleAdultWeight > 0 {
		setParts = append(setParts, "average_male_adult_weight = ?")
		args = append(args, breed.AverageMaleAdultWeight)
	}
	if breed.AverageFemaleAdultWeight > 0 {
		setParts = append(setParts, "average_female_adult_weight = ?")
		args = append(args, breed.AverageFemaleAdultWeight)
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("aucun champ à mettre à jour")
	}

	query := "UPDATE breeds SET " + strings.Join(setParts, ", ") + " WHERE id = ?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la mise à jour de la race: %w", err)
	}

	return r.GetByID(id)
}

func (r *BreedRepository) Delete(id int) error {
	query := "DELETE FROM breeds WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression de la race: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de la suppression: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("race non trouvée")
	}

	return nil
}

func (r *BreedRepository) ImportFromCSV(breeds []Breed) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("erreur lors du début de la transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO breeds (species, pet_size, name, average_male_adult_weight, average_female_adult_weight) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE species=VALUES(species), pet_size=VALUES(pet_size), average_male_adult_weight=VALUES(average_male_adult_weight), average_female_adult_weight=VALUES(average_female_adult_weight)")
	if err != nil {
		return fmt.Errorf("erreur lors de la préparation de la requête: %w", err)
	}
	defer stmt.Close()

	for _, breed := range breeds {
		_, err := stmt.Exec(breed.Species, breed.PetSize, breed.Name, breed.AverageMaleAdultWeight, breed.AverageFemaleAdultWeight)
		if err != nil {
			return fmt.Errorf("erreur lors de l'insertion de la race %s: %w", breed.Name, err)
		}
	}

	return tx.Commit()
}
