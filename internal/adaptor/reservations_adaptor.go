package adaptor

import (
	"aplikasi-pos-team-boolean/internal/dto"
	"aplikasi-pos-team-boolean/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ReservationsAdaptor menangani semua request HTTP untuk reservations
type ReservationsAdaptor struct {
	reservationsUsecase usecase.ReservationsUseCase
	logger              *zap.Logger
}

// NewReservationsAdaptor membuat instance baru dari ReservationsAdaptor
func NewReservationsAdaptor(reservationsUsecase usecase.ReservationsUseCase, logger *zap.Logger) *ReservationsAdaptor {
	return &ReservationsAdaptor{
		reservationsUsecase: reservationsUsecase,
		logger:              logger,
	}
}

// GetAllReservations menangani request GET /reservations
func (h *ReservationsAdaptor) GetAllReservations(c *gin.Context) {
	h.logger.Debug("GetAllReservations handler called", zap.String("client_ip", c.ClientIP()))

	// Call usecase
	reservations, err := h.reservationsUsecase.GetAllReservations(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get reservations",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to get reservations: " + err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Reservations retrieved successfully",
		zap.Int("count", len(reservations)),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Reservations retrieved successfully",
		"data":    reservations,
	})
}

// CreateReservation menangani request POST /reservations
func (h *ReservationsAdaptor) CreateReservation(c *gin.Context) {
	h.logger.Debug("CreateReservation handler called", zap.String("client_ip", c.ClientIP()))

	var req dto.ReservationCreateRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for create reservation",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Call usecase
	reservation, err := h.reservationsUsecase.CreateReservation(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create reservation",
			zap.String("customer_name", req.CustomerName),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Reservation created successfully",
		zap.Int64("id", reservation.ID),
		zap.String("customer_name", req.CustomerName),
	)

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Reservation created successfully",
		"data":    reservation,
	})
}

// UpdateReservation menangani request PUT /reservations/:id
func (h *ReservationsAdaptor) UpdateReservation(c *gin.Context) {
	h.logger.Debug("UpdateReservation handler called", zap.String("client_ip", c.ClientIP()))

	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid reservation ID",
			zap.String("id", idStr),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid reservation ID",
			"data":    nil,
		})
		return
	}

	var req dto.ReservationUpdateRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for update reservation",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Call usecase
	if err := h.reservationsUsecase.UpdateReservation(c.Request.Context(), uint(id), req); err != nil {
		h.logger.Error("Failed to update reservation",
			zap.Uint64("id", id),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Reservation updated successfully",
		zap.Uint64("id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Reservation updated successfully",
		"data":    nil,
	})
}

// DeleteReservation menangani request DELETE /reservations/:id
func (h *ReservationsAdaptor) DeleteReservation(c *gin.Context) {
	h.logger.Debug("DeleteReservation handler called", zap.String("client_ip", c.ClientIP()))

	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid reservation ID",
			zap.String("id", idStr),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid reservation ID",
			"data":    nil,
		})
		return
	}

	// Call usecase
	if err := h.reservationsUsecase.DeleteReservation(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete reservation",
			zap.Uint64("id", id),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Reservation deleted successfully",
		zap.Uint64("id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Reservation deleted successfully",
		"data":    nil,
	})
}
