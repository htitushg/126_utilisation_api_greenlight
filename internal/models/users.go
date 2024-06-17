package models

import (
	//"125_isbn_new/internal/utils"
	"125_isbn_new/internal/utils"
	"125_isbn_new/internal/validator"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

// Define a new User struct. Notice how the field names and types align
// with the columns in the database "users" table?
/*
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}
*/
// AnonymousUser variable.
var AnonymousUser = &User{}

// Définissez une nouvelle structure UserModel qui encapsule
// un pool de connexions à la base de données.
type UsersModel struct {
	DB *sql.DB
}
type LivresModel struct {
	DB *sql.DB
}
type AuteursModel struct {
	DB *sql.DB
}
type EditeursModel struct {
	DB *sql.DB
}
type MoviesModel struct {
	DB *sql.DB
}

/*
Supprimez le champ de structure FieldErrors explicite et intégrez à la place
le validateur de structure. L'intégration signifie que notre
snippetCreateForm "hérite" de tous les champs et méthodes
de notre structure Validator (y compris le champ FieldErrors).
*/
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

// Create a new userSignupForm struct.
type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// Create a new userSignupForm struct.
type isbnForm struct {
	Isbn                string `form:"isbn"`
	validator.Validator `form:"-"`
}

/* // Create a new userLoginForm struct.
type userLoginForm struct {
	ID                  string `form:"id"`
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	Token               string `form:"token"`
	validator.Validator `form:"-"`
} */

// IsAnonymous Check if a User instance is the AnonymousUser.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// Nous utiliserons la méthode Insert pour ajouter un nouvel
// enregistrement à la table « utilisateurs ».
func (m *UsersModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
VALUES(?, ?, ?, UTC_TIMESTAMP())`
	// Use the Exec() method to insert the user details and hashed password
	// into the users table.
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// If this returns an error, we use the errors.As() function to check
		// whether the error has the type *mysql.MySQLError. If it does, the
		// error will be assigned to the mySQLError variable. We can then check
		// whether or not the error relates to our users_uc_email key by
		// checking if the error code equals 1062 and the contents of the error
		// message string. If it does, we return an ErrDuplicateEmail error.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// We'll use the Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UsersModel) Authenticate(email, password string) (CreateUserMovie, error) {
	// Retrieve the id and hashed password associated with the given email. If
	// no matching email exists we return the ErrInvalidCredentials error.
	var user CreateUserMovie
	var hashedPassword []byte
	stmt := "SELECT id, name, hashed_password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&user.User_id, &user.Name, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrInvalidCredentials
		} else {
			return user, err
		}
	}
	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return user, ErrInvalidCredentials
		} else {
			return user, err
		}
	}
	// Otherwise, the password is correct. Return the user ID.
	return user, nil
}

// We'll use the Exists method to check if a user exists with a specific ID.
func (m *UsersModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

var TempUsers []TempUser
var TempModifUsers []TempUser
var TempPassUserLosts []TempUser

// Delete TempUser
func deleteTempUser(temp TempUser) {
	for i, user := range TempUsers {
		if user == temp {
			TempUsers = append(TempUsers[:i], TempUsers[i+1:]...)
		}
	}
}

// Create User where it is confirmed
func (m *UsersModel) PushTempUser(id string) {
	log.Printf("TempUsers: %#v\n", TempUsers)
	log.Printf("id: %#v\n", id)
	for _, temp := range TempUsers {
		if temp.ConfirmID == id {
			temp.User.Id = m.GetIdNewUser()
			m.CreateUser(temp.User)
			deleteTempUser(temp)
		}
	}
}

// Supprime les enregistrement provisoires qui n'ont pas été validés après 12h
// func (m *UserModel) ManageTempUsers() {
func ManageTempUsers() {
	duration := utils.SetDailyTimer(2)
	for {
		for _, user := range TempUsers {
			if time.Since(user.CreationTime) > time.Hour*12 {
				deleteTempUser(user)
			}
		}
		for _, user := range TempModifUsers {
			if time.Since(user.CreationTime) > time.Hour*12 {
				deleteTempModifUsers(user)
			}
		}
		time.Sleep(duration)
		duration = time.Hour * 24
	}
}

// Delete TempModifUsers
func deleteTempModifUsers(temp TempUser) {
	for i, user := range TempModifUsers {
		if user == temp {
			TempModifUsers = append(TempModifUsers[:i], TempModifUsers[i+1:]...)
		}
	}
}

// Update User where it is confirmed
func (m *UsersModel) PushTempModifUser(id string) {
	//log.Println(GetCurrentFuncName())
	//log.Printf("PushTempModifUser id= %#v\n", id)
	for _, temp := range TempModifUsers {
		if temp.ConfirmID == id {
			m.UpdateUser(temp.User)
			deleteTempModifUsers(temp)
		}
	}
}

// Delete TempModifUsers
// func (m *UserModel) deleteTempPassUserLost(temp TempUser) {
func deleteTempPassUserLost(temp TempUser) {
	for i, user := range TempPassUserLosts {
		if user == temp {
			TempPassUserLosts = append(TempPassUserLosts[:i], TempPassUserLosts[i+1:]...)
		}
	}
}

// Update User where it is confirmed
func (m *UsersModel) PushTempPassUserLost(id string) {
	//log.Println(GetCurrentFuncName())
	//log.Printf("PushTempModifUser id= %#v\n", id)
	for _, temp := range TempPassUserLosts {
		if temp.ConfirmID == id {
			m.UpdateUser(temp.User)
			deleteTempPassUserLost(temp)
		}
	}
}
