package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
)

// Fonction qui vérifie que ni l'utilisateur, ni le courriel n'exitent dans la base
func (m *SnippetModel) UserOrEmailExist(name string, email string) bool {
	var err error

	var rows *sql.Rows
	resultat := false
	query := `SELECT * FROM users WHERE name = ? OR email = ?`

	// Préparation de la requête
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	// Exécution de la requête avec les valeurs des variables
	rows, err = stmt.Query(name, email)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	UnUser := User{}
	for rows.Next() {
		err = rows.Scan(&UnUser.Name)
		if err != nil {
			resultat = true
		} else {
			resultat = false
		}
	}
	return resultat
}

// Fonction qui vérifie si l'utilisateur exite dans la base
func (m *SnippetModel) UserExist(pseudo string) bool {
	var err error

	resultat := false
	query := `SELECT COUNT(*) FROM users WHERE pseudo = ?`

	// Préparation de la requête
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	// Exécution de la requête avec les valeurs des variables
	var count int
	err = stmt.QueryRow(pseudo).Scan(&count)
	if err != nil {
		panic(err)
	}

	// Test du résultat
	if count == 0 {
		// L'utilisateur n'existe pas
		println("Le pseudo 'admin' n'existe pas dans la table users.")
		resultat = false
	} else {
		// L'utilisateur existe
		println("Le pseudo 'admin' existe dans la table users.")
		resultat = true
	}
	return resultat
}

// Fonction qui vérifie si courriel exite dans la base
func (m *SnippetModel) EmailExist(email string) (user User, ok bool) {
	var err error

	resultat := false
	query := `SELECT * FROM users WHERE email = ?`

	// Préparation de la requête
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	// Exécution de la requête avec les valeurs des variables
	err = stmt.QueryRow(email).Scan(&user.Id, &user.Name, &user.Email, &user.Droits, &user.HashedPwd, &user.Salt)
	if err != nil {
		// L'utilisateur n'existe pas
		log.Printf("Le pseudo 'admin' n'existe pas dans la table users.\n")
		resultat = false
	} else {
		// L'utilisateur existe
		log.Printf("Le pseudo 'admin' existe dans la table users.\n")
		resultat = true
	}
	return user, resultat
}

// GetUser
// returns the models.User which models.User.Pseudo matches the `Pseudo` argument.
func (m *UsersModel) GetUser(name string) (user User, ok bool) {
	// DB simulation
	// Il faut ici se connecter à la base et vérifier si l'utilisateur
	// est bien enregistré
	var err error

	rows, err := m.DB.Query("SELECT id, name, email, hashed_password FROM users where name = ? ", name)
	if err != nil {
		panic(err)
	}
	ok = true
	defer rows.Close()

	i := 0
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.HashedPwd)
		if err != nil {
			ok = false
		}
		i++
	}
	return user, ok
}

// SelectUserwithId
// returns the models.User which models.User.Id matches the `Id` argument.
func (m *UsersModel) SelectUserwithId(id int) (user User, ok bool) {
	// DB simulation
	// Il faut ici se connecter à la base et vérifier si l'utilisateur
	// est bien enregistré
	var err error

	rows, err := m.DB.Query("SELECT id, name, email, hashed_password, created FROM users where id = ? ", id)
	if err != nil {
		panic(err)
	}
	ok = true
	defer rows.Close()

	i := 0
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.HashedPwd, &user.Created)
		if err != nil {
			ok = false
		}
		i++
	}
	return user, ok
}

// Fonction de création d'un utlisateur dans la base
func (m *UsersModel) CreateUser(user User) {
	var err error

	query := "INSERT INTO users (id, pseudo, email, droits, hashedPwd, 	salt )  VALUES 	(?, ?, ?, ?, ?, ?)"
	insertResult, err := m.DB.ExecContext(context.Background(), query, user.Id, user.Name, user.Email, user.Droits, user.HashedPwd, user.Salt)
	if err != nil {
		log.Fatalf("Impossible d'inserer le nouvel utilisateur: 	%s\n", err)
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: 	%s", err)
	}
	log.Printf("inserted id: %d", id)
}

