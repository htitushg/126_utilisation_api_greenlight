package main

import (
	"125_isbn_new/internal/models"
	"context"
	"net/http"
)

const isAuthenticatedContextKey = contextKey("isAuthenticated")

// Define a custom contextKey type, with the underlying type string.
// Définissez un type contextKey personnalisé, avec la chaîne de type sous-jacente.
type contextKey string

// Convert the string "user" to a contextKey type and assign it to the userContextKey
// constant. We'll use this constant as the key for getting and setting user information
// in the request context.
// Convertit la chaîne "user" en un type contextKey et l'attribue à la constante
// userContextKey. Nous utiliserons cette constante comme clé pour obtenir et définir les informations utilisateur
// dans le contexte de la requête.
const userContextKey = contextKey("user")

// The contextSetUser() method returns a new copy of the request with the provided
// User struct added to the context. Note that we use our userContextKey constant as the
// key.
// La méthode contextSetUser() renvoie une nouvelle copie de la requête avec la
// Structure utilisateur ajoutée au contexte.
// Notez que nous utilisons notre constante userContextKey comme clé.
func (app *application) contextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// The contextGetUser() retrieves the User struct from the request context. The only
// time that we'll use this helper is when we logically expect there to be User struct
// value in the context, and if it doesn't exist it will firmly be an 'unexpected' error.
// As we discussed earlier in the book, it's OK to panic in those circumstances.
// ContextGetUser() récupère la structure User à partir du contexte de la requête.
// Le seul moment où nous utiliserons cet assistant
// est le moment où nous nous attendons logiquement à ce qu'il y ait un identifiant utilisateur
// dans le contexte, et s'il n'existe pas, ce sera clairement une erreur "inattendue".
// Comme nous l'avons expliqué plus tôt dans le livre, il est normal de paniquer dans ces circonstances.
func (app *application) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
