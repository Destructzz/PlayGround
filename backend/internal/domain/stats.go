package domain

import "time"

type SummaryStats struct {
	TotalBookings   int64  `json:"total_bookings"`
	TotalRevenue    string `json:"total_revenue"`
	ActiveUsers     int64  `json:"active_users"`
	PendingBookings int64  `json:"pending_bookings"`
}

type TimeSeriesData struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

type ZoneStats struct {
	ZoneType string `json:"zone_type"`
	Count    int64  `json:"count"`
}

type AdminStatsResponse struct {
	Summary  SummaryStats     `json:"summary"`
	Revenue  []TimeSeriesData `json:"revenue"`
	Bookings []TimeSeriesData `json:"bookings"`
	ByZone   []ZoneStats      `json:"by_zone"`
}
