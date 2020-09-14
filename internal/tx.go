package internal

import (
	"database/sql/driver"
	"log"
)

type Tx struct {
	driver.Tx
}

func (c *Tx) Commit() (err error) {
	log.Println("Commit")

	return c.Tx.Commit()
}

func (c *Tx) Rollback() error {
	log.Println("Rollback")

	return c.Tx.Rollback()
}
