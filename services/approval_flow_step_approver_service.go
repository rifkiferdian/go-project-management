package services

import (
	"errors"
	"fmt"
	"gobase-app/models"
	"gobase-app/repositories"
	"strings"
)

type ApprovalFlowStepApproverService struct {
	Repo *repositories.ApprovalFlowStepApproverRepository
}

func (s *ApprovalFlowStepApproverService) GetApprovers() ([]models.ApprovalFlowStepApprover, error) {
	return s.Repo.GetAll()
}

func (s *ApprovalFlowStepApproverService) GetStepOptions() ([]models.ApprovalFlowStepOption, error) {
	return s.Repo.GetStepOptions()
}

func (s *ApprovalFlowStepApproverService) GetUserOptions() ([]models.LookupOption, error) {
	return s.Repo.GetUserOptions()
}

func (s *ApprovalFlowStepApproverService) GetRoleOptions() ([]models.LookupOption, error) {
	return s.Repo.GetRoleOptions()
}

func (s *ApprovalFlowStepApproverService) GetDivisionOptions() ([]models.LookupOption, error) {
	return s.Repo.GetDivisionOptions()
}

func (s *ApprovalFlowStepApproverService) CreateApprover(stepID int, approverType string, userID, roleID, divisionID int, isActive bool) error {
	approverType = strings.ToLower(strings.TrimSpace(approverType))
	userID, roleID, divisionID, err := s.validateApproverInput(stepID, approverType, userID, roleID, divisionID, 0)
	if err != nil {
		return err
	}

	return s.Repo.Create(stepID, approverType, userID, roleID, divisionID, isActive)
}

func (s *ApprovalFlowStepApproverService) UpdateApprover(id, stepID int, approverType string, userID, roleID, divisionID int, isActive bool) error {
	if id <= 0 {
		return errors.New("approval flow step approver tidak valid")
	}

	approverType = strings.ToLower(strings.TrimSpace(approverType))
	userID, roleID, divisionID, err := s.validateApproverInput(stepID, approverType, userID, roleID, divisionID, id)
	if err != nil {
		return err
	}

	return s.Repo.Update(id, stepID, approverType, userID, roleID, divisionID, isActive)
}

func (s *ApprovalFlowStepApproverService) DeleteApprover(id int) error {
	if id <= 0 {
		return errors.New("approval flow step approver tidak valid")
	}
	return s.Repo.Delete(id)
}

func (s *ApprovalFlowStepApproverService) validateApproverInput(stepID int, approverType string, userID, roleID, divisionID int, excludeID int) (int, int, int, error) {
	if stepID <= 0 {
		return 0, 0, 0, errors.New("step wajib dipilih")
	}

	stepExists, err := s.Repo.ExistsStepByID(stepID)
	if err != nil {
		return 0, 0, 0, err
	}
	if !stepExists {
		return 0, 0, 0, fmt.Errorf("step dengan id %d tidak ditemukan", stepID)
	}

	switch approverType {
	case "user":
		if userID <= 0 {
			return 0, 0, 0, errors.New("user approver wajib dipilih")
		}
		ok, err := s.Repo.ExistsUserByID(userID)
		if err != nil {
			return 0, 0, 0, err
		}
		if !ok {
			return 0, 0, 0, fmt.Errorf("user approver id %d tidak ditemukan", userID)
		}
		roleID = 0
		divisionID = 0
	case "role":
		if roleID <= 0 {
			return 0, 0, 0, errors.New("role approver wajib dipilih")
		}
		ok, err := s.Repo.ExistsRoleByID(roleID)
		if err != nil {
			return 0, 0, 0, err
		}
		if !ok {
			return 0, 0, 0, fmt.Errorf("role approver id %d tidak ditemukan", roleID)
		}
		userID = 0
		divisionID = 0
	case "division":
		if divisionID <= 0 {
			return 0, 0, 0, errors.New("division approver wajib dipilih")
		}
		ok, err := s.Repo.ExistsDivisionByID(divisionID)
		if err != nil {
			return 0, 0, 0, err
		}
		if !ok {
			return 0, 0, 0, fmt.Errorf("division approver id %d tidak ditemukan", divisionID)
		}
		userID = 0
		roleID = 0
	default:
		return 0, 0, 0, errors.New("tipe approver tidak valid")
	}

	dupExists, err := s.Repo.ExistsDuplicate(stepID, approverType, userID, roleID, divisionID, excludeID)
	if err != nil {
		return 0, 0, 0, err
	}
	if dupExists {
		return 0, 0, 0, errors.New("approver sudah terdaftar pada step ini")
	}

	return userID, roleID, divisionID, nil
}
