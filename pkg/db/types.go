// db package type definitions
package db

import "time"

type Character struct {
	FirstName   string
	LastName    string
	NickName    string
	MonthBudget float64
	Balance     float64
	Coins       int
	XP          int
	NextLevelXP int
	Level       int
	PP          int
	MaxPP       int
	HP          int
	MaxHP       int
	Karma       int
	BirthDate   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
