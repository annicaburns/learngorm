package query

/*
* By default, Gorm does not inflate the entire graph of objects that are related to a parent entity
	* Use Eager Loading in scenarios where we want to inflate child objects
* You can select result subsets for chores like pagination
* You can shape results if you want data structures tha don't match those defined by the Go application
* Can also pass Raw SQL to the database
*/

import (
	"fmt"
	"time"

	// Anonymous import - package just needs to initialize in order to establish itself as a database driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

// RetrieveSimple demonstrates some basic query language
func RetrieveSimple() {
	// Only seed the database once
	// SeedDB()

	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	// Queries for individual records
	// Create an empty user struct
	// user := UserQuery{}
	// First asks Gorm to find the first record in that table (ordered ASC by primary key) and inflate that record into our empty user struct
	// db.First(&user)
	// FirstOrInit initializes the object provided if it doesn't find a the object provided. But it does not create the object in the database - just initializes it on the Go side.
	// db.FirstOrInit(&user, &UserQuery{Username: "lprosser"})
	// FirstOrCreate actually creates the new record in the database
	// db.FirstOrCreate(&user, &UserQuery{Username: "lprosser"})
	// Last asks Gorm to find the last record in the table (again ordered by primary key)
	// db.Last(&user)
	// fmt.Println(user)

	// Using the Find method to select multiple records (recordsets)
	// users := []UserQuery{}
	// Find method can be called with just the first parameter
	// db.Find(&users)
	// Or with the second optional where clause parameter
	// db.Find(&users, &UserQuery{Username: "fprefect"})
	// Or with a map as the second parameter to choose field names at run time because this uses database field names instead of Go property names
	// db.Find(&users, map[string]interface{}{"username": "fprefect"})
	// And finally by using straight SQL
	// db.Find(&users, "username = ?", "fprefect")

	// Where method - Use for more complex fetches
	// db.Where("username = ?", "fprefect").Find(&users)
	// Can also use the Go property names or using a map as a parameter like we did with Find
	// db.Where(&UserQuery{Username: "fprefect"}).Find(&users)
	// db.Where(map[string]interface{}{"username": "fprefect"}).Find(&users)
	// Can use the raw SQL method to find a group of users by passing a slice of username strings
	// db.Where("username in (?)", []string{"adent", "mrobot"}).Find(&users)
	// Or using Like
	// db.Where("username like ?", "%mac%").Find(&users)
	// Can also string together multiple elements in a where clause
	// db.Where("username like ? and first_name = ?", "%e%", "Ford").Find(&users)
	// Can query on time
	// db.Where("created_at < ?", time.Now()).Find(&users)
	// db.Where("created_at between ? and ?", time.Now().Add(-30*24*time.Hour), time.Now()).Find(&users)
	// Not method - The Where clause is only looking for positive matches, so use the Not method for the reverse
	// db.Not("username = ?", "adent").Find(&users)
	// Or method. Chain this on to a where clause to combine two different fetches
	// db.Where("username = ?", "fprefect").Or("username = ?", "tmacmillan").Find(&users)

	// Gorm defaults to lazy loading to prefer speed over convenience
	// Preload method indicates to gorm which child objects we want to inflate (load)
	// db.Preload("CalendarQuery").Find(&users)
	// Can take this as deep into the object graph as desired
	// Just make sure any preloaded fields are indexed or the query will be very slow
	// Also note in the console dump that the AppointmentQueries - which are being stored as pointers - print out the memory address
	// db.Preload("CalendarQuery.AppointmentQuerys").Find(&users)

	// Limiting and ordering data
	// Keep in mind that .Find is always called AFTER the constraints are described. It basically says "ok, Gorm, go do something"
	// Limit method is like top
	// db.Limit(2).Find(&users)
	// Order method can be combine with Limit
	// db.Limit(3).Order("first_name").Find(&users)
	// To use pagination, we also need to use offset. Remember all shaping statements need to come before .Find.
	// The order of the "shapers" (Limit and Offset in this case) doesn't matter
	// db.Limit(2).Offset(2).Order("first_name").Find(&users)

	// Selecting Data Subsets - inflating only certain properties in a returned object
	// Select method takes a slice of strings that indicate the properties desired
	// Be judicious with Select because it confuses end users of the API and quickly becomes a maintenance problem
	// db.Select([]string{"first_name", "last_name"}).Find(&users)

	// for _, user := range users {
	// This will print the user object data but will not inflate the child calendar objects
	// Gorm is lazy about what it loads and will only inflate child objects if explicitly requested
	// The empty calendar object will show up after each user in the console
	// 	fmt.Printf("\n%v\n", user)
	// }

	// Plucking is used to select certain object properties and store them into a structure
	/*
		usernames := []string{}
		db.Model(&UserQuery{}).Pluck("username", &usernames)
		for _, username := range usernames {
			fmt.Printf("\n%v\n", username)
		}
	*/
	// Can also dynamically load the plucked data into a new structure to make things more clear to the end user
	// Use the Select method to select the fields
	// Use Scan instead of Find in this case because Find expects a parameter that is directly mapped to an underlying table
	/*
		userVMs := []UserViewModel{}
		db.Model(&UserQuery{}).Select([]string{"first_name", "last_name"}).Scan(&userVMs)
		for _, user := range userVMs {
			fmt.Printf("\n%v\n", user)
		}
	*/

	// Count method lets you find out how many objects are in a table that meet a certain criteria without inflating all those objects
	/*
		var userCount int
		db.Model(&UserQuery{}).Count(&userCount)
		fmt.Println(userCount)
	*/

	// Attrs method is for cases where you want to provide defaults or overide data without constraining the where clause
	// Combine a Where method with a call to .Attrs to provide a default with which to inflate an object if .FirstOrInit finds nothing
	/*
		user := &UserQuery{}
		db.Where("username = ?", "goober").Attrs(&UserQuery{FirstName: "Eddie", LastName: "Nothing"}).FirstOrInit(&user)
		fmt.Printf("\n%v\n", user)
	*/
	// Assign method overides just specific properties
	// But remember this is only in the application layer - didn't change anything in the underlying database
	user := &UserQuery{}
	db.Where("username = ?", "adent").Assign(&UserQuery{FirstName: "Eddie"}).FirstOrInit(&user)
	fmt.Printf("\n%v\n", user)
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

// UserViewModel is specific to this class file
type UserViewModel struct {
	FirstName string
	LastName  string
}
