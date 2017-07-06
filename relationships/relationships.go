package relationships

import (
	"fmt"

	"time"

	"github.com/jinzhu/gorm"

	// Anonymous import - package just needs to initialize in order to establish itself as a database driver
	_ "github.com/go-sql-driver/mysql"
)

// BasicRelationships demonstrates some basic principles of relationships in GORM
func BasicRelationships() {
	// Need to include the ?parseTime=true parameter to connection string if you want the gorm model fields to update updated_at when anything related to the object changes
	// https://github.com/jinzhu/gorm/issues/958
	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	db.DropTableIfExists(&Appointment{})
	db.CreateTable(&Appointment{})
	db.DropTableIfExists(&Calendar{})
	db.CreateTable(&Calendar{})
	db.DropTableIfExists(&RelationshipUser{})
	db.CreateTable(&RelationshipUser{})
	db.DropTableIfExists(&TaskList{})
	db.CreateTable(&TaskList{})
	// .Debug method logs the SQL statements as they are being made
	db.Debug().Model(&Calendar{}).
		AddForeignKey("relationship_user_id", "relationship_users(id)", "CASCADE", "CASCADE")

	users := []RelationshipUser{
		{Username: "fprefect"},
		{Username: "tmacmillan"},
		{Username: "mrobot"},
	}
	for i := range users {
		db.Save(&users[i])
	}
	// Interestingly... as we add each appointment with it's associated list of users, the updated_at field of each user will get updated because something "related to the user" has changed.
	db.Save(&RelationshipUser{
		Username: "adent",
		Calendar: Calendar{
			Name: "Improbable Events",
			Appointments: []Appointment{
				{Subject: "Spontaneous Whale Generation", Description: "easy", StartTime: time.Now(), Attendees: users},
				{Subject: "Saved from Vacuum of Space", Description: "hard", StartTime: time.Now(), Attendees: users},
			},
		},
		TaskList: TaskList{
			Name: "Urgent ToDos",
			Appointments: []Appointment{
				{Subject: "Submit Expenses", Description: "easy", StartTime: time.Now()},
				{Subject: "Jira Work", Description: "hard", StartTime: time.Now()},
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

// ModelAssociationMethod demonstrates capabilities of the Model function's Association Method
func ModelAssociationMethod() {
	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	BasicRelationships()

	// The parameter for the Association method scopes subsequent code to the named field of the scoped table (Calendar in this case)
	db.Model(&Calendar{}).Association("Appointments")
	appointments := []Appointment{}
	db.Find(&appointments)
	fmt.Println(appointments)
	// .Find(&appointments)
	// .Append(&appointments)
	// .Delete(&appointments)
	// .Replace(&appointments)
	// .Count()
	// .Clear()
}

// RelationshipUser is used to demonstrate relationship stuff
type RelationshipUser struct {
	gorm.Model
	Username  string
	FirstName string
	LastName  string
	// This embedded Calendar object, along with the RelationshipUserID field in the Calendar table establishes a HasOne relationship (one-to-one)
	// If we want an OwnedBy one-to-one relationship, add a CalendarID field here and remove the RelationshipUserID field in the Calendar table
	Calendar Calendar
	TaskList TaskList
}

// Calendar is used to demonstrate a one-to-one relationship with RelationshipUser
type Calendar struct {
	gorm.Model
	Name string
	// Named this way, GORM can infer during inserts that this field is a foreign key for the RelationshipUser table
	// But GORM won't automatically add a FK constraint >> That has to be explicitly added.
	RelationshipUserID uint
	Appointments       []Appointment `gorm:"polymorphic:owner"`
}

// Appointment is used to demonstrate a many-to-one relationship with Calendar
type Appointment struct {
	gorm.Model
	Subject     string
	Description string
	StartTime   time.Time
	Length      uint
	// Go's version of polymorphism. This allows us to map EITHER the id of a Calendar OR the id of a TaskList into the same field
	OwnerID   uint
	OwnerType string
	// Creating a many-to-many relationship by specifying the name of the lookup table that we want GORM to create.
	// In this case: "appointment_user"
	Attendees []RelationshipUser `gorm:"many2many:appoinment_user"`
}

// TaskList is a second container for Appointments - demonstrating Go's version of polymorphism
type TaskList struct {
	gorm.Model
	Name               string
	RelationshipUserID uint
	Appointments       []Appointment `gorm:"polymorphic:owner"`
}
