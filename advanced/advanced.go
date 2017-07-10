package advanced

/*
* All Callbacks use this format:
	* func (e *entity)NamingConvention() error { }
	* any error in the callback chain aborts the entire transaction and you can address/react to the error as soon as it occurs
	// BeforeSave
	// BeforeCreate
		// Save Before Associations
		// Save Self
		// Save After Associations
	// After Create
	// After Save

	// Before Delete
		// Delete
	// After Delete

	// After Find

	// There is a demo segment, but I didn't watch it
*/

/*
*Scope
	* This allows you to codify a "where" scoping rule that you can write once and apply anywhere
	* Takes the form of a function
*/

import (
	"fmt"

	// Anonymous import - package just needs to initialize in order to establish itself as a database driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"

	"github.com/annicaburns/learngorm/query"
)

// Scope demonstrates how to codify a scope
func Scope() {

	// Only seed the database once
	// query.SeedDB()

	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	appointments := []query.AppointmentQuery{}

	db.Scopes(LongMeetings).Find(&appointments)

	for _, appointment := range appointments {
		fmt.Printf("\n%v\n", appointment)
	}

}

// LongMeetings is an example of a scoping function
func LongMeetings(db *gorm.DB) *gorm.DB {
	return db.Where("length > ?", 60)
}
