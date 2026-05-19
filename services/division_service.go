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

func (s *DivisionService) CreateDivision(name, prefix string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama divisi wajib diisi")
	}
	prefix = normalizeDivisionPrefix(prefix)
	if err := validateDivisionPrefix(prefix); err != nil {
		return err
	}

	exists, err := s.Repo.ExistsByName(name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("divisi %s sudah ada", name)
	}

	if prefix != "" {
		existsPrefix, err := s.Repo.ExistsByPrefix(prefix)
		if err != nil {
			return err
		}
		if existsPrefix {
			return fmt.Errorf("prefix divisi %s sudah dipakai", prefix)
		}
	}

	return s.Repo.Create(name, prefix)
}

func (s *DivisionService) UpdateDivision(id int, name, prefix string) error {
	if id <= 0 {
		return errors.New("divisi tidak valid")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama divisi wajib diisi")
	}
	prefix = normalizeDivisionPrefix(prefix)
	if err := validateDivisionPrefix(prefix); err != nil {
		return err
	}

	exists, err := s.Repo.ExistsByNameExceptID(name, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("divisi %s sudah ada", name)
	}

	if prefix != "" {
		existsPrefix, err := s.Repo.ExistsByPrefixExceptID(prefix, id)
		if err != nil {
			return err
		}
		if existsPrefix {
			return fmt.Errorf("prefix divisi %s sudah dipakai", prefix)
		}
	}

	return s.Repo.Update(id, name, prefix)
}

func (s *DivisionService) DeleteDivision(id int) error {
	if id <= 0 {
		return errors.New("divisi tidak valid")
	}
	return s.Repo.Delete(id)
}

func normalizeDivisionPrefix(prefix string) string {
	return strings.ToUpper(strings.TrimSpace(prefix))
}

func validateDivisionPrefix(prefix string) error {
	if prefix == "" {
		return nil
	}
	if len(prefix) > 10 {
		return errors.New("prefix divisi maksimal 10 karakter")
	}
	for _, ch := range prefix {
		if (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') {
			return errors.New("prefix divisi hanya boleh huruf A-Z dan angka 0-9 tanpa spasi")
		}
	}
	return nil
}
