package models

type Dashboard struct {
	CustomerCount       uint `json:"customer_count"`
	StationCount        uint `json:"station_count"`
	ProductCount        uint `json:"product_count"`
	StationProductCount uint `json:"station_product_count"`
}

type DashboardTasks struct {
	MonthCount uint `json:"month_count"`
	WeekCount  uint `json:"week_count"`
	DayCount   uint `json:"day_count"`
}
