package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

const (
	databaseName = "revenue.db"
)

type DbConfig struct {
	db  *bun.DB
	ctx context.Context
}

type Revenue struct {
	bun.BaseModel `bun:"table:revenue"`
	ID            int     `bun:"id,pk,autoincrement"`
	Month         int     `bun:"month,notnull"`
	Year          int     `bun:"year,notnull"`
	Amount        float64 `bun:"amount"`
}

func (c *DbConfig) InitDb() error {
	sqldb, err := sql.Open(sqliteshim.ShimName, databaseName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	// defer sqldb.Close()

	c.db = bun.NewDB(sqldb, sqlitedialect.New())
	c.ctx = context.Background()

	if err := c.createRevenueTable(); err != nil {
		return fmt.Errorf("failed to create revenue table: %w", err)
	}

	revenues := generateRevenueData()
	if err := c.upsertRevenues(revenues); err != nil {
		log.Printf("failed to upsert revenues: %v", err)
		return fmt.Errorf("failed to upsert revenues: %w", err)
	}

	log.Println("Database initialized successfully")
	log.Printf("added revenues: %d", len(revenues))

	return nil
}

func (c *DbConfig) createRevenueTable() error {
	// Create the table
	_, err := c.db.NewCreateTable().
		Model((*Revenue)(nil)).
		IfNotExists().
		Exec(c.ctx)
	if err != nil {
		return fmt.Errorf("failed to create revenue table: %w", err)
	}

	// Create unique index for month and year combination
	_, err = c.db.NewCreateIndex().
		Model((*Revenue)(nil)).
		Index("idx_month_year").
		Column("month", "year").
		Unique().
		IfNotExists().
		Exec(c.ctx)
	if err != nil {
		return fmt.Errorf("failed to create unique index on month and year: %w", err)
	}

	return nil
}

func (c *DbConfig) upsertRevenues(revenues []Revenue) error {
	_, err := c.db.NewInsert().Model(&revenues).On("CONFLICT (month, year) DO UPDATE SET amount = excluded.amount").Exec(c.ctx)
	if err != nil {
		return fmt.Errorf("failed to upsert revenues: %w", err)
	}
	return nil
}

func generateRevenueData() []Revenue {
	return []Revenue{
		{Month: 1, Year: 2023, Amount: 1000.00},
		{Month: 2, Year: 2023, Amount: 1500.00},
		{Month: 3, Year: 2023, Amount: 2000.00},
		{Month: 4, Year: 2023, Amount: 2500.00},
		{Month: 5, Year: 2023, Amount: 3000.00},
		{Month: 6, Year: 2023, Amount: 3500.00},
	}
}

func (c *DbConfig) DropRevenueTable() error {
	_, err := c.db.NewDropTable().Model((*Revenue)(nil)).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to drop revenue table: %w", err)
	}
	return nil
}

// GetRevenueByMonthYear retrieves revenue for a specific month and year
func (c *DbConfig) GetRevenueByMonthYear(month, year int) (*Revenue, error) {
	revenue := &Revenue{}
	log.Printf("month: %d, year: %d", month, year)

	err := c.db.NewSelect().
		Model(revenue).
		Where("month = ? AND year = ?", month, year).
		Scan(c.ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no revenue found for month %d and year %d", month, year)
			return nil, fmt.Errorf("no revenue found for month %d and year %d", month, year)
		}
		log.Printf("failed to get revenue: %v", err)
		return nil, fmt.Errorf("failed to get revenue: %w", err)
	}

	log.Printf("revenue: %+v", revenue)
	return revenue, nil
}

func (c *DbConfig) Close() error {
	return c.db.Close()
}
