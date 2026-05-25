package services

import (
	"errors"
	"fmt"
	"gobase-app/models"
	"gobase-app/repositories"
	"strings"
)

type TicketTemplateService struct {
	Repo *repositories.TicketTemplateRepository
}

func (s *TicketTemplateService) GetSets() ([]models.TicketTemplateSet, error) {
	return s.Repo.GetSets()
}

func (s *TicketTemplateService) CreateSet(name, purpose, description string, isActive bool) error {
	name = strings.TrimSpace(name)
	purpose = normalizePurpose(purpose)
	description = strings.TrimSpace(description)

	if name == "" {
		return errors.New("nama set wajib diisi")
	}
	if err := validatePurpose(purpose); err != nil {
		return err
	}

	exists, err := s.Repo.ExistsSetByNamePurpose(name, purpose)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("set %s (%s) sudah ada", name, purpose)
	}

	return s.Repo.CreateSet(name, purpose, description, isActive)
}

func (s *TicketTemplateService) UpdateSet(id int, name, purpose, description string, isActive bool) error {
	if id <= 0 {
		return errors.New("set template tidak valid")
	}

	name = strings.TrimSpace(name)
	purpose = normalizePurpose(purpose)
	description = strings.TrimSpace(description)

	if name == "" {
		return errors.New("nama set wajib diisi")
	}
	if err := validatePurpose(purpose); err != nil {
		return err
	}

	exists, err := s.Repo.ExistsSetByNamePurposeExceptID(name, purpose, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("set %s (%s) sudah ada", name, purpose)
	}

	return s.Repo.UpdateSet(id, name, purpose, description, isActive)
}

func (s *TicketTemplateService) DeleteSet(id int) error {
	if id <= 0 {
		return errors.New("set template tidak valid")
	}
	return s.Repo.DeleteSet(id)
}

func (s *TicketTemplateService) GetEpics() ([]models.TicketTemplateEpic, error) {
	return s.Repo.GetEpics()
}

func (s *TicketTemplateService) CreateEpic(
	setID int,
	name, description string,
	startOffsetDays, dueOffsetDays *int,
	sortOrder int,
	isActive bool,
) error {
	if setID <= 0 {
		return errors.New("set template wajib dipilih")
	}
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		return errors.New("nama epic template wajib diisi")
	}

	setExists, err := s.Repo.SetExists(setID)
	if err != nil {
		return err
	}
	if !setExists {
		return errors.New("set template tidak ditemukan")
	}

	exists, err := s.Repo.ExistsEpicBySetAndName(setID, name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("epic %s pada set ini sudah ada", name)
	}

	return s.Repo.CreateEpic(setID, name, description, startOffsetDays, dueOffsetDays, normalizeOrder(sortOrder), isActive)
}

func (s *TicketTemplateService) UpdateEpic(
	id, setID int,
	name, description string,
	startOffsetDays, dueOffsetDays *int,
	sortOrder int,
	isActive bool,
) error {
	if id <= 0 {
		return errors.New("epic template tidak valid")
	}
	if setID <= 0 {
		return errors.New("set template wajib dipilih")
	}
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		return errors.New("nama epic template wajib diisi")
	}

	epicExists, err := s.Repo.EpicExists(id)
	if err != nil {
		return err
	}
	if !epicExists {
		return errors.New("epic template tidak ditemukan")
	}

	setExists, err := s.Repo.SetExists(setID)
	if err != nil {
		return err
	}
	if !setExists {
		return errors.New("set template tidak ditemukan")
	}

	exists, err := s.Repo.ExistsEpicBySetAndNameExceptID(setID, name, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("epic %s pada set ini sudah ada", name)
	}

	return s.Repo.UpdateEpic(id, setID, name, description, startOffsetDays, dueOffsetDays, normalizeOrder(sortOrder), isActive)
}

func (s *TicketTemplateService) DeleteEpic(id int) error {
	if id <= 0 {
		return errors.New("epic template tidak valid")
	}
	return s.Repo.DeleteEpic(id)
}

func (s *TicketTemplateService) GetItems() ([]models.TicketTemplateItem, error) {
	return s.Repo.GetItems()
}

