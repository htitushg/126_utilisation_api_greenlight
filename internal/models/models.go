package models

import (
	"125_isbn_new/internal/validator"
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"time"
)

type Middleware func(handler http.HandlerFunc) http.HandlerFunc
type TokenStruct struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
	Status int       `json:"status"`
}
type Session struct {
	UserID         int       `json:"user_id"`
	ConnectionID   int       `json:"connection_id"`
	Pseudo         string    `json:"Pseudo"`
	Droits         string    `json:"Droits"`
	IpAddress      string    `json:"ip_address"`
	ExpirationTime time.Time `json:"expiration_time"`
}
type Logs []Log

type Log struct {
	Time       time.Time `json:"time"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	ReqId      int       `json:"req_id,omitempty"`
	User       Session   `json:"user,omitempty"`
	ClientIP   string    `json:"client_ip,omitempty"`
	ReqMethod  string    `json:"req_method,omitempty"`
	ReqURL     string    `json:"req_url,omitempty"`
	HttpStatus int       `json:"http_status,omitempty"`
	ErrOutput  string    `json:"output,omitempty"`
}

// Create a new userLoginForm struct.
type UserLoginForm struct {
	ID                  string `form:"id"`
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	Token               string `form:"token"`
	validator.Validator `form:"-"`
}

var (
	Ctx context.Context
	//Db  *sql.DB
)
var Tmpl = make(map[string]*template.Template)

type TConf_Mysql struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

// **** 2024-05-17-19h54
// Define a BaseModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// **** 2024-05-17-19h54 FIN

var Conf_Mysql TConf_Mysql

type data struct {
	Connect  bool
	Date     string
	Username string
	Email    string
	Droits   string
	Message  string
}

type Credentials struct {
	Username string
	Password string
	Droits   string
}

/*
	 type CreateUserMovie struct {
		Name                string    `json:"name"`
		User_id             int       `json:"id"`
		Token               string    `json:"token"`
		Expiry              time.Time `json:"expiry"`
		Scope               string    `json:"-"`
		validator.Validator `form:"-"`
	}
*/
type CreateUserMovie struct {
	User_id             int       `json:"id"`
	Name                string    `json:"name"`
	Email               string    `json:"email"`
	Token               string    `json:"token"`
	Expiry              time.Time `json:"expiry"`
	RToken              string    `json:"rtoken"`
	RExpiry             time.Time `json:"rexpiry"`
	Scope               string    `json:"-"`
	Activated           bool      `json:"activated"`
	CreatedAt           time.Time `json:"createdat"`
	Password            string    `json:"-"`
	validator.Validator `form:"-"`
}
type Movie struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"-"`
	Title     string `json:"title"`
	Year      int32  `json:"year,omitempty"`
	// Use the Runtime type instead of int32. Note that the omitempty directive will
	// still work on this: if the Runtime field has the underlying value 0, then it will
	// be considered empty and omitted -- and the MarshalJSON() method we just made
	// won't be called at all.
	Runtime Runtime  `json:"runtime,omitempty"` //`json:"-"` //Runtime  `json:"runtime,omitempty"`
	Genres  []string `json:"genres,omitempty"`
	Version int32    `json:"version"`
}

// ############################################################################################
type AuthenticationToken struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}
type RefreshToken struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}
type AuthenticationTokenApi struct {
	Id                  int                 `json:"id"`
	AuthenticationToken AuthenticationToken `json:"authentication_token"`
	RefreshToken        RefreshToken        `json:"refresh_token"`
}

// #############################################################################################
type AuthenticateUserApi struct {
	ID     int       `json:"id"`
	Scope  string    `json:"scope"`
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}
type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Droits    string    `json:"Droits"`
	HashedPwd string    `json:"hashedpwd"`
	Salt      string    `json:"salt"`
	Email     string    `json:"email"`
	Created   time.Time `json:"created"`
}

type TempUser struct {
	ConfirmID    string
	CreationTime time.Time
	User         User
}

type MailConfig struct {
	Email    string `json:"email_addr"`
	Auth     string `json:"email_auth"`
	Hostname string `json:"host"`
	Port     int    `json:"port"`
}

var Droits = []string{"normal", "admin"}

// structures de la base liées aux livres
type Livre struct {
	Livre_Id     int
	Idg          string
	Titre        string
	Isbn         string
	Thumbnail    string
	Editeur      Editeur
	Language     string
	Publish_date string
	Nb_pages     int
	CreatedAt    string
	Resume       string
	Description  string

	Auteurs      []Auteur
	LivreAuteurs []LivreAuteur
}
type LivreAuteur struct {
	Livre_Id  int
	Auteur_Id int
}
type Auteur struct {
	Auteur_Id   int
	Nom         string
	Description string
	CreatedAt   string
}
type Editeur struct {
	Editeur_Id  int
	Nom         string
	Description string
	CreatedAt   string
}
type ApiBook struct {
	Kind       string `json:"kind"`
	TotalItems int    `json:"totalItems"`
	Items      []struct {
		Kind       string `json:"kind"`
		ID         string `json:"id"`
		Etag       string `json:"etag"`
		SelfLink   string `json:"selfLink"`
		VolumeInfo struct {
			Title               string   `json:"title"`
			Subtitle            string   `json:"subtitle"`
			Authors             []string `json:"authors"`
			PublishedDate       string   `json:"publishedDate"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
			ReadingModes struct {
				Text  bool `json:"text"`
				Image bool `json:"image"`
			} `json:"readingModes"`
			PageCount           int    `json:"pageCount"`
			PrintType           string `json:"printType"`
			MaturityRating      string `json:"maturityRating"`
			AllowAnonLogging    bool   `json:"allowAnonLogging"`
			ContentVersion      string `json:"contentVersion"`
			PanelizationSummary struct {
				ContainsEpubBubbles  bool `json:"containsEpubBubbles"`
				ContainsImageBubbles bool `json:"containsImageBubbles"`
			} `json:"panelizationSummary"`
			Language            string `json:"language"`
			PreviewLink         string `json:"previewLink"`
			InfoLink            string `json:"infoLink"`
			CanonicalVolumeLink string `json:"canonicalVolumeLink"`
		} `json:"volumeInfo"`
		SaleInfo struct {
			Country     string `json:"country"`
			Saleability string `json:"saleability"`
			IsEbook     bool   `json:"isEbook"`
		} `json:"saleInfo"`
		AccessInfo struct {
			Country                string `json:"country"`
			Viewability            string `json:"viewability"`
			Embeddable             bool   `json:"embeddable"`
			PublicDomain           bool   `json:"publicDomain"`
			TextToSpeechPermission string `json:"textToSpeechPermission"`
			Epub                   struct {
				IsAvailable bool `json:"isAvailable"`
			} `json:"epub"`
			Pdf struct {
				IsAvailable bool `json:"isAvailable"`
			} `json:"pdf"`
			WebReaderLink       string `json:"webReaderLink"`
			AccessViewStatus    string `json:"accessViewStatus"`
			QuoteSharingAllowed bool   `json:"quoteSharingAllowed"`
		} `json:"accessInfo"`
	} `json:"items"`
}
