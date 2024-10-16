// db package is used to interact with the database
package db

// imports the necessary packages
// time package is used to manipulate time
import (
	"aio/helpers"
	"aio/logger"
	"time"
)

// CharacterMgr struct contains all the information about the character.
// It is used to manage the character in the application.
type CharacterMgr struct {
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

// NewCharacterMgr function creates a new CharacterMgr object.
// It returns a pointer to the new CharacterMgr object.
func NewCharacterMgr() *CharacterMgr {
	var birth, created, updated string
	c := &CharacterMgr{}
	row, err := get("characters_get")

	if err == nil {
		err = row.Scan(
			&c.FirstName,
			&c.LastName,
			&c.NickName,
			&birth,
			&c.MonthBudget,
			&c.Balance,
			&c.Coins,
			&c.XP,
			&c.NextLevelXP,
			&c.Level,
			&c.PP,
			&c.MaxPP,
			&c.HP,
			&c.MaxHP,
			&c.Karma,
			created,
			updated,
		)
	}

	logger.Fatal("Error scanning character from the database", err)
	c.BirthDate = helpers.TimeDBParse(birth)
	c.CreatedAt = helpers.TimeDBParse(created)
	c.UpdatedAt = helpers.TimeDBParse(updated)

	return c
}

// Death function kills the character.
// It sets all the character's stats to intial value and decreases the karma by 10.
func (c *CharacterMgr) Death() {
	c.XP = 0
	c.NextLevelXP = 50
	c.Level = 1
	c.HP = 100
	c.MaxHP = 100
	c.PP = 50
	c.MaxPP = 50
	c.Karma = c.Karma - 10
	c.Coins = 0

	err := gitFlow(func() error {
		return do("characters_death")
	})
	logger.Fatal("Error killing character", err)
}
