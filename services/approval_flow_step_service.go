package services

import (
	"errors"
	"fmt"
	"gobase-app/models"
	"gobase-app/repositories"
	"strings"
)

type ApprovalFlowStepService struct {
	Repo *repositories.ApprovalFlowStepRepository
}

func (s *ApprovalFlowStepService) GetApprovalFlowSteps() ([]models.ApprovalFlowStep, error) {
	return s.Repo.GetAll()
}

func (s *ApprovalFlowStepService) GetFlowOptions() ([]models.ApprovalFlowOption, error) {
	return s.Repo.GetFlowOptions()
}

func (s *ApprovalFlowStepService) CreateApprovalFlowStep(flowID, stepOrder int, stepName, approvalRule string, isActive bool) error {
	stepName, approvalRule, err := normalizeStepInput(stepName, approvalRule)
	if err != nil {
		return err
	}
	if err := validateStepKeyInput(s.Repo, flowID, stepOrder, 0); err != nil {
		return err
	}

	return s.Repo.Create(flowID, stepOrder, stepName, approvalRule, isActive)
}

func (s *ApprovalFlowStepService) UpdateApprovalFlowStep(id, flowID, stepOrder int, stepName, approvalRule string, isActive bool) error {
	if id <= 0 {
		return errors.New("approval flow step tidak valid")
	}

	stepName, approvalRule, err := normalizeStepInput(stepName, approvalRule)
	if err != nil {
		return err
	}
	if err := validateStepKeyInput(s.Repo, flowID, stepOrder, id); err != nil {
		return err
	}

	return s.Repo.Update(id, flowID, stepOrder, stepName, approvalRule, isActive)
}

func (s *ApprovalFlowStepService) DeleteApprovalFlowStep(id int) error {
	if id <= 0 {
		return errors.New("approval flow step tidak valid")
	}
	return s.Repo.Delete(id)
}

func normalizeStepInput(stepName, approvalRule string) (string, string, error) {
	stepName = strings.TrimSpace(stepName)
	approvalRule = strings.ToLower(strings.TrimSpace(approvalRule))

	if stepName == "" {
		return "", "", errors.New("nama step wajib diisi")
	}
	if len(stepName) > 255 {
		return "", "", errors.New("nama step maksimal 255 karakter")
	}
	if approvalRule != "any" && approvalRule != "all" {
		return "", "", errors.New("approval rule harus any atau all")
	}

	return stepName, approvalRule, nil
}

func validateStepKeyInput(repo *repositories.ApprovalFlowStepRepository, flowID, stepOrder, excludeID int) error {
	if flowID <= 0 {
		return errors.New("flow wajib dipilih")
	}
	if stepOrder <= 0 {
		return errors.New("step order wajib lebih dari 0")
	}

	flowExists, err := repo.ExistsFlowByID(flowID)
	if err != nil {
		return err
	}
	if !flowExists {
		return fmt.Errorf("flow dengan id %d tidak ditemukan", flowID)
	}

	var exists bool
	if excludeID > 0 {
		exists, err = repo.ExistsStepOrderExceptID(flowID, stepOrder, excludeID)
	} else {
		exists, err = repo.ExistsStepOrder(flowID, stepOrder)
	}
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("step order %d sudah dipakai pada flow ini", stepOrder)
	}

	return nil
}
