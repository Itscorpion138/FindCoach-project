package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Users struct {
	ID         uuid.UUID
	Name       string
	LastName   string
	Age        int
	Height     int
	Weight     float64
	Gender     string
	Plan       string
	Created    time.Time
	SkillLevel string
}

func main() {
	err := godotenv.Load("database.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connstr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("SSL_MODE"),
	)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	createInfoTable(db)
	fmt.Print("What operation you want to do today? (addUser/deleteUser/editUser): ")
	var operation string
	fmt.Scanln(&operation)
	switch operation {
	case "addUser":
		newUser(db)
	case "deleteUser":
		fmt.Print("please enter the id: ")
		var IDstr string
		fmt.Scanln(&IDstr)
		ID, err := uuid.Parse(IDstr)
		if err != nil {
			log.Fatal(err)
		}
		deleteUserDataByID(db, ID)
	case "editUser":
		fmt.Print("please enter the id: ")
		var IDstr string
		fmt.Scanln(&IDstr)
		ID, err := uuid.Parse(IDstr)
		if err != nil {
			log.Fatal(err)
		}
		editUserDataByID(db, ID)
	}

	// createInfoTable(db)

}

func createInfoTable(db *sql.DB) {
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	if err != nil {
		log.Fatal(err)
	}

	query := `CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL,
    lastname TEXT NOT NULL,
    age INT NOT NULL,
    height INT NOT NULL,
    weight DECIMAL(5,2) NOT NULL,  -- allows two decimal places, e.g. 72.35
    gender TEXT NOT NULL CHECK (gender IN ('Male','Female')),
    plan TEXT NOT NULL CHECK (plan IN ('normal','semi-interactive','fully-interactive')),
    created TIMESTAMP DEFAULT now(),
    skilllevel TEXT NOT NULL CHECK (skilllevel IN ('Beginner','Intermediate','Master'))
);`

	_, err2 := db.Exec(query)
	if err2 != nil {
		log.Fatal(err2)
	}
}
func insertUserData(db *sql.DB, user Users) uuid.UUID {
	var pk uuid.UUID

	query := `
    INSERT INTO users (name, lastname, age, height, weight, gender, skilllevel, plan)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id`

	err := db.QueryRow(query, user.Name, user.LastName, user.Age, user.Height, user.Weight, user.Gender, user.SkillLevel, user.Plan).Scan(&pk)
	if err != nil {
		log.Fatal(err)
	}
	return pk
}
func deleteUserDataByID(db *sql.DB, id uuid.UUID) {

	query := `DELETE FROM users
	WHERE id = $1`

	_, err := db.Exec(query, id)
	if err != nil {
		log.Fatal(err)
	}
}
func editUserDataByID(db *sql.DB, id uuid.UUID) {
	fmt.Print("What value you want to change: (Name, LastName, Age, Height, Weight, SkillLevel, plan, gender)\n")
	var column string
	fmt.Scanln(&column)

	// Whitelist valid columns
	validColumns := map[string]bool{
		"Name": true, "Lastname": true, "Age": true, "Height": true,
		"Weight": true, "Skilllevel": true, "Gender": true, "Plan": true,
	}

	if !validColumns[column] {
		log.Fatalf("Invalid column name: %s", column)
	}

	fmt.Print("What do you want to change it to? ")
	var newValue string
	fmt.Scanln(&newValue)

	// Validate if applicable
	validGenders := map[string]bool{"Male": true, "Female": true}
	validSkills := map[string]bool{"Beginner": true, "Intermediate": true, "Master": true}
	validPlans := map[string]bool{"normal": true, "semi-interactive": true, "fully-interactive": true}

	switch column {
	case "Gender":
		if !validGenders[newValue] {
			log.Fatal("Invalid gender")
		}
	case "Skilllevel":
		if !validSkills[newValue] {
			log.Fatal("Invalid skill level")
		}
	case "Plan":
		if !validPlans[newValue] {
			log.Fatal("Invalid plan")
		}
	}

	query := fmt.Sprintf(`UPDATE users
SET %s = $1
WHERE id = $2;`, column)

	_, err := db.Exec(query, newValue, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("successfully changed.\n")
}
func newUser(db *sql.DB) {

	for {
		var user Users

		fmt.Println("Enter new user details:")

		fmt.Print("Name: ")
		fmt.Scanln(&user.Name)
		fmt.Print("Last Name: ")
		fmt.Scanln(&user.LastName)
		fmt.Print("Age: ")
		fmt.Scanln(&user.Age)
		fmt.Print("Height: ")
		fmt.Scanln(&user.Height)
		fmt.Print("Weight: ")
		fmt.Scanln(&user.Weight)
		fmt.Print("Gender: ")
		fmt.Scanln(&user.Gender)
		validGenders := map[string]bool{"Male": true, "Female": true}
		if !validGenders[user.Gender] {
			log.Fatal("Invalid gender")
		}
		fmt.Print("Skill Level (Beginner/Intermediate/Master): ")
		fmt.Scanln(&user.SkillLevel)
		validSkills := map[string]bool{"Beginner": true, "Intermediate": true, "Master": true}
		if !validSkills[user.SkillLevel] {
			log.Fatal("Invalid skill level")
		}

		fmt.Print("What plan would you like to choose? (normal/semi-interactive/fully-interactive) :")
		fmt.Scanln(&user.Plan)
		validPlans := map[string]bool{"normal": true, "semi-interactive": true, "fully-interactive": true}
		if !validPlans[user.Plan] {
			log.Fatal("Invalid plan")
		}

		pk := insertUserData(db, user)

		fmt.Printf("User inserted successfully with id of %v!\n", pk)

		fmt.Print("Add another user? (y/n): ")
		var choice string
		fmt.Scanln(&choice)
		if choice != "y" && choice != "Y" {
			break
		}
	}

	fmt.Println("Exiting program.")
}
