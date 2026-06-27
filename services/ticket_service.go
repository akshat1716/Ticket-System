package services

import (
	"errors"

	"ticket-system/models"
	"ticket-system/utils"

	"gorm.io/gorm"
)

var (
	ErrTicketNotFound = errors.New("Ticket not found")
	ErrForbidden      = errors.New("Forbidden")
)

type TicketService struct {
	db *gorm.DB
}

func NewTicketService(db *gorm.DB) *TicketService {
	return &TicketService{db: db}
}

type CreateTicketInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateStatusInput struct {
	Status string `json:"status"`
}

func (s *TicketService) Create(userID uint, input CreateTicketInput) (*models.Ticket, error) {
	if err := utils.ValidateTitle(input.Title); err != nil {
		return nil, err
	}

	ticket := models.Ticket{
		Title:       input.Title,
		Description: input.Description,
		Status:      models.StatusOpen,
		UserID:      userID,
	}

	if err := s.db.Create(&ticket).Error; err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (s *TicketService) ListByUser(userID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	if err := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&tickets).Error; err != nil {
		return nil, err
	}
	if tickets == nil {
		tickets = []models.Ticket{}
	}
	return tickets, nil
}

func (s *TicketService) GetByID(ticketID, userID uint) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	if ticket.UserID != userID {
		return nil, ErrForbidden
	}

	return &ticket, nil
}

func (s *TicketService) UpdateStatus(ticketID, userID uint, input UpdateStatusInput) (*models.Ticket, error) {
	if err := utils.ValidateStatus(input.Status); err != nil {
		return nil, err
	}

	var ticket models.Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	if ticket.UserID != userID {
		return nil, ErrForbidden
	}

	if err := validateStatusTransition(ticket.Status, input.Status); err != nil {
		return nil, err
	}

	ticket.Status = input.Status
	if err := s.db.Save(&ticket).Error; err != nil {
		return nil, err
	}

	return &ticket, nil
}

func validateStatusTransition(current, next string) error {
	if current == next {
		return errors.New("Status is already " + current)
	}

	switch current {
	case models.StatusOpen:
		if next != models.StatusInProgress {
			return errors.New("Invalid status transition: open can only transition to in_progress")
		}
	case models.StatusInProgress:
		if next != models.StatusClosed {
			return errors.New("Invalid status transition: in_progress can only transition to closed")
		}
	case models.StatusClosed:
		return errors.New("Invalid status transition: closed tickets cannot be updated")
	default:
		return errors.New("Invalid current status")
	}

	return nil
}
