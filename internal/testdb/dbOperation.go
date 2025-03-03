package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseName = "test.db"
)

type User struct {
	ID   int
	Name string
	Age  int
}

func DbOperation() {
	log.Printf("DbOperation")

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable(db)

	user := User{Name: "Chetan", Age: 32}

	id, err := createUser(db, user)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Inserted user with ID:", id)

	user, err = readUser(db, id)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Read User:%+v\n", user)

	// err = updateUser(db, 30, "Mahendra Singh")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Updated User:%+v\n", user)

	//update user by Age

	err = updateUserByAge(db, 25, "Mahendra Singh")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Updated User by Age:%+v\n", user)

	// err = deleteUser(db, 9)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func openDB() (*sql.DB, error) {
	return sql.Open("sqlite3", databaseName)
}

func createTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, age INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
}

func createUser(db *sql.DB, user User) (int64, error) {
	result, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func readUser(db *sql.DB, id int64) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, name, age FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		return user, err
	}
	return user, nil
}

// func updateUser(db *sql.DB, id int64, name string) error {
// 	_, err := db.Exec("UPDATE users SET name = ? WHERE age = ?", name, id)
// 	return err
// }

//how to give new value to updateUser()

func updateUserByAge(db *sql.DB, age int, name string) error {
	_, err := db.Exec("UPDATE users SET age = ? WHERE name = ?", age, name)
	return err
}

func deleteUser(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
