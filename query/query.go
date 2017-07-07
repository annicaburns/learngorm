package query

/*
* By default, Gorm does not inflate the entire graph of objects that are related to a parent entity
	* Use Eager Loading in scenarios where we want to inflate child objects
* You can select result subsets for chores like pagination
* You can shape results if you want data structures tha don't match those defined by the Go application
* Can also pass Raw SQL to the database
*/

import (
	"time"

	// Anonymous import - package just needs to initialize in order to establish itself as a database driver
	_ "github.com/go-sql-driver/mysql"

	"fmt"

	"github.com/jinzhu/gorm"
)

// RetrieveSimple demonstrates some basic query language
func RetrieveSimple() {

	// SeedDB()

	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	// Basic Query for first record in a table
	// Create an empty user struct
	user := UserQuery{}
	// First asks Gorm to find the first record in that table (ordered ASC by primary key) and inflate that record into our empty user struct
	// db.First(&user)
	// FirstOrInit initializes the object provided if it doesn't find a the object provided. But it does not create the object in the database - just initializes it on the Go side.
	// db.FirstOrInit(&user, &UserQuery{Username: "lprosser"})
	// FirstOrCreate actually creates the new record in the database
	// db.FirstOrCreate(&user, &UserQuery{Username: "lprosser"})
	// Last asks Gorm to find the last record in the table (again ordered by primary key)
	db.Last(&user)
	fmt.Println(user)
}

// SeedDB can be used from any package
func SeedDB() {
	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	db.DropTableIfExists(&UserQuery{})
	db.CreateTable(&UserQuery{})
	db.DropTableIfExists(&CalendarQuery{})
	db.CreateTable(&CalendarQuery{})
	db.DropTableIfExists(&AppointmentQuery{})
	db.CreateTable(&AppointmentQuery{})

	users := map[string]*UserQuery{
		"adent":       &UserQuery{Username: "adent", FirstName: "Arthur", LastName: "Dent"},
		"fprefect":    &UserQuery{Username: "fprefect", FirstName: "Ford", LastName: "Prefect"},
		"tmacmillan":  &UserQuery{Username: "tmacmillan", FirstName: "Tricia", LastName: "Macmillan"},
		"zbeeblebrox": &UserQuery{Username: "zbeeblebrox", FirstName: "Zaphod", LastName: "Beeblebrox"},
		"mrobot":      &UserQuery{Username: "mrobot", FirstName: "Marvin", LastName: "Robot"},
	}

	for _, user := range users {
		user.CalendarQuery = CalendarQuery{Name: "Calendar"}
	}

	users["adent"].AddAppointment(&AppointmentQuery{
		Subject:   "Save House",
		StartTime: parseTime("1979-07-02 08:00"),
		Length:    60,
	})
	users["fprefect"].AddAppointment(&AppointmentQuery{
		Subject:   "Get a drink at a local pub",
		StartTime: parseTime("1979-07-02 10:00"),
		Length:    11,
		Attendees: []*UserQuery{users["adent"]},
	})
	users["fprefect"].AddAppointment(&AppointmentQuery{
		Subject:   "Hitch a ride",
		StartTime: parseTime("1979-07-02 10:12"),
		Length:    60,
		Attendees: []*UserQuery{users["adent"]},
	})
	users["fprefect"].AddAppointment(&AppointmentQuery{
		Subject:   "Attend a poetry reading",
		StartTime: parseTime("1979-07-02 11:00"),
		Length:    30,
		Attendees: []*UserQuery{users["adent"]},
	})
	users["fprefect"].AddAppointment(&AppointmentQuery{
		Subject:   "Get thrown into Space",
		StartTime: parseTime("1979-07-02 11:40"),
		Length:    5,
		Attendees: []*UserQuery{users["adent"]},
	})
	users["fprefect"].AddAppointment(&AppointmentQuery{
		Subject:   "Get saved from Space",
		StartTime: parseTime("1979-07-02 11:45"),
		Length:    1,
		Attendees: []*UserQuery{users["adent"]},
	})
	users["zbeeblebrox"].AddAppointment(&AppointmentQuery{
		Subject:   "Explore Planet Builder's Homeworld",
		StartTime: parseTime("1979-07-03 11:00"),
		Length:    240,
		Attendees: []*UserQuery{users["adent"], users["fprefect"], users["tmacmillan"], users["mrobot"]},
	})

	for _, user := range users {
		db.Save(&user)
	}
}

func parseTime(rawTime string) time.Time {
	// Apparently it has to be this exact date ???? WTF ???
	const timeLayout = "2006-01-02 15:04"
	t, _ := time.Parse(timeLayout, rawTime)
	return t
}

// UserQuery is specific to this class file
type UserQuery struct {
	gorm.Model
	Username      string
	FirstName     string
	LastName      string
	CalendarQuery CalendarQuery
}

// AddAppointment is a helper function
func (user *UserQuery) AddAppointment(appointment *AppointmentQuery) {
	user.CalendarQuery.AppointmentQuerys = append(user.CalendarQuery.AppointmentQuerys, appointment)
}

// CalendarQuery is specific to this class file
type CalendarQuery struct {
	gorm.Model
	Name              string
	UserQueryID       uint
	AppointmentQuerys []*AppointmentQuery
}

// AppointmentQuery is specific to this class file
type AppointmentQuery struct {
	gorm.Model
	Subject         string
	Description     string
	StartTime       time.Time
	Length          uint
	CalendarQueryID uint
	Attendees       []*UserQuery `gorm:"many2many:appointment_query_user_query"`
}
