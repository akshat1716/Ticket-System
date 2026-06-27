package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"ticket-system/middleware"
	"ticket-system/services"
	"ticket-system/utils"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	ticketService *services.TicketService
}

func NewTicketController(ticketService *services.TicketService) *TicketController {
	return &TicketController{ticketService: ticketService}
}

func (ctrl *TicketController) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input services.CreateTicketInput
	if !utils.BindJSON(c, &input) {
		return
	}

	ticket, err := ctrl.ticketService.Create(userID, input)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func (ctrl *TicketController) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tickets, err := ctrl.ticketService.ListByUser(userID)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "Failed to fetch tickets")
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func (ctrl *TicketController) GetByID(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ticketID, err := parseTicketID(c)
	if err != nil {
		return
	}

	ticket, err := ctrl.ticketService.GetByID(ticketID, userID)
	if err != nil {
		handleTicketError(c, err)
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func (ctrl *TicketController) UpdateStatus(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ticketID, err := parseTicketID(c)
	if err != nil {
		return
	}

	var input services.UpdateStatusInput
	if !utils.BindJSON(c, &input) {
		return
	}

	ticket, err := ctrl.ticketService.UpdateStatus(ticketID, userID, input)
	if err != nil {
		handleTicketError(c, err)
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func parseTicketID(c *gin.Context) (uint, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "Invalid ticket ID")
		return 0, err
	}
	return uint(id), nil
}

func handleTicketError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrTicketNotFound):
		utils.JSONError(c, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrForbidden):
		utils.JSONError(c, http.StatusForbidden, err.Error())
	default:
		utils.JSONError(c, http.StatusBadRequest, err.Error())
	}
}
