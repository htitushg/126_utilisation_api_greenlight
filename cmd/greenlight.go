package main

import (
	"125_isbn_new/internal/models"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// #####################################################################
func (app application) InfoUserApi(name string, email string, password string) (nom string, err error) {

	data := map[string]string{
		"name":     name,
		"email":    email,
		"password": password,
	}

	// Encodage de la structure User en JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nom, err
	}

	// Création de l'URL de l'API
	url := "https://localhost:4000/v1/user/info"

	// Création de la requête HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nom, err
	}

	// Définir le type de contenu du body
	req.Header.Set("Content-Type", "application/json")
	// Ajouté le 17/06/2024 9h36
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Fin de Ajouté le 17/06/2024 9h36
	// Création d'un client HTTP
	//client := &http.Client{}

	// Envoi de la requête et récupération de la réponse
	resp, err := client.Do(req)
	if err != nil {
		return nom, err
	}
	defer resp.Body.Close()

	type jsonInfoUserApi struct {
		User struct {
			Name      string `json:"Name"`
			Email     string `json:"Email"`
			Password  string `json:"Password"`
			Activated bool   `json:"Activated"`
			CreatedAt string `json:"CreatedAt"`
			Code      int    `json:"Code"`
		} `json:"user"`
	}
	// Vérification du code de statut HTTP

	// Décodage du corps de la réponse JSON dans une nouvelle structure User
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nom, err
	}
	// Convertir le corps en chaîne de caractères
	bodyString := string(body)

	// Rechercher le mot "error"
	if strings.Contains(bodyString, name) {
		// Le nom figure bien dans la réponse
		var jsonInfoApi jsonInfoUserApi
		err = json.Unmarshal(body, &jsonInfoApi)
		if err != nil {
			return nom, nil
		}
		// Traitement du code de retour
		// Si 0 email non trouvé
		// Si -1 erreur serveur
		// Si 1 erreur Mot de passe incorrect
		// Si 2 Ok utilisateur reconnu
		// Si 4 le nom de l'utilisateur n'est pas correct
		switch jsonInfoApi.User.Code {
		case -1:
			err = models.ErrErreurServer
		case 0:
			err = models.ErrEmailNonTrouve
		case 1:
			err = models.ErrMdPIncorrect
		case 2:
			err = models.UserOk
		case 4:
			err = models.ErrNomIncorrect
		default:

		}
		return name, err
	} else {
		return nom, err
	}
}

// #####################################################################
// CreateUserApi : fonction de création d'un utlisateur de l'API Greenlight
func (app application) CreateUserApi(name string, email string, password string, ID int) (err error) {

	data := map[string]string{
		"name":     name,
		"email":    email,
		"password": password,
	}

	// Encodage de la structure User en JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Création de l'URL de l'API
	url := "https://localhost:4000/v1/users"

	// Création de la requête HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	// Définir le type de contenu du body
	req.Header.Set("Content-Type", "application/json")
	// Ajouté le 17/06/2024 9h36
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Fin de Ajouté le 17/06/2024 9h36
	// Création d'un client HTTP
	//client := &http.Client{}

	// Envoi de la requête et récupération de la réponse
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Vérification du code de statut HTTP
	/* if resp.StatusCode != http.StatusCreated {
		return cmovie, err
	} */

	// Décodage du corps de la réponse JSON dans une nouvelle structure User
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}
	// Convertir le corps en chaîne de caractères
	bodyString := string(body)

	// Rechercher le mot "error"
	if strings.Contains(bodyString, "error") {
		// Gérer l'erreur ici (par exemple, retourner une erreur personnalisée)
		fmt.Println("Error detected in response body")
		type ApiError struct {
			Error struct {
				Email string `json:"email"`
			} `json:"error"`
		}
		var apiError ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return err
		}
		return fmt.Errorf(apiError.Error.Email)
	} else {
		// tester si le body contient le mot error
		// Créer la structure User pour représenter les données à envoyer
		type CretateUser struct {
			User struct {
				ID        int       `json:"id"`
				CreatedAt time.Time `json:"created_at"`
				Name      string    `json:"name"`
				Email     string    `json:"email"`
				Activated bool      `json:"activated"`
			} `json:"user"`
		}
		var createUserApi CretateUser
		err = json.Unmarshal(body, &createUserApi)
		if err != nil {
			return err
		}
		// Renvoyer l'utilisateur créé
		fmt.Println(createUserApi)

		return nil
	}
}

