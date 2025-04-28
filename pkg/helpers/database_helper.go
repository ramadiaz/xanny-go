package helpers

import (
	"gorm.io/gorm"
)

func CommitOrRollback(tx *gorm.DB) {
	err := recover()
	if err != nil || tx.Error != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}