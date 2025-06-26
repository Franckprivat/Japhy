CREATE TABLE IF NOT EXISTS breeds (
    id INT AUTO_INCREMENT PRIMARY KEY,
    species VARCHAR(50) NOT NULL,
    pet_size VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL UNIQUE,
    average_male_adult_weight INT NOT NULL,
    average_female_adult_weight INT NOT NULL
);