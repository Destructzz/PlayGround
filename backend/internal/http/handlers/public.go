package handlers

import (
	"backend/internal/http/response"
	"backend/internal/service"
	"backend/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Public struct {
	publicService *service.PublicService
}

func NewPublic(publicService *service.PublicService) *Public {
	return &Public{publicService: publicService}
}

// Home returns summary data for the homepage cards.
// @Summary     Home catalog
// @Description Returns homepage summaries for gaming/lounge/event
// @Tags        public
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Router      /api/v1/public/home [get]
func (p *Public) Home(c *gin.Context) {
	data, err := p.publicService.HomeCatalog(c.Request.Context())
	if err != nil {
		zap.L().Warn("home catalog error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("catalog_error", "Failed to load home catalog", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("gaming", data.Gaming),
		response.WithData("lounge", data.Lounge),
		response.WithData("event", data.Event),
	).JSON(c)
}

// Lounge returns the lounge catalog for the /lounge page.
// @Summary     Lounge catalog
// @Description Returns lounge zones with services and availability
// @Tags        public
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Router      /api/v1/public/lounge [get]
func (p *Public) Lounge(c *gin.Context) {
	items, err := p.publicService.LoungeCatalog(c.Request.Context())
	if err != nil {
		zap.L().Warn("lounge catalog error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("catalog_error", "Failed to load lounge catalog", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("zones", items),
	).JSON(c)
}

// Event returns the event catalog for the /event page.
// @Summary     Event catalog
// @Description Returns event zones with services and availability
// @Tags        public
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Router      /api/v1/public/event [get]
func (p *Public) Event(c *gin.Context) {
	items, err := p.publicService.EventCatalog(c.Request.Context())
	if err != nil {
		zap.L().Warn("event catalog error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("catalog_error", "Failed to load event catalog", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("zones", items),
	).JSON(c)
}

// Gaming returns the gaming catalog for the /gaming page.
// @Summary     Gaming catalog
// @Description Returns gaming zones with places, configurations, and availability
// @Tags        public
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Router      /api/v1/public/gaming [get]
func (p *Public) Gaming(c *gin.Context) {
	data, err := p.publicService.GamingCatalog(c.Request.Context())
	if err != nil {
		zap.L().Warn("gaming catalog error", zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("catalog_error", "Failed to load gaming catalog", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("zone_tags", data.ZoneTags),
		response.WithData("zones", data.Zones),
		response.WithData("configurations", data.Configurations),
	).JSON(c)
}

// GamingAvailability returns taken gaming intervals for a zone and date.
// @Summary     Gaming availability
// @Description Returns booked place intervals for a gaming zone on the selected date
// @Tags        public
// @Produce     json
// @Param       zone_id query int true "Zone ID"
// @Param       date query string true "Date in YYYY-MM-DD"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/public/gaming/availability [get]
func (p *Public) GamingAvailability(c *gin.Context) {
	rawZoneID := c.Query("zone_id")
	if rawZoneID == "" {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("invalid_zone_id", "zone_id query parameter is required", nil),
		).JSON(c)
		return
	}

	zoneID, err := pkg.ParsePositiveInt64(rawZoneID)
	if err != nil || zoneID <= 0 {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("invalid_zone_id", "zone_id must be a positive integer", nil),
		).JSON(c)
		return
	}

	rawDate := c.Query("date")
	if rawDate == "" {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("invalid_date", "date query parameter is required", nil),
		).JSON(c)
		return
	}

	selectedDate, err := pkg.ParseDateYYYYMMDD(rawDate)
	if err != nil {
		response.NewResponseBuilder(
			response.WithStatus(http.StatusBadRequest),
			response.WithError("invalid_date", "date must be in YYYY-MM-DD format", nil),
		).JSON(c)
		return
	}

	data, err := p.publicService.GamingAvailability(c.Request.Context(), zoneID, selectedDate)
	if err != nil {
		zap.L().Warn("gaming availability error", zap.Int64("zone_id", zoneID), zap.Error(err))
		response.NewResponseBuilder(
			response.WithStatus(http.StatusInternalServerError),
			response.WithError("availability_error", "Failed to load gaming availability", nil),
		).JSON(c)
		return
	}

	response.NewResponseBuilder(
		response.WithData("zone_id", data.ZoneID),
		response.WithData("date", data.Date),
		response.WithData("bookings", data.Bookings),
	).JSON(c)
}