// Fonction de mise à jour d'un utilisateur
func (m *UsersModel) UpdateUser(user User) {
	var err error

	//Prépare Update Db
	var stmt *sql.Stmt
	stmt, err = m.DB.Prepare("update users set pseudo=?,email=?,droits=? ,hashedPwd=?,salt=? where id=?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	// execute Update Db
	res, err := stmt.Exec(user.Name, user.Email, user.Droits, user.HashedPwd, user.Salt, user.Id)
	if err != nil {
		panic(err)
	}

	a, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	log.Printf("Result = %v\n", a)
	log.Printf("inserted id: %d", user.Id)
}
func (m *UsersModel) DeleteUserWithId(id int) {
	var err error

	query := "DELETE FROM users WHERE id=?"

	result, err := m.DB.Exec(query, id)
	if err != nil {
		log.Fatalf("Impossible d'effacer l'utilisateur: 	%s\n", err)
	}
	log.Printf("Résultat de l'effacement: %v\n", result)
}

// GetIdNewUser
// returns first unused id in database users
func (m *UsersModel) GetIdNewUser() (id int) {
	// DB simulation
	// Il faut ici se connecter à la base et vérifier si l'utilisateur
	// est bien enregistré
	var err error
	var users []User
	var user User

	rows, err := m.DB.Query("SELECT id FROM users ")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&user.Id)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
		i++
	}
	var idFound bool
	for id = 1; !idFound; id++ {
		idFound = true
		for _, user := range users {
			if user.Id == id {
				idFound = false
			}
		}
	}
	id--
	return id
}
func (m *UsersModel) GetUsers() (users []User, ok bool) {
	// DB simulation
	// Il faut ici se connecter à la base et vérifier si l'utilisateur
	// est bien enregistré
	var err error

	rows, err := m.DB.Query("SELECT * FROM users ")
	if err != nil {
		panic(err)
	}
	ok = true
	defer rows.Close()
	var user User
	i := 0
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Droits, &user.HashedPwd, &user.Salt)
		if err != nil {
			ok = false
		}
		users = append(users, user)
		i++
	}
	return users, ok
}

// GetIdNewEditeur
// returns first unused id in database users
func (m *EditeursModel) GetIdNewEditeur() (id int) {
	// DB simulation
	// Il faut ici se connecter à la base et vérifier si l'utilisateur
	// est bien enregistré
	var err error
	var editeurs []Editeur
	var editeur Editeur

	rows, err := m.DB.Query("SELECT editeur_id FROM editeurs ")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&editeur.Editeur_Id)
		if err != nil {
			panic(err)
		}
		editeurs = append(editeurs, editeur)
		i++
	}
	var idFound bool
	for id = 1; !idFound; id++ {
		idFound = true
		for _, editeur := range editeurs {
			if editeur.Editeur_Id == id {
				idFound = false
			}
		}
	}
	id--
	return id
}