func (s *TicketTemplateService) CreateItem(
	setID int,
	title, description string,
	templateEpicID, defaultTypeID, defaultPriorityID, defaultStatusID, defaultOwnerID, defaultResponsibleID *int,
	estimation *float64,
	startOffsetDays, dueOffsetDays *int,
	sortOrder int,
	isActive bool,
) error {
	if setID <= 0 {
		return errors.New("set template wajib dipilih")
	}
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	if title == "" {
		return errors.New("judul ticket template wajib diisi")
	}

	setExists, err := s.Repo.SetExists(setID)
	if err != nil {
		return err
	}
	if !setExists {
		return errors.New("set template tidak ditemukan")
	}

	if templateEpicID != nil && *templateEpicID > 0 {
		belongs, err := s.Repo.EpicBelongsToSet(*templateEpicID, setID)
		if err != nil {
			return err
		}
		if !belongs {
			return errors.New("epic template tidak sesuai dengan set")
		}
	}

	sortOrder = normalizeOrder(sortOrder)
	if estimation != nil && *estimation < 0 {
		return errors.New("estimasi tidak boleh negatif")
	}

	return s.Repo.CreateItem(
		setID,
		title,
		description,
		templateEpicID,
		defaultTypeID,
		defaultPriorityID,
		defaultStatusID,
		defaultOwnerID,
		defaultResponsibleID,
		estimation,
		startOffsetDays,
		dueOffsetDays,
		sortOrder,
		isActive,
	)
}

func (s *TicketTemplateService) UpdateItem(
	id, setID int,
	title, description string,
	templateEpicID, defaultTypeID, defaultPriorityID, defaultStatusID, defaultOwnerID, defaultResponsibleID *int,
	estimation *float64,
	startOffsetDays, dueOffsetDays *int,
	sortOrder int,
	isActive bool,
) error {
	if id <= 0 {
		return errors.New("item template tidak valid")
	}
	if setID <= 0 {
		return errors.New("set template wajib dipilih")
	}

	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	if title == "" {
		return errors.New("judul ticket template wajib diisi")
	}

	itemExists, err := s.Repo.ItemExists(id)
	if err != nil {
		return err
	}
	if !itemExists {
		return errors.New("item template tidak ditemukan")
	}

	setExists, err := s.Repo.SetExists(setID)
	if err != nil {
		return err
	}
	if !setExists {
		return errors.New("set template tidak ditemukan")
	}

	if templateEpicID != nil && *templateEpicID > 0 {
		belongs, err := s.Repo.EpicBelongsToSet(*templateEpicID, setID)
		if err != nil {
			return err
		}
		if !belongs {
			return errors.New("epic template tidak sesuai dengan set")
		}
	}

	sortOrder = normalizeOrder(sortOrder)
	if estimation != nil && *estimation < 0 {
		return errors.New("estimasi tidak boleh negatif")
	}

	return s.Repo.UpdateItem(
		id,
		setID,
		title,
		description,
		templateEpicID,
		defaultTypeID,
		defaultPriorityID,
		defaultStatusID,
		defaultOwnerID,
		defaultResponsibleID,
		estimation,
		startOffsetDays,
		dueOffsetDays,
		sortOrder,
		isActive,
	)
}

func (s *TicketTemplateService) DeleteItem(id int) error {
	if id <= 0 {
		return errors.New("item template tidak valid")
	}
	return s.Repo.DeleteItem(id)
}

func (s *TicketTemplateService) GetTicketTypeOptions() ([]models.TicketTemplateOption, error) {
	return s.Repo.GetTicketTypeOptions()
}

func (s *TicketTemplateService) GetTicketPriorityOptions() ([]models.TicketTemplateOption, error) {
	return s.Repo.GetTicketPriorityOptions()
}

func (s *TicketTemplateService) GetTicketStatusOptions() ([]models.TicketTemplateOption, error) {
	return s.Repo.GetTicketStatusOptions()
}

func (s *TicketTemplateService) GetUserOptions() ([]models.TicketTemplateOption, error) {
	return s.Repo.GetUserOptions()
}

func (s *TicketTemplateService) GetEpicOptions() ([]models.TicketTemplateOption, error) {
	return s.Repo.GetEpicOptions()
}

func validatePurpose(purpose string) error {
	switch purpose {
	case "new_project", "new_feature", "bugfix", "custom":
		return nil
	default:
		return errors.New("purpose tidak valid")
	}
}

func normalizePurpose(purpose string) string {
	return strings.ToLower(strings.TrimSpace(purpose))
}
