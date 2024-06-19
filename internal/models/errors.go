package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// Add a new ErrInvalidCredentials error. We'll use this later if a user
	// tries to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// Add a new ErrDuplicateEmail error. We'll use this later if a user
	// tries to signup with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
	// Add a new ErrDuplicateIsbn error. We'll use this later if a user
	// tries to enter with an isbn code that's already in use.
	ErrDuplicateIsbn   = errors.New("models: duplicate isbn")
	ErrDuplicateIdg    = errors.New("models: duplicate idg")
	ErrDuplicateId     = errors.New("models: duplicate id")
	ErrDuplicateUserId = errors.New("models: duplicate user_id")
	ErrErreurServer    = errors.New("erreur du serveur api")  //-1
	ErrEmailNonTrouve  = errors.New("email non trouvé")       //0
	ErrMdPIncorrect    = errors.New("mot de passe incorrect") //1
	UserOk             = errors.New("utilisateur reconnu")    //2
	ErrNomIncorrect    = errors.New("nom incorrect")          //4
)
