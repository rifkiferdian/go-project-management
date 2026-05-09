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

func (s *ProjectService) GetDivisionOptions() ([]models.DivisionOption, error) {
	return s.Repo.GetDivisionOptions()
}

func (s *ProjectService) GetPriorityOptions() ([]models.ProjectPriorityOption, error) {
	return s.Repo.GetPriorityOptions()
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
	divisionIDs := uniqueInt64s(input.DivisionIDs)

	if name == "" {
		return input, errors.New("nama project wajib diisi")
	}
	if input.OwnerID <= 0 {
		return input, errors.New("owner project wajib dipilih")
	}
	if input.DeveloperID <= 0 {
		return input, errors.New("developer project wajib dipilih")
	}
	if len(divisionIDs) == 0 {
		return input, errors.New("minimal pilih 1 divisi requester")
	}
	if input.StatusID <= 0 {
		return input, errors.New("status project wajib dipilih")
	}
	if input.PriorityID <= 0 {
		return input, errors.New("prioritas project wajib dipilih")
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

	existingDivisionMap, err := s.Repo.FindExistingDivisionIDs(divisionIDs)
	if err != nil {
		return input, err
	}
	var missingDivisions []string
	for _, divisionID := range divisionIDs {
		if !existingDivisionMap[divisionID] {
			missingDivisions = append(missingDivisions, fmt.Sprintf("%d", divisionID))
		}
	}
	if len(missingDivisions) > 0 {
		return input, fmt.Errorf("divisi requester tidak ditemukan: %s", strings.Join(missingDivisions, ", "))
	}

	isITDeveloper, err := s.Repo.IsUserInDivision(input.DeveloperID, "IT")
	if err != nil {
		return input, err
	}
	if !isITDeveloper {
		return input, errors.New("developer harus user dari divisi IT")
	}

	input.Name = name
	input.Description = description
	input.DivisionIDs = divisionIDs
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
	divisionIDs := uniqueInt64s(input.DivisionIDs)

	if name == "" {
		return input, errors.New("nama project wajib diisi")
	}
	if input.OwnerID <= 0 {
		return input, errors.New("owner project wajib dipilih")
	}
	if input.DeveloperID <= 0 {
		return input, errors.New("developer project wajib dipilih")
	}
	if len(divisionIDs) == 0 {
		return input, errors.New("minimal pilih 1 divisi requester")
	}
	if input.StatusID <= 0 {
		return input, errors.New("status project wajib dipilih")
	}
	if input.PriorityID <= 0 {
		return input, errors.New("prioritas project wajib dipilih")
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

	existingDivisionMap, err := s.Repo.FindExistingDivisionIDs(divisionIDs)
	if err != nil {
		return input, err
	}
	var missingDivisions []string
	for _, divisionID := range divisionIDs {
		if !existingDivisionMap[divisionID] {
			missingDivisions = append(missingDivisions, fmt.Sprintf("%d", divisionID))
		}
	}
	if len(missingDivisions) > 0 {
		return input, fmt.Errorf("divisi requester tidak ditemukan: %s", strings.Join(missingDivisions, ", "))
	}

	isITDeveloper, err := s.Repo.IsUserInDivision(input.DeveloperID, "IT")
	if err != nil {
		return input, err
	}
	if !isITDeveloper {
		return input, errors.New("developer harus user dari divisi IT")
	}

	input.Name = name
	input.Description = description
	input.DivisionIDs = divisionIDs
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