// Renvoie la liste des editeurs
func (m *EditeursModel) GetEditeurs() (editeurs []Editeur) {
	var err error
	var editeur Editeur

	rows, err := m.DB.Query("SELECT * FROM editeurs ")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&editeur.Editeur_Id, &editeur.Nom, &editeur.CreatedAt, &editeur.Description)
		if err != nil {
			panic(err)
		}
		editeurs = append(editeurs, editeur)
		i++
	}

	return editeurs
}
func (m *AuteursModel) GetAuteurs() (auteurs []Auteur) {
	var auteur Auteur

	rows, err := m.DB.Query("SELECT * FROM auteurs ")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&auteur.Auteur_Id, &auteur.Nom, &auteur.CreatedAt, &auteur.Description)
		if err != nil {
			panic(err)
		}
		auteurs = append(auteurs, auteur)
		i++
	}

	return auteurs
}
func (m *AuteursModel) GetIdNewAuteur() (id int) {
	// DB simulation
	// Il faut ici se connecter à la base et vérifier si l'utilisateur
	// est bien enregistré
	var err error
	var auteurs []Auteur
	var auteur Auteur

	rows, err := m.DB.Query("SELECT auteur_id FROM auteurs ")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&auteur.Auteur_Id)
		if err != nil {
			panic(err)
		}
		auteurs = append(auteurs, auteur)
		i++
	}
	var idFound bool
	for id = 1; !idFound; id++ {
		idFound = true
		for _, auteur := range auteurs {
			if auteur.Auteur_Id == id {
				idFound = false
			}
		}
	}
	id--
	return id
}
func (m *LivresModel) GetIdNewLivre() (id int) {
	// DB simulation
	// Il faut ici se connecter à la base et vérifier si l'utilisateur
	// est bien enregistré
	var err error
	var livres []Livre
	var livre Livre

	rows, err := m.DB.Query("SELECT livre_id FROM livres ")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&livre.Livre_Id)
		if err != nil {
			panic(err)
		}
		livres = append(livres, livre)
		i++
	}
	var idFound bool
	for id = 1; !idFound; id++ {
		idFound = true
		for _, livre := range livres {
			if livre.Livre_Id == id {
				idFound = false
			}
		}
	}
	id--
	return id
}
func (m *EditeursModel) CreateEditeur(editeur Editeur) int {
	var err error

	DatedeCreation := []byte(time.Now().Format("2006-01-02"))
	id := m.GetIdNewEditeur()
	fmt.Printf("id = %d, nom = %v, description = %v\n", id, editeur.Nom, editeur.Description)
	query := "INSERT INTO editeurs (editeur_id, nom, description, created_at)  VALUES 	(?, ?, ?, ?)"
	insertResult, err := m.DB.ExecContext(context.Background(), query, id, editeur.Nom, editeur.Description, DatedeCreation)
	if err != nil {
		log.Fatalf("Impossible d'inserer le nouvel editeur: 	%s\n", err)
	}
	insertid, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: 	%s", err)
	}
	log.Printf("inserted id: %d", insertid)
	return id
}
func (m *AuteursModel) CreateAuteur(auteur Auteur) int {
	var err error
	DatedeCreation := []byte(time.Now().Format("2006-01-02"))
	id := m.GetIdNewAuteur()
	query := "INSERT INTO auteurs (auteur_id, nom, description, created_at)  VALUES 	(?, ?, ?, ?)"
	insertResult, err := m.DB.ExecContext(context.Background(), query, id, auteur.Nom, auteur.Description, DatedeCreation)
	if err != nil {
		log.Fatalf("Impossible d'inserer le nouvel auteur: 	%s\n", err)
	}

	insertid, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: 	%s", err)
	}
	log.Printf("inserted id: %d", insertid)
	return id
}
func (m *AuteursModel) LivreAuteurExist(livre_Id, auteur_Id int) bool {

	query := `SELECT * FROM livre_auteur WHERE id_livre = ? AND id_auteur = ?`
	// Préparation de la requête
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	// Exécution de la requête avec les valeurs des variables
	rows, err := stmt.Query(livre_Id, auteur_Id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Vérification du nombre de lignes renvoyées
	count := 0
	for rows.Next() {

		count++
	}

	// Affichage d'un message en fonction du nombre de lignes renvoyées
	if count == 0 {
		return false
	} else {
		// Boucle pour parcourir les lignes
		var livaut LivreAuteur
		rows, err = stmt.Query(livre_Id, auteur_Id)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			_ = rows.Scan(&livaut.Livre_Id, &livaut.Auteur_Id)
		}
		return true
	}
}

func (m *AuteursModel) CreateLivreAuteur(livre_id, auteur_id int) bool {
	var err error

	query := "INSERT INTO livre_auteur (id_livre, id_auteur)  VALUES 	(?, ?)"
	insertResult, err := m.DB.ExecContext(context.Background(), query, livre_id, auteur_id)
	if err != nil {
		log.Fatalf("Impossible d'inserer le nouvel editeur: 	%s\n", err)
	}

	insertid, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: 	%s", err)
	}
	log.Printf("inserted id: %d", insertid)
	return true
}
func (m *LivresModel) LivreExist(isbn string) int {
	var err error

	query := `SELECT livre_id, isbn FROM livres WHERE isbn=?`
	// Préparation de la requête
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	// Exécution de la requête avec les valeurs des variables
	rows, err := stmt.Query(isbn)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Vérification du nombre de lignes renvoyées
	count := 0
	var livre Livre
	for rows.Next() {
		count++
		_ = rows.Scan(&livre.Livre_Id, &livre.Isbn)
	}
	return livre.Livre_Id
}

