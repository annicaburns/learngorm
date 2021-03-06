package crud

import (
	"time"
	// Anonymous import - package just needs to initialize in order to establish itself as a database driver
	_ "github.com/go-sql-driver/mysql"

	"fmt"

	"github.com/jinzhu/gorm"
)

// CreateWithChildRecords demonstrates some basic operations
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

// UpdateRecords demonstrates some update operations
func UpdateRecords() {
	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	db.DropTableIfExists(&CruddyUser2{})
	db.CreateTable(&CruddyUser2{})

	user := CruddyUser2{
		FirstName: "Arthur",
		LastName:  "Dent",
	}
	db.Create(&user)
	fmt.Println(user)
	fmt.Println()

	// user.FirstName = "Zaphod"
	// user.LastName = "Beeblebrox"
	// db.Save(&user)
	// fmt.Println(user)

	// Or use reflection to iterate through properties and use the .Update method dynamically change the values at run time of a new field
	// Also note that we are scoping the model here to a specific user object
	// db.Model(&user).Update("first_name", "Marker")

	// Can also use .Updates with a map to update muliple properties at once
	db.Model(&user).Updates(
		map[string]interface{}{
			"first_name": "Zaphod",
			"last_name":  "Beeblebrox",
		})

	fmt.Println(user)

}

// BatchUpdates demonstrates scoping to a particular table and making a bunch of updates
func BatchUpdates() {
	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	db.DropTableIfExists(&CruddyUser3{})
	db.CreateTable(&CruddyUser3{})

	db.Create(&CruddyUser3{
		FirstName: "Tricia",
		LastName:  "Dent",
		Salary:    50000,
	})

	db.Create(&CruddyUser3{
		FirstName: "Arthur",
		LastName:  "Dent",
		Salary:    30000,
	})

	// The Table method scopes calls to a particular table - but we must speak about things (and call things) from the database perspective
	// The Model method uses the semantics (and naming) of GO
	db.Table("cruddy_user3").Where("last_name = ?", "Dent").Update("last_name", "Macmillan-Dent")

	db.Table("cruddy_user3").Where("salary > ?", 40000).Update("salary", gorm.Expr("salary + 5000"))

}

// DeleteRecords demonstrates hard and soft deletes
func DeleteRecords() {
	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	db.DropTableIfExists(&CruddyUser4{})
	db.CreateTable(&CruddyUser4{})
	db.DropTableIfExists(&CruddyUser2{})
	db.CreateTable(&CruddyUser2{})

	// Note that the CruddyUser4 model has an ID field instead of the GORM model fields - a hard delete will occur
	user := CruddyUser4{
		FirstName: "Ford",
		LastName:  "Prefect",
	}

	db.Create(&user)

	// By passing in the user object, GORM is smart enough to write a SQL statement based on the properties on that user object
	// In this case it will know that we are using an object with an ID field and it's where clause will use the ID
	db.Debug().Delete(&user)

	// A soft delete will happen here because we are deleting a model object that includes the GORM model fields instead of an ID field
	modelUser := CruddyUser2{
		FirstName: "Ford",
		LastName:  "Prefect",
	}

	db.Create(&modelUser)
	db.Debug().Delete(&modelUser)

	// GORM considers this soft deleted record as truly deleted and will not return it from most queries
	user2 := CruddyUser2{}
	db.First(&user2)
	fmt.Println()
	fmt.Println(user2)

}

// BatchDeletes demonstrates the obvious
func BatchDeletes() {
	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	db.DropTableIfExists(&CruddyUser2{})
	db.CreateTable(&CruddyUser2{})

	db.Create(&CruddyUser2{
		FirstName: "Tricia",
		LastName:  "Macmillan-Dent",
	})

	db.Create(&CruddyUser2{
		FirstName: "Arthur",
		LastName:  "Dent",
	})

	db.Where("last_name LIKE ?", "Mac%").Delete(&CruddyUser2{})

}

// Transactions demonstrates the obvious
func Transactions() {
	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	db.DropTableIfExists(&CruddyUser2{})
	db.CreateTable(&CruddyUser2{})

	user := CruddyUser2{
		FirstName: "Marvin",
		LastName:  "Robot",
	}

	// Create a transaction object
	transaction := db.Begin()

	if err = transaction.Create(&user).Error; err != nil {
		transaction.Rollback()
	}

	user.LastName = "The Happy Robot"
	// Intentionally rollback the transaction
	if err = transaction.Save(&user).Error; err == nil {
		transaction.Rollback()
	}

	transaction.Commit()
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

// CruddyUser2 is used with the UpdateRecordsFunction
type CruddyUser2 struct {
	gorm.Model
	FirstName string
	LastName  string
}

// BeforeUpdate method is a built in callback that fires before the .Save, .Update and .Updates methods are called
// Use .Update or .UpdateColumns if we DON'T want to trigger these callbacks - these methods don't fire them
func (user *CruddyUser2) BeforeUpdate() error {
	fmt.Println("Before Update")

	return nil
}

// AfterUpdate method is the same, just firing right after the updates occur
func (user *CruddyUser2) AfterUpdate() error {
	println("After Update")
	return nil
}

// CruddyUser3 is used with BatchUpdates function
type CruddyUser3 struct {
	gorm.Model
	FirstName string
	LastName  string
	Salary    uint
}

// CruddyUser4 is used with DeleteRecords function
type CruddyUser4 struct {
	ID        uint
	FirstName string
	LastName  string
}