// ##################################################################
// ActivateUserApi : fonction d'activation d'un utilisateur de l'Api Greenlight
func (app application) ActivateUserApi(tokenapi models.AuthenticateUserApi) (cmovie models.CreateUserMovie, err error) {
	type ActivateUser struct {
		Token string `json:"token"`
	}
	var activateUser ActivateUser
	activateUser.Token = tokenapi.Token
	// Encodage du champ Token de la structure tokenapi models.UserLoginForm en JSON
	jsonData, err := json.Marshal(activateUser)
	if err != nil {
		return cmovie, err
	}

	// Création de l'URL de l'API
	url := "https://localhost:4000/v1/users/activated"

	// Création de la requête HTTP PUT
	req, err := http.NewRequest("PUT", url, bytes.NewReader(jsonData))
	if err != nil {
		return cmovie, err
	}

	// Définir le type de contenu du body à "application/json"
	req.Header.Set("Content-Type", "application/json")
	// Ajouté le 17/06/2024 9h36
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Fin de Ajouté le 17/06/2024 9h36
	// Création d'un client HTTP
	//client := &http.Client{}

	// Envoi de la requête et récupération de la réponse
	resp, err := client.Do(req)
	if err != nil {
		return cmovie, err
	}
	defer resp.Body.Close()
	// Créer une nouvelle structure CreateUser pour récupérer la réponse de l'API
	type CreateUser struct {
		User struct {
			ID        int       `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			Name      string    `json:"name"`
			Email     string    `json:"email"`
			Activated bool      `json:"activated"`
		} `json:"user"`
	}
	var createUserApi CreateUser
	// Lecture du contenu de resp.Body dans un tableau d'octets
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return cmovie, err
	}
	// Décodage du corps de la réponse JSON dans createUserApi
	err = json.Unmarshal(bodyBytes, &createUserApi)
	if err != nil {
		return cmovie, err
	}
	//Vérification du code de statut HTTP
	if resp.StatusCode != http.StatusOK {
		return cmovie, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	// Renvoyer la structure
	fmt.Println(createUserApi)
	// Mettre mamovie, dans les champs movie
	cmovie.Name = createUserApi.User.Name
	cmovie.Email = createUserApi.User.Email
	cmovie.Activated = createUserApi.User.Activated
	cmovie.CreatedAt = createUserApi.User.CreatedAt
	return cmovie, err
}

// ##################################################################
// AuthenticateUserApi : fonction d'authentification de l'utilisateur de l'Api greenlight
func (app application) AuthenticateUserApi(email string, password string, ID int) (cmovie models.CreateUserMovie, err error) {
	type AuthenticateUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var authenticateUser AuthenticateUser
	authenticateUser.Email = email
	authenticateUser.Password = password
	// Encodage du champ Token de la structure tokenapi models.UserLoginForm en JSON
	jsonData, err := json.Marshal(authenticateUser)
	if err != nil {
		return cmovie, err
	}

	// Création de l'URL de l'API
	url := "https://localhost:4000/v1/tokens/authentication"

	// Création de la requête HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return cmovie, err
	}

	// Définir le type de contenu du body à "application/json"
	req.Header.Set("Content-Type", "application/json")
	// Ajouté le 17/06/2024 9h36
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Fin de Ajouté le 17/06/2024 9h36
	// Création d'un client HTTP
	//client := &http.Client{}

	// Envoi de la requête et récupération de la réponse
	resp, err := client.Do(req)
	if err != nil {
		return cmovie, err
	}
	defer resp.Body.Close()
	// Créer une nouvelle structure CreateUser pour récupérer la réponse de l'API
	type AuthenticateApi struct {
		Authentication_Token struct {
			Token  string    `json:"token"`
			Expiry time.Time `json:"expiry"`
		} `json:"authentication_token"`
	}
	var authenticateTokenApi AuthenticateApi
	// Lecture du contenu de resp.Body dans un tableau d'octets
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return cmovie, err
	}
	// Décodage du corps de la réponse JSON dans createUserApi
	err = json.Unmarshal(bodyBytes, &authenticateTokenApi)
	if err != nil {
		return cmovie, err
	}
	//Vérification du code de statut HTTP
	if resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 202 {
		// Renvoyer la structure
		fmt.Println(authenticateTokenApi)
		// Mettre mamovie, dans les champs movie
		cmovie.User_id = ID
		cmovie.Email = email
		cmovie.Token = authenticateTokenApi.Authentication_Token.Token
		cmovie.Expiry = authenticateTokenApi.Authentication_Token.Expiry
		return cmovie, nil
	} else {
		return cmovie, fmt.Errorf("status code: %d", resp.StatusCode)
	}

}

// ##################################################################
// ChangePwdUserApi : fonction de changement du Password de l'utilisateur de l'Api greenlight
func (app application) ChangePwdUserApi(email string, password string, newpassword string, ID int) (cmovie models.CreateUserMovie, ok bool) {
	type ConnectApi struct {
		Authentication_Token struct {
			Token  string    `json:"token"`
			Expiry time.Time `json:"expiry"`
		} `json:"authentication_token"`
	}
	var tokenApi ConnectApi
	// On construit l'URL
	url := "https://localhost:4000/v1/tokens/authentication"
	// Il s'agit d'une methode GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// On indique le token du demandeur dans le Header
	req.Header.Set("email", email)
	req.Header.Set("password", password)
	// Ajouté le 17/06/2024 9h36
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Fin de Ajouté le 17/06/2024 9h36
	//client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.StatusCode)
		return cmovie, false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return cmovie, false
	}
	// Décoder le JSON de la réponse
	err = json.Unmarshal(body, &tokenApi)
	if err != nil {
		log.Println(err)
		return cmovie, false
	}
	fmt.Println(tokenApi)
	// Mettre mamovie, dans les champs movie
	cmovie.User_id = ID
	return cmovie, true
}

// ##################################################################
func (app application) GetMovie(id string, token string) (movie models.Movie, ok bool) {
	// Acquisition du token

	type MyJsonName struct {
		Movie struct {
			Genres  []string `json:"genres"`
			ID      int64    `json:"id"`
			Runtime string   `json:"runtime"`
			Title   string   `json:"title"`
			Version int32    `json:"version"`
			Year    int32    `json:"year"`
		} `json:"movie"`
	}
	var mamovie MyJsonName
	// On construit l'URL
	url := fmt.Sprintf("https://localhost:4000/v1/movies/%s", id)
	// Il s'agit d'une methode GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// On indique le token du demandeur dans le Header
	entete := "Bearer " + token
	req.Header.Set("Authorization", entete)

	// Ajouté le 17/06/2024 9h36
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Fin de Ajouté le 17/06/2024 9h36
	//client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Stockage des en-têtes
	headers := resp.Header
	fmt.Printf("Headers = %v\n", headers)

	fmt.Printf("Status = %v , StatusCode = %v \n", resp.Status, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		// traiter les différents cas d'erreurs (404, ...)
		fmt.Println("Error:", resp.StatusCode)
		return movie, false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return movie, false
	}
	// Décoder le JSON de la réponse
	err = json.Unmarshal(body, &mamovie)
	if err != nil {
		log.Println(err)
		return movie, false
	}
	fmt.Println(mamovie)
	// Mettre mamovie, dans les champs movie
	movie.ID = mamovie.Movie.ID
	movie.CreatedAt = time.Now().Format("02/01/2006")
	runtimeArr := strings.Split(mamovie.Movie.Runtime, " ")
	runtime, err := strconv.ParseInt(runtimeArr[0], 10, 32)
	if err != nil {
		log.Println(err)
		return movie, false
	}
	movie.Runtime = models.Runtime(runtime)
	movie.Title = mamovie.Movie.Title
	movie.Version = mamovie.Movie.Version
	movie.Year = mamovie.Movie.Year
	movie.Genres = mamovie.Movie.Genres

	return movie, true
}
func (app application) GetMovies(token string) (movies []models.Movie, ok bool) {
	type ListeMovies struct {
		Metadata struct {
			CurrentPage  int `json:"current_page"`
			PageSize     int `json:"page_size"`
			FirstPage    int `json:"first_page"`
			LastPage     int `json:"last_page"`
			TotalRecords int `json:"total_records"`
		} `json:"metadata"`
		Movies []struct {
			ID      int64    `json:"id"`
			Title   string   `json:"title"`
			Year    int32    `json:"year"`
			Genres  []string `json:"genres"`
			Version int32    `json:"version"`
			Runtime string   `json:"runtime"`
		} `json:"movies"`
	}
	//var m models.UsersModel
	//user:=app.m.GetUserWithId)
	var mamovie models.Movie
	var mesmovies ListeMovies
	// On construit l'URL
	url := "https://localhost:4000/v1/movies"
	// Il s'agit d'une methode GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return movies, false
	}
	// On indique le token du demandeur dans le Header
	entete := "Bearer " + token
	req.Header.Set("Authorization", entete)

	// Ajouté le 17/06/2024 9h36
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Fin de Ajouté le 17/06/2024 9h36
	//client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return movies, false
	}
	defer resp.Body.Close()
	fmt.Printf("Status = %v , StatusCode = %v \n", resp.Status, resp.StatusCode)
	// Stockage des en-têtes
	headers := resp.Header
	fmt.Printf("Headers = %v\n", headers)

	if resp.StatusCode != http.StatusOK {
		// traiter les cas d'erreur

		fmt.Println("Error:", resp.StatusCode)
		return movies, false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return movies, false
	}
	// Décoder le JSON de la réponse
	err = json.Unmarshal(body, &mesmovies)
	if err != nil {
		log.Println(err)
		return movies, false
	}
	fmt.Println(mesmovies)
	// Mettre mesmovies.Movies, dans les champs mamovie
	// Puis ajouter mamovie dans movies
	for _, m := range mesmovies.Movies {
		mamovie.ID = m.ID
		mamovie.CreatedAt = time.Now().Format("02/01/2006")
		runtimeArr := strings.Split(m.Runtime, " ")
		runtime, err := strconv.ParseInt(runtimeArr[0], 10, 32)
		if err != nil {
			log.Println(err)
			return movies, false
		}
		mamovie.Runtime = models.Runtime(runtime)
		mamovie.Title = m.Title
		mamovie.Version = m.Version
		mamovie.Year = m.Year
		mamovie.Genres = m.Genres
		movies = append(movies, mamovie)
	}

	return movies, true
}
