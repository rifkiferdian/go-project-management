package services

import (
	"database/sql"
	"errors"
	"fmt"
	"gobase-app/models"
	"gobase-app/repositories"
	"strings"
)

type ProjectService struct {
	Repo *repositories.ProjectRepository
}

func (s *ProjectService) GetProjects() ([]models.Project, error) {
	return s.Repo.GetAll()
}

func (s *ProjectService) GetProject(id int) (*models.Project, error) {
	if id <= 0 {
		return nil, errors.New("project id tidak valid")
	}

	project, err := s.Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("project dengan id %d tidak ditemukan", id)
		}
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) GetStatusOptions() ([]models.ProjectStatusOption, error) {
	return s.Repo.GetStatusOptions()
}

func (s *ProjectService) CreateProject(input models.ProjectCreateInput) error {
	params, err := s.validateCreateInput(input)
	if err != nil {
		return err
	}
	return s.Repo.Create(params)
}

func (s *ProjectService) UpdateProject(input models.ProjectUpdateInput) error {
	params, err := s.validateUpdateInput(input)
	if err != nil {
		return err
	}
	return s.Repo.Update(params)
}

func (s *ProjectService) DeleteProject(id int) error {
	if id <= 0 {
		return errors.New("project id tidak valid")
	}
	return s.Repo.Delete(id)
}

func (s *ProjectService) validateCreateInput(input models.ProjectCreateInput) (models.ProjectCreateInput, error) {
	name := strings.TrimSpace(input.Name)
	prefix := strings.ToUpper(strings.TrimSpace(input.TicketPrefix))
	description := strings.TrimSpace(input.Description)
	statusType := normalizeProjectStatusType(input.StatusType)
	projectType := normalizeProjectType(input.Type)

	if name == "" {
		return input, errors.New("nama project wajib diisi")
	}
	if input.OwnerID <= 0 {
		return input, errors.New("owner project wajib dipilih")
	}
	if input.StatusID <= 0 {
		return input, errors.New("status project wajib dipilih")
	}
	if prefix == "" || len(prefix) > 3 {
		return input, errors.New("ticket prefix wajib diisi maksimal 3 karakter")
	}

	exists, err := s.Repo.ExistsByTicketPrefix(prefix)
	if err != nil {
		return input, err
	}
	if exists {
		return input, fmt.Errorf("ticket prefix %s sudah digunakan", prefix)
	}

	input.Name = name
	input.Description = description
	input.TicketPrefix = prefix
	input.StatusType = statusType
	input.Type = projectType
	return input, nil
}

func (s *ProjectService) validateUpdateInput(input models.ProjectUpdateInput) (models.ProjectUpdateInput, error) {
	if input.ID <= 0 {
		return input, errors.New("project tidak valid")
	}

	name := strings.TrimSpace(input.Name)
	prefix := strings.ToUpper(strings.TrimSpace(input.TicketPrefix))
	description := strings.TrimSpace(input.Description)
	statusType := normalizeProjectStatusType(input.StatusType)
	projectType := normalizeProjectType(input.Type)

	if name == "" {
		return input, errors.New("nama project wajib diisi")
	}
	if input.OwnerID <= 0 {
		return input, errors.New("owner project wajib dipilih")
	}
	if input.StatusID <= 0 {
		return input, errors.New("status project wajib dipilih")
	}
	if prefix == "" || len(prefix) > 3 {
		return input, errors.New("ticket prefix wajib diisi maksimal 3 karakter")
	}

	exists, err := s.Repo.ExistsByTicketPrefixExceptID(prefix, input.ID)
	if err != nil {
		return input, err
	}
	if exists {
		return input, fmt.Errorf("ticket prefix %s sudah digunakan", prefix)
	}

	input.Name = name
	input.Description = description
	input.TicketPrefix = prefix
	input.StatusType = statusType
	input.Type = projectType
	return input, nil
}

func normalizeProjectType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "scrum":
		return "scrum"
	default:
		return "kanban"
	}
}

func normalizeProjectStatusType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "custom":
		return "custom"
	default:
		return "default"
	}
}
