package models

import "time"

type DeductionConfigEntity struct {
	ID        int       `postgres:"id"`
	Personal  float64   `postgres:"personal"`
	KReceipt  float64   `postgres:"kreceipt"`
	CreatedAt time.Time `postgres:"created_at"`
	UpdatedAt time.Time `postgres:"updated_at"`
}

type DeductionConfig struct {
	ID        int       `postgres:"id"`
	Personal  float64   `postgres:"personal"`
	KReceipt  float64   `postgres:"kreceipt"`
	CreatedAt time.Time `postgres:"created_at"`
	UpdatedAt time.Time `postgres:"updated_at"`
}