// Créer un nouveau livre dans la base
func (m *LivresModel) CreateLivre(new_livre Livre) error {
	DatedeCreation := []byte(time.Now().Format("2006-01-02"))
	// Vérifier si les champs uniques le sont bien
	// avant de lancer l'insertion du livre
	query := `SELECT livre_id, idg, isbn FROM livres WHERE isbn=? || idg = ? || livre_id = ?`
	formatedQuery := fmt.Sprintf(query)
	stmt, err := m.DB.Prepare(formatedQuery)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	// Exécution de la requête avec les valeurs des variables
	rows, err := stmt.Query(new_livre.Isbn, new_livre.Idg, new_livre.Livre_Id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Vérification du nombre de lignes renvoyées
	count := 0
	var livre Livre
	for rows.Next() {
		count++
		_ = rows.Scan(&livre.Livre_Id, &livre.Isbn)
	}
	if count == 0 {

		//id := GetIdNewLivre()
		// il faut créer si nécessaire l'éditeur dans la table editeurs
		// car en raison des liasons il n'est pas possible de créer un livre si l'éditeur n'existe pas
		// Il faut créer le livre dans la table livres
		// Il faut créer si nécessaire les auteurs dans la table auteurs

		query := `INSERT INTO livres (
		livre_id, idg, titre, isbn, thumbnail,
		editeur_id, language, publish_date,
		nb_pages, created_at, resume, description )  
		VALUES 	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		formatedQuery := fmt.Sprintf(query)
		insertResult, err := m.DB.ExecContext(context.Background(), formatedQuery,
			new_livre.Livre_Id, new_livre.Idg, new_livre.Titre, new_livre.Isbn, new_livre.Thumbnail,
			new_livre.Editeur.Editeur_Id, new_livre.Language, new_livre.Publish_date,
			new_livre.Nb_pages, DatedeCreation, new_livre.Resume, new_livre.Description)

		//query, new_livre.Livre_Id, new_livre.Idg, new_livre.Titre, new_livre.Isbn, new_livre.Thumbnail, new_livre.Editeur.Editeur_Id, new_livre.Language, new_livre.Publish_date, new_livre.Nb_pages, new_livre.Publish_date, new_livre.Nb_pages, DatedeCreation, new_livre.Resume, new_livre.Description)
		if err != nil {
			// If this returns an error, we use the errors.As() function to check
			// whether the error has the type *mysql.MySQLError. If it does, the
			// error will be assigned to the mySQLError variable. We can then check
			// whether or not the error relates to our users_uc_email key by
			// checking if the error code equals 1062 and the contents of the error
			// message string. If it does, we return an ErrDuplicateEmail error.
			var mySQLError *mysql.MySQLError
			if errors.As(err, &mySQLError) {
				if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "isbn") {
					return ErrDuplicateIsbn
				}
				if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "idg") {
					return ErrDuplicateIdg
				}
			} else {
				return err
			}
		}
		insertid, err := insertResult.LastInsertId()
		if err != nil {
			log.Fatalf("impossible to retrieve last inserted id: 	%s", err)
		}
		log.Printf("inserted id: %d", insertid)
		return err
	} else {
		return ErrDuplicateIsbn
	}
}
func (m *LivresModel) GetLivre(isbn string) (livre Livre) {

	rows, err := m.DB.Query("SELECT * FROM livres WHERE livres.isbn=?", isbn)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&livre.Livre_Id, &livre.Idg, &livre.Titre, &livre.Isbn, &livre.Thumbnail, &livre.Editeur.Editeur_Id, &livre.Language, &livre.Publish_date, &livre.Nb_pages, &livre.CreatedAt, &livre.Resume, &livre.Description)
		i++
	}

	return livre
}
func (m *LivresModel) GetLivres() (livres []Livre) {
	var livre Livre
	rows, err := m.DB.Query("SELECT * FROM livres ")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&livre.Livre_Id, &livre.Idg, &livre.Titre, &livre.Isbn, &livre.Thumbnail, &livre.Editeur.Editeur_Id, &livre.Language, &livre.Publish_date, &livre.Nb_pages, &livre.CreatedAt, &livre.Resume, &livre.Description)
		if err != nil {
			panic(err)
		}
		livres = append(livres, livre)
		i++
	}

	return livres
}

// Cette fonction retourne les livres et les auteurs associés
func (m *AuteursModel) GetLivreetEditeurAuteurs(isbn string) (livre Livre) {
	// requête sql : recherche le livre qui a le code 'isbn'
	rows, err := m.DB.Query("SELECT * FROM livres WHERE livres.isbn=?", isbn)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&livre.Livre_Id, &livre.Idg, &livre.Titre, &livre.Isbn, &livre.Thumbnail, &livre.Editeur.Editeur_Id, &livre.Language, &livre.Publish_date, &livre.Nb_pages, &livre.CreatedAt, &livre.Resume, &livre.Description)
		i++
	}
	log.Printf("desLivresA.Livre_Id : %v, desLivresA.Titre: %v, desLivresA.Editeur_Id: %v, desLivresA.Isbn: %v\n",
		livre.Livre_Id, livre.Titre, livre.Editeur.Editeur_Id, livre.Isbn)
	// Acquisition des auteurs du livre
	sql_auteurs := `SELECT auteurs.auteur_id, auteurs.nom, auteurs.created_at, auteurs.description 	FROM livres
	JOIN livre_auteur ON livres.livre_id = livre_auteur.id_livre
	JOIN auteurs ON livre_auteur.id_auteur = auteurs.auteur_id
	WHERE livres.livre_id=?`
	rowsA, err := m.DB.Query(sql_auteurs, livre.Livre_Id)
	if err != nil {
		panic(err)
	}
	for rowsA.Next() {
		var auteur Auteur
		err = rowsA.Scan(&auteur.Auteur_Id, &auteur.Nom, &auteur.CreatedAt, &auteur.Description)
		livre.Auteurs = append(livre.Auteurs, auteur)
		if err != nil {
			panic(err)
		}
	}
	defer rowsA.Close()
	// Acquisition des éditeurs du livre
	sql_editeur := `SELECT editeurs.editeur_id, editeurs.nom FROM livres 
	JOIN editeurs ON livres.editeur_id = editeurs.editeur_id
	WHERE livres.editeur_id=?`
	rowsA, err = m.DB.Query(sql_editeur, livre.Editeur.Editeur_Id)
	if err != nil {
		panic(err)
	}
	for rowsA.Next() {
		var editeur Editeur
		err = rowsA.Scan(&editeur.Editeur_Id, &editeur.Nom)
		livre.Editeur = editeur
		if err != nil {
			panic(err)
		}
	}
	defer rowsA.Close()
	return livre
}

// Cette fonction retourne les livres et les auteurs associés
func (m *AuteursModel) GetLivresetEditeursAuteurs() (desLivres []Livre) {

	Unlivre := Livre{}
	// requête sql
	sql_livres := `SELECT * FROM livres`
	rowsL, err := m.DB.Query(sql_livres)
	if err != nil {
		panic(err)
	}
	i := 0
	//j := 0
	for rowsL.Next() {
		err = rowsL.Scan(&Unlivre.Livre_Id, &Unlivre.Idg, &Unlivre.Titre, &Unlivre.Isbn, &Unlivre.Thumbnail, &Unlivre.Editeur.Editeur_Id, &Unlivre.Language, &Unlivre.Publish_date, &Unlivre.Nb_pages, &Unlivre.CreatedAt, &Unlivre.Resume, &Unlivre.Description)
		if err != nil {
			panic(err)
		}
		desLivres = append(desLivres, Unlivre)
		//log.Printf("desLivresA.Livre_Id : %v, desLivresA.Titre: %v, desLivresA.Editeur_Id: %v, desLivresA.Isbn: %v\n",
		//	desLivres[i].Livre_Id, desLivres[i].Titre, desLivres[i].Editeur.Editeur_Id, desLivres[i].Isbn)
		// Acquisition des auteurs du livre
		sql_auteurs := `SELECT auteurs.auteur_id, auteurs.nom, auteurs.created_at, auteurs.description 	FROM livres
		JOIN livre_auteur ON livres.livre_id = livre_auteur.id_livre
		JOIN auteurs ON livre_auteur.id_auteur = auteurs.auteur_id
		WHERE livres.livre_id=?`
		rowsA, err := m.DB.Query(sql_auteurs, desLivres[i].Livre_Id)
		if err != nil {
			panic(err)
		}
		for rowsA.Next() {
			var auteur Auteur
			err = rowsA.Scan(&auteur.Auteur_Id, &auteur.Nom, &auteur.CreatedAt, &auteur.Description)
			desLivres[i].Auteurs = append(desLivres[i].Auteurs, auteur)
			if err != nil {
				panic(err)
			}
		}
		defer rowsA.Close()
		// Acquisition des éditeurs du livre
		sql_editeur := `SELECT editeurs.editeur_id, editeurs.nom FROM livres 
		JOIN editeurs ON livres.editeur_id = editeurs.editeur_id
		WHERE livres.editeur_id=?`
		rowsA, err = m.DB.Query(sql_editeur, Unlivre.Editeur.Editeur_Id)
		if err != nil {
			panic(err)
		}
		for rowsA.Next() {
			var editeur Editeur
			err = rowsA.Scan(&editeur.Editeur_Id, &editeur.Nom)
			desLivres[i].Editeur = editeur
			if err != nil {
				panic(err)
			}
		}
		i++
		defer rowsA.Close()

	}
	defer rowsL.Close()
	return desLivres
}
func (m *EditeursModel) EditeurExist(nomediteur string) (editeur Editeur, ok bool) {

	query := `SELECT * FROM editeurs WHERE nom=?`
	// Préparation de la requête
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	// Exécution de la requête avec les valeurs des variables
	rows, err := stmt.Query(nomediteur)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Vérification du nombre de lignes renvoyées
	count := 0
	for rows.Next() {

		count++
	}

	// Affichage d'un message en fonction du nombre de lignes renvoyées
	if count == 0 {
		return editeur, false
	} else {
		// Boucle pour parcourir les lignes
		rows, err = stmt.Query(nomediteur)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			_ = rows.Scan(&editeur.Editeur_Id, &editeur.Nom, &editeur.CreatedAt, &editeur.Description)
		}
		return editeur, true
	}
}
func (m *AuteursModel) AuteurExist(nomauteur string) (auteur Auteur, ok bool) {
	query := `SELECT * FROM auteurs WHERE nom=?`
	// Préparation de la requête
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	// Exécution de la requête avec les valeurs des variables
	rows, err := stmt.Query(nomauteur)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Vérification du nombre de lignes renvoyées
	count := 0
	for rows.Next() {
		err = rows.Scan(&auteur.Auteur_Id, &auteur.Nom, &auteur.CreatedAt, &auteur.Description)
		count++
	}

	// Affichage d'un message en fonction du nombre de lignes renvoyées
	if count == 0 {
		return auteur, false
	} //else {
	// Boucle pour parcourir les lignes
	//rows, err = stmt.Query(nomauteur)
	//if err != nil {
	//	panic(err)
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	auteur = rows.Scan(&auteur.Auteur_Id, &auteur.Nom, &auteur.CreatedAt, &auteur.Description)
	//}
	return auteur, true
	//}
}

// ####################################################################################
func (m *MoviesModel) LireJetonDansBase(user_id int) (token AuthenticateUserApi, err error) {

	query := `SELECT id, token , expiry FROM tokens WHERE id=?`

	// Préparation de la requête
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return token, err
	}
	defer stmt.Close()
	// Exécution de la requête avec les valeurs des variables
	rows, err := stmt.Query(user_id)
	if err != nil {
		return token, err
	}
	defer rows.Close()
	// Vérification du nombre de lignes renvoyées
	count := 0
	for rows.Next() {
		err = rows.Scan(&token.ID, &token.Token, &token.Expiry)
		if err != nil {
			fmt.Printf("Erreur : %v\n", err)
			return token, err
		}
		count++
	}
	if count > 1 {
		err = nil
	}
	return token, err
}

// ####################################################################################
func (m *MoviesModel) EcrireJetonDansBase(token AuthenticateUserApi) (err error) {

	// vérification de l'existence d'un jeton pour cet utilisateur
	var token2 AuthenticateUserApi
	token2, err = m.LireJetonDansBase(token.ID)
	if err != nil {
		return err
	}
	if token2.ID > 0 {
		query := `UPDATE tokens SET token=?, id= ? , expiry= ? WHERE id = ?`
		res, err := m.DB.Exec(query, token.Token, token.ID, token.Expiry, token.ID)
		if err != nil {
			log.Fatal(err)
		}

		count, err := res.RowsAffected()

		if err != nil {
			log.Fatal(err)
		}

		log.Println(count)
		return err

	} else {
		query := `INSERT INTO tokens (token, id, expiry ) VALUES (?, ?, ?)`
		//formatedQuery := fmt.Sprintf(query)
		insertResult, err := m.DB.ExecContext(context.Background(), query,
			token.Token, token.ID, token.Expiry)
		if err != nil {
			// If this returns an error, we use the errors.As() function to check
			// whether the error has the type *mysql.MySQLError. If it does, the
			// error will be assigned to the mySQLError variable. We can then check
			// whether or not the error relates to our users_uc_email key by
			// checking if the error code equals 1062 and the contents of the error
			// message string. If it does, we return an ErrDuplicateEmail error.
			var mySQLError *mysql.MySQLError
			if errors.As(err, &mySQLError) {
				if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "token.ID") {
					return ErrDuplicateUserId
				}
			} else {
				return err
			}
		}
		insertid, err := insertResult.LastInsertId()
		if err != nil {
			log.Fatalf("impossible to retrieve last inserted id: 	%s", err)
		}
		log.Printf("inserted id: %d", insertid)
		return err
	}
}

// ####################################################################################
func (m *MoviesModel) MovieExist(id int64) Movie {
	var err error

	query := `SELECT * FROM movies WHERE id=?`
	// Préparation de la requête
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	// Exécution de la requête avec les valeurs des variables
	rows, err := stmt.Query(id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Vérification du nombre de lignes renvoyées
	count := 0
	var movie Movie
	for rows.Next() {
		count++
		_ = rows.Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, &movie.Genres, &movie.Version)
	}
	return movie
}

func (m *MoviesModel) GetMovie(id int64) (Movie, bool) {
	var err error

	query := `SELECT * FROM movies WHERE id=?`
	// Préparation de la requête
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	// Exécution de la requête avec les valeurs des variables
	rows, err := stmt.Query(id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Vérification du nombre de lignes renvoyées
	count := 0
	var movie Movie
	for rows.Next() {
		count++
		_ = rows.Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, &movie.Genres, &movie.Version)
	}
	if count < 1 {
		return movie, false
	}
	return movie, true
}
func (m *MoviesModel) CreateMovie(movie Movie) int64 {
	var err error
	DatedeCreation := []byte(time.Now().Format("2006-01-02"))
	id := m.GetIdNewMovie()
	query := "INSERT INTO movies (id, createdat, title, year, runtime, genres, version)  VALUES 	(?, ?, ?, ?, ?, ?, ?)"
	insertResult, err := m.DB.ExecContext(context.Background(), query, &movie.ID, DatedeCreation, &movie.Title, &movie.Year, &movie.Runtime, &movie.Genres, &movie.Version)
	if err != nil {
		log.Fatalf("Impossible d'inserer le nouveau film: 	%s\n", err)
	}

	insertid, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: 	%s", err)
	}
	log.Printf("inserted id: %d", insertid)
	return id
}
func (m *MoviesModel) GetIdNewMovie() (id int64) {
	// DB simulation
	// Il faut ici se connecter à la base et vérifier si l'utilisateur
	// est bien enregistré
	var err error
	var mesmovies []Movie
	var mamovie Movie

	rows, err := m.DB.Query("SELECT id FROM movies ")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var i int64 = 0
	for rows.Next() {
		err = rows.Scan(&mamovie.ID)
		if err != nil {
			panic(err)
		}
		mesmovies = append(mesmovies, mamovie)
		i++
	}
	var idFound bool
	for id = 1; !idFound; id++ {
		idFound = true
		for _, mamovie := range mesmovies {
			if mamovie.ID == id {
				idFound = false
			}
		}
	}
	id--
	return id
}
