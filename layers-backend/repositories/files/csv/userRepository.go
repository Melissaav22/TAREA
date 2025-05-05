package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"layersapi/entities"
	"os"
	"os/user"
	"strings"
	"time"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u UserRepository) GetAll() ([]entities.User, error) {
	file, err := os.Open(`data/data.csv`)
	if err != nil {
		return []entities.User{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return []entities.User{}, err
	}

	var result []entities.User

	for i, record := range records {
		if i == 0 {
			continue
		}

		createdAt, _ := time.Parse(time.RFC3339, record[3])
		updatedAt, _ := time.Parse(time.RFC3339, record[4])
		meta := entities.Metadata{
			CreatedAt: createdAt.String(),
			UpdatedAt: updatedAt.String(),
			CreatedBy: record[5],
			UpdatedBy: record[6],
		}
		result = append(result, entities.NewUser(record[0], record[1], record[2], meta))
	}

	return result, nil
}

func (u UserRepository) GetById(id string) (entities.User, error) {
	file, err := os.Open("data/data.csv")
	if err != nil {
		return entities.User{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return entities.User{}, err
	}

	for i, record := range records {
		if i == 0 {
			continue
		} else if strings.TrimSpace(record[0]) == strings.TrimSpace(id) {
			createdAt, _ := time.Parse(time.RFC3339, record[3])
			updatedAt, _ := time.Parse(time.RFC3339, record[4])
			createdBy := strings.Trim(record[5], `"`)
			updatedBy := strings.Trim(record[6], `"`)

			meta := entities.Metadata{
				CreatedAt: createdAt.String(),
				UpdatedAt: updatedAt.String(),
				CreatedBy: createdBy,
				UpdatedBy: updatedBy,
			}
			return entities.NewUser(record[0], record[1], record[2], meta), nil
		}

	}

	return entities.User{}, errors.New("user not found")
}

func (u UserRepository) Create(user entities.User) error {
	file, err := os.OpenFile("data/data.csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	newUser := []string{
		user.Id,
		user.Name,
		user.Email,
		user.Metadata.CreatedAt,
		user.Metadata.UpdatedAt,
		"webapp",
		"webapp",
	}

	if err := writer.Write(newUser); err != nil {
		return err
	}

	return nil
}

func (u UserRepository) Update(id, name, email string) error {

	file, err := os.Open("data/data.csv")
	if err != nil {
		fmt.Println("Opening in file has gone wrong", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Reading in file has gone wrong", err)
	}

	updated := false
	for i, record := range records {
		if i == 0 {
			continue
		}
		if record[0] == id {
			record[1] = name
			record[2] = email
			record[4] = time.Now().Format(time.RFC3339)
			current, err := user.Current()
			if err != nil {
				return err
			}
			record[6] = current.Name
			updated = true
			break
		}
	}

	if !updated {
		return errors.New("user to update not found")
	}

	temporary, err := os.Create("data/temp.csv")
	if err != nil {
		fmt.Println("creating temporary has gone wrong", err)
	}
	defer temporary.Close()

	writer := csv.NewWriter(temporary)
	err = writer.WriteAll(records)
	if err != nil {
		fmt.Println("writing in temporary has gone wrong", err)

	}

	file.Close()
	temporary.Close()

	err = os.Rename("data/temp.csv", "data/data.csv")
	if err != nil {
		fmt.Println("renaming has gone wrong", err)
	}

	return nil
}
