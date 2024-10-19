// db package is used to interact with the database
package db

import (
	"aio/pkg/log"
	"aio/pkg/utils/tm"
)

// NewCharacterMgr function creates a new CharacterMgr object.
// It returns a pointer to the new CharacterMgr object.
func CharGet() (*Character, error) {
	var birth, created, updated string
	c := &Character{}
	row, err := get("characters_get")
	if err != nil {
		log.Err("failed to get the character")
		return nil, err
	}

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

	if err != nil {
		log.Err("failed to scan the character")
		return nil, err
	}

	c.BirthDate, err = tm.DBParse(birth)
	if err != nil {
		log.Err("failed to parse the birth date")
		return nil, err
	}

	c.CreatedAt, err = tm.DBParse(created)
	if err != nil {
		log.Err("failed to parse the created date")
		return nil, err
	}

	c.UpdatedAt, err = tm.DBParse(updated)
	if err != nil {
		log.Err("failed to parse the updated date")
		return nil, err
	}

	return c, nil
}

// Death function kills the character.
// It sets all the character's stats to intial value and decreases the karma by 10.
func (c *Character) Death() error {
	c.XP = 0
	c.NextLevelXP = 50
	c.Level = 1
	c.HP = 100
	c.MaxHP = 100
	c.PP = 50
	c.MaxPP = 50
	c.Karma = c.Karma - 10
	c.Coins = 0

	err := do("characters_death")
	if err != nil {
		log.Err("failed to kill the character")
		return err
	}

	return nil
}
