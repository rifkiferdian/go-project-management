package services

import (
	"errors"
	"fmt"
	"gobase-app/models"
	"gobase-app/repositories"
	"strings"
)

type DivisionService struct {
	Repo *repositories.DivisionRepository
}

func (s *DivisionService) GetDivisions() ([]models.Division, error) {
	return s.Repo.GetAll()
}

func (s *DivisionService) CreateDivision(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama divisi wajib diisi")
	}

	exists, err := s.Repo.ExistsByName(name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("divisi %s sudah ada", name)
	}

	return s.Repo.Create(name)
}

func (s *DivisionService) UpdateDivision(id int, name string) error {
	if id <= 0 {
		return errors.New("divisi tidak valid")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama divisi wajib diisi")
	}

	exists, err := s.Repo.ExistsByNameExceptID(name, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("divisi %s sudah ada", name)
	}

	return s.Repo.Update(id, name)
}

func (s *DivisionService) DeleteDivision(id int) error {
	if id <= 0 {
		return errors.New("divisi tidak valid")
	}
	return s.Repo.Delete(id)
}
