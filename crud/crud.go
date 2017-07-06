package crud

import (
	"time"
	// Anonymous import - package just needs to initialize in order to establish itself as a database driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

// CreateWithChildRecords demonstrates some basic create, update, delete operations
func CreateWithChildRecords() {
	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	db.DropTableIfExists(&CruddyUser{})
	db.CreateTable(&CruddyUser{})
	db.DropTableIfExists(&CruddyAppointment{})
	db.CreateTable(&CruddyAppointment{})

	user := CruddyUser{
		FirstName: "Arthur",
		LastName:  "Dent",
	}

	appointments := []CruddyAppointment{
		CruddyAppointment{Subject: "First"},
		CruddyAppointment{Subject: "Second", Attendees: []*CruddyUser{&user}},
		CruddyAppointment{Subject: "Third"},
	}

	user.Appointments = appointments

	// It may not be obvious in production code if the record you want to create exists yet - NewRecord returns a bool with the answer
	if db.NewRecord(&user) {
		db.Create(&user)
	}

	// fmt.Println(db.NewRecord(&user))

}

// CruddyUser is specific to this class file
type CruddyUser struct {
	gorm.Model
	FirstName    string
	LastName     string
	Appointments []CruddyAppointment
}

// CruddyAppointment is specific to this class file
type CruddyAppointment struct {
	gorm.Model
	CruddyUserID uint
	Subject      string
	Description  string
	StartTime    *time.Time
	Length       uint
	Attendees    []*CruddyUser
}
