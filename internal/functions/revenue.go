package functions

import (
	"log"
	"tempfunctiontools/internal/database"
)

func GetRevenue(month int, year int, db *database.DbConfig) (float64, error) {
	log.Printf("month: %d, year: %d", month, year)

	// get the revenue
	rev, err := db.GetRevenueByMonthYear(month, year)
	if err != nil {
		return 0, err
	}

	return rev.Amount, nil
}
