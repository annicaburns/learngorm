package dbSchema

import (
	"fmt"

	"github.com/jinzhu/gorm"

	// anonymous import - package just needs to initialize in order to establish itself as a database driver
	_ "github.com/go-sql-driver/mysql"
)

// BasicMethods demonstrates GORM basic schema and query methods
func BasicMethods() {
	// gorm.Open() to connect to a database

	// other methods operate on the instance of the database connection
	//.CreateTable()
	//.DropTable()
	//.DropTableIfExists
	//.Create (dbase record)
	//.Save (dbase record) - update
	//.Where (dbase record)
	//.First (dbase record)
	//.Last (dbase record)
	//.Delete (dbase record)

	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	db.DropTableIfExists(&User{})
	db.CreateTable(&User{})

	for _, user := range users {
		db.Create(&user)
	}

	// u := User{Username: "tmacmillan"}
	// db.Where(&u).First(&u)
	// fmt.Println(u)
	// u.LastName = "Beeblebrox"
	// db.Save(&u)

	// user := User{}
	// db.Where(&u).First(&user)
	// fmt.Println(user)

	// db.Where(&User{Username: "adent"}).Delete(&User{})

	fmt.Println("done")

}

// CustomizationMethods demonstrates methods to customize schema and do something outside of the GO conventions
func CustomizationMethods() {
	//.SingularTable(true)
	// do not pluralize the table name during table creations.
	// Without this customization, struct User would create a table called "users"

	// Implement the TableName method shown below to produce a name that doesn't match your model name

	// See gorm and sql tags in the User struct definition for more ways to customize schema
}

// TableName method is used to customize how GO names yourn table.
// Implement the TableName function to produce a name other than the model name
// func (u User) TableName() string {
// 	return "alternate_name"
// }

// User is our user object
type User struct {
	// If a field is named "ID" gorm will interpret this as a primary key and will make it auto-incrementing and create the PK constratints
	// ID uint
	// Or use an field name with the gorm "primary_key" tag
	// UserID uint `gorm:"primary_key"`
	// embed gorm.Model struct to compose the default GORM fields into your model (id, created_at, updated_at, deleted_at)
	// GORM will automatically keep these fields up to date during inserts and updates and deletes and will use a "soft delete" model
	gorm.Model
	// capitalized fields will show up in the database in all lowercase
	// use 'type' tag to specify a type other than the default (default string type is VARCHAR(255))
	// use 'not null' tag to add a null constraint
	Username string `sql:"type:VARCHAR(15);not null"`
	// pascal-cased names like this will be converted to this format: first_name
	// use 'size' tag to accept the default type, but use a different size
	// use the DEFAULT tag to provide a default in cases where the user hasn't
	FirstName string `sql:"size:100;not null;DEFAULT:'Annica'"`
	// To exert control over the exact name of the field, use the gorm "column" tag
	// use 'unique' tag to add a unique constraint, and 'unique_index' to add a unique index to take advantage of the unique constraint
	// although these seemed to be redundant in MySQL - we ended up with two unique indexes so far as I could tell in workbench
	// If you are only needing to enforce uniqueness as a business rule rather than the field needing to be searched or used for sorting then I'd use the constraint, again to make the intended use more obvious when someone else looks at your table definition.
	// Note that if you use both a unique constraint and a unique index on the same field the database will not be bright enough to see the duplication, so you will end up with two indexes which will consume extra space and slow down row inserts/updates.
	LastName string `sql:"unique" gorm:"column:LastName"`
	// GORM AUTO_INCREMENT tag - this didn't work for me with MySQL
	Count int `gorm:"AUTO_INCREMENT"`
	// Indicate a Temp field that will be created and used but persisted to the database schema - essentially ignored by GORM
	TempField bool `gorm:"-"`
	// To add a unique constratint
}

var users = []User{
	User{Username: "adent", FirstName: "Arthur", LastName: "Dent"},
	User{Username: "fprefect", FirstName: "Ford", LastName: "Prefect"},
	User{Username: "tmacmillan", FirstName: "Tricia", LastName: "Macmillan"},
	User{Username: "mrobot", FirstName: "Marvin", LastName: "Robot"},
	User{Username: "smith", LastName: "Smith"},
}

// EmbedChildObjects demonstrates embedding child objects
func EmbedChildObjects() {
	db, err := gorm.Open("mysql", "gorm:gorm@tcp(localhost:23306)/gorm")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	db.DropTableIfExists(&CleanUser{})
	db.CreateTable(&CleanUser{})

	// Call .AddIndex method on the database instance to manually create indexes so we can name them and include multiple columns
	// And because we know the names, we can remove them with the .RemoveIndex method
	// Column names are the way the database sees them
	// the Model method scopes any subsequent calls to the table represented by the empty interface supplied as the argument
	db.Model(&CleanUser{}).AddIndex("idx_first_name", "first_name")

	for _, f := range db.NewScope(&CleanUser{}).Fields() {
		fmt.Println(f.Name)
	}

}

// CleanUser is used to demonstrate embedding child object
type CleanUser struct {
	// use gorm "embedded" tag to implement a psuedo inheritance - flattens the fields of the embedded child object into the parent object
	Model     gorm.Model `gorm:"embedded"`
	FirstName string
	LastName  string
}

var cleanUsers = []CleanUser{
	CleanUser{FirstName: "Arthur", LastName: "Dent"},
	CleanUser{FirstName: "Ford", LastName: "Prefect"},
	CleanUser{FirstName: "Tricia", LastName: "Macmillan"},
}
