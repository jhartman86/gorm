package tests_test

import (
	"testing"

	"gorm.io/gorm"
)

func TestComputedCols(t *testing.T) {
	type (
		WithComputedCols struct {
			gorm.Model
			FirstName string `gorm:"type:varchar;"`
			LastName  string `gorm:"type:varchar;"`
			// FullName  string `gorm:"->;type:varchar GENERATED ALWAYS AS (first_name || ' ' || last_name) STORED"`
			FullName string `gorm:"type:varchar;generated:ALWAYS AS (first_name || ' ' || last_name) STORED;"`
		}
	)

	if err := DB.Migrator().DropTable(WithComputedCols{}); err != nil {
		t.Fatalf("Failed to drop table")
	}

	if err := DB.Debug().AutoMigrate(WithComputedCols{}); err != nil {
		t.Fatalf("Failed to auto migrate, but got error %v", err)
	}

	t.Run(`WithoutValue`, func(t *testing.T) {
		r := WithComputedCols{FirstName: `lorem`, LastName: `ipsum`}
		if err := DB.Debug().Create(&r).Error; err != nil {
			t.Fatalf("Failed creating fixture record, err %v", err)
		}
		if r.FullName != `lorem ipsum` {
			t.Fatalf("Failed on create return")
		}
	})

	t.Run(`WithValue`, func(t *testing.T) {
		r := WithComputedCols{FirstName: `lorem`, LastName: `ipsum`, FullName: `ShouldNotAttemptToWrite`}
		if err := DB.Debug().Create(&r).Error; err != nil {
			t.Fatalf("Failed creating fixture record, err %v", err)
		}
		if r.FullName != `lorem ipsum` {
			t.Fatalf("Failed on create return")
		}
	})

	// reloaded := WithComputedCols{}
	// if err := DB.First(&reloaded, r.ID).Error; err != nil {
	// 	t.Fatalf("Failed creating fixture record, err %v", err)
	// }
	// if reloaded.FullName != `lorem ipsum` {
	// 	t.Fatalf("Computed column does not match")
	// }
	// if r.FullName != reloaded.FullName {
	// 	t.Fatalf("Computed column not returned during onCreate")
	// }
}
