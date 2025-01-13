package v1

import (
	"go-clean-arch/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents user http service
type apiServer struct {
	useCase usecase.HealthUsecase
}

type APIServer interface {
	ConfigRoute(eg *echo.Group)
}

func NewAPIServer(useCase usecase.HealthUsecase) APIServer {
	return &apiServer{
		useCase: useCase,
	}
}

// ConfigRoute creates new health http service
func (a *apiServer) ConfigRoute(eg *echo.Group) {
	route := eg.Group("/health")

	// Liveness
	// @Summary Fetch a list of all books.
	// @Description Fetch a list of all books.
	// @Tags Book
	// @Accept */*
	// @Security Bearer Authentication
	// @Param Authorization header string true "'Bearer _YOUR_TOKEN_'"
	// @Produce json
	// @Success 200
	// @Failure 500 {object} swagger.Error
	// @Router /api/v1/book [get]
	route.GET("/liveness", a.liveness)

	// Liveness
	// @Summary Fetch a list of all books.
	// @Description Fetch a list of all books.
	// @Tags Book
	// @Accept */*
	// @Security Bearer Authentication
	// @Param Authorization header string true "'Bearer _YOUR_TOKEN_'"
	// @Produce json
	// @Success 200
	// @Failure 500 {object} swagger.Error
	// @Router /api/v1/book [get]
	route.GET("/readiness", a.readiness)
}

func (h *apiServer) readiness(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *apiServer) liveness(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
