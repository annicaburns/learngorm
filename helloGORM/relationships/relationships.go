package relationships

import (
	"time"

	"github.com/jinzhu/gorm"

	// anonymous import - package just needs to initialize in order to establish itself as a database driver

	_ "github.com/go-sql-driver/mysql"
)

// BasicRelationships demonstrates some basic principles of relationships in GORM
func BasicRelationships() {
	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm")
	if err != nil {
		panic(err.Error())
	}

	db.DropTableIfExists(&Appointment{})
	db.CreateTable(&Appointment{})
	db.DropTableIfExists(&Calendar{})
	db.CreateTable(&Calendar{})
	db.DropTableIfExists(&User{})
	db.CreateTable(&User{})

	// .Debug method logs the SQL statements as they are being made

	db.Debug().Model(&Calendar{}).
		AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	db.Debug().Save(&User{
		Username: "adent",
		Calendar: Calendar{
			Name: "Improbable Events",
			Appointments: []Appointment{
				{Subject: "Spontaneous Whale Generation", StartTime: time.Now()},
				{Subject: "Saved from Vacuum of Space", StartTime: time.Now()},
			},
		},
	})

	// To query the user and calendar records we just created, create empty structs and inflate them with the .FirstMethod
	// u := User{}
	// c := Calendar{}
	// db.First(&u)
	// note that by default GORM does not inflate child objects without extra work being done
	// fmt.Println(u)
	// fmt.Println()
	// fmt.Println(c)

}

// User is used to demonstrate relationship stuff
type User struct {
	gorm.Model
	Username  string
	FirstName string
	LastName  string
	// This embedded Calendar object, along with the UserID field in the Calendar table establishes a HasOne relationship (one-to-one)
	// If we want an OwnedBy one-to-one relationship, add a CalendarID field here and remove the UserID field in the Calendar table
	Calendar Calendar
}

// Calendar is used to demonstrate a one-to-one relationship with User
type Calendar struct {
	gorm.Model
	Name string
	// Named this way, GORM can infer during inserts that this field is a foreign key for the User table
	// But GORM won't automatically add a FK constraint >> That has to be explicitly added.
	UserID       uint
	Appointments []Appointment
}

// Appointment is used to demonstrate a many-to-one relationship with Calendar
type Appointment struct {
	gorm.Model
	Subject     string
	Description string
	StartTime   time.Time
	Length      uint
	CalendarID  uint
	// Creating a many-to-many relationship by specifying the name of the lookup table that we want GORM to create.
	// In this case: "appointment_user"
	Attendees []User `gorm:"many2many:appoinment_user"`
}
