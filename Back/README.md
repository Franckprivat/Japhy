# Japhy Backend Test

## Description
Ce projet est une API REST en Go permettant de gérer les races d'animaux (chiens et chats) pour un back-office de startup pet food. L'API permet de créer, lire, mettre à jour, supprimer et filtrer les races, stockées en base MySQL. Les données initiales peuvent être importées depuis un fichier CSV.

## Stack technique
- Go
- MySQL
- Docker & Docker Compose
- [Gorilla Mux](https://github.com/gorilla/mux) pour le routing
- [Charmbracelet Log](https://github.com/charmbracelet/log) pour les logs
- Tests unitaires avec [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)

## Prérequis
- [Docker](https://www.docker.com/products/docker-desktop/)
- [Git](https://git-scm.com/downloads)
- (optionnel) [Go](https://golang.org/dl/) pour développement local

## Installation & Lancement

1. **Cloner le projet**
   ```sh
   git clone git@github.com:Franckprivat/Backend---Japhy.git
   cd Backend
   ```
2. **Construire et lancer les services**
   ```sh
   docker compose build
   docker compose up -d
   ```
3. **Vérifier que l'API est en ligne**
   ```sh
   curl http://localhost:50010/health
   # ou via Postman
   ```

## Utilisation de l'API

### Endpoints principaux

- `GET    /breeds` : Liste toutes les races (filtres possibles)
- `POST   /breeds` : Crée une nouvelle race
- `GET    /breeds/{id}` : Détail d'une race
- `PUT    /breeds/{id}` : Met à jour une race
- `DELETE /breeds/{id}` : Supprime une race
- `POST   /import-breeds` : Importe les races depuis le CSV
- `GET    /health` : Vérifie que l'API tourne

### Exemple de requête POST (création d'une race)
```json
POST http://localhost:50010/breeds
Content-Type: application/json
{
  "species": "dog",
  "pet_size": "medium",
  "name": "Border Collie",
  "average_male_adult_weight": 20,
  "average_female_adult_weight": 18
}
```

### Exemple de filtre
```sh
GET http://localhost:50010/breeds?species=dog&weight_min=10&weight_max=30
```

## Importer les races depuis le CSV
- Placez votre fichier `breeds.csv` à la racine du projet.
- Appelez l'endpoint :
  ```sh
  curl -X POST http://localhost:50010/import-breeds
  ```

## Lancer les tests unitaires
```sh
go test ./...
```
- Les tests unitaires utilisent des mocks et ne nécessitent pas de base MySQL réelle.

## Structure du projet
```
Backend/
  ├── internal/
  │   ├── handlers/         # Handlers HTTP
  │   ├── repository/       # Accès base de données
  │   └── service/          # Services (CSV, etc.)
  ├── database_actions/     # Migrations SQL
  ├── breeds.csv            # Données de races (CSV)
  ├── main.go               # Point d'entrée
  └── Dockerfile, docker-compose.yml
```

Contact
Pour toute question ou suggestion :

Email : franck.kiemde@epitech.eu
---
