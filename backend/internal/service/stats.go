package service

import (
	"backend/internal/domain"
	"backend/internal/repo/sqlc"
	"backend/pkg"
	"context"
)

type StatsService struct {
	queries *sqlc.Queries
}

func NewStatsService(queries *sqlc.Queries) *StatsService {
	return &StatsService{queries: queries}
}

func (s *StatsService) GetAdminStats(ctx context.Context) (domain.AdminStatsResponse, error) {
	summary, err := s.queries.GetSummaryStats(ctx)
	if err != nil {
		return domain.AdminStatsResponse{}, err
	}

	revenueRows, err := s.queries.GetRevenueLast30Days(ctx)
	if err != nil {
		return domain.AdminStatsResponse{}, err
	}

	bookingRows, err := s.queries.GetBookingsLast30Days(ctx)
	if err != nil {
		return domain.AdminStatsResponse{}, err
	}

	zoneRows, err := s.queries.GetBookingsByZoneType(ctx)
	if err != nil {
		return domain.AdminStatsResponse{}, err
	}

	res := domain.AdminStatsResponse{
		Summary: domain.SummaryStats{
			TotalBookings:   summary.TotalBookings,
			TotalRevenue:    pkg.NumericToString(summary.TotalRevenue),
			ActiveUsers:     summary.ActiveUsers,
			PendingBookings: summary.PendingBookings,
		},
		Revenue:  make([]domain.TimeSeriesData, 0, len(revenueRows)),
		Bookings: make([]domain.TimeSeriesData, 0, len(bookingRows)),
		ByZone:   make([]domain.ZoneStats, 0, len(zoneRows)),
	}

	for _, r := range revenueRows {
		res.Revenue = append(res.Revenue, domain.TimeSeriesData{
			Date:  r.Date.Time,
			Value: pkg.NumericToFloat64(r.Revenue),
		})
	}

	for _, r := range bookingRows {
		res.Bookings = append(res.Bookings, domain.TimeSeriesData{
			Date:  r.Date.Time,
			Value: float64(r.Count),
		})
	}

	for _, r := range zoneRows {
		res.ByZone = append(res.ByZone, domain.ZoneStats{
			ZoneType: string(r.ZoneType),
			Count:    r.Count,
		})
	}

	return res, nil
}
