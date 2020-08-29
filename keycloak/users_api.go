package keycloak

import (
	"fmt"
	"log"
)

// UserInterface provides access to methods in the KeyCloak REST Client
type UserInterface interface {
	GetUserID(username string) string
	GetUser(username string) User
}

// User represents a user in KeyCloak
type User struct {
	ID        string `json:"id"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// GetUserID gets the ID of a user
func (keyCloakClient Client) GetUserID(username string) string {
	return keyCloakClient.GetUser(username).ID
}

// GetUser gets a User from KeyCloak
func (keyCloakClient Client) GetUser(username string) (user User) {
	var keyCloakUsers []User
	keyCloakClient.doGetRequest(fmt.Sprintf("users?username=%s", username), &keyCloakUsers)

	for _, user := range keyCloakUsers {
		if user.UserName == username {
			return user
		}
	}

	log.Fatalf("User %q was not found in KeyCloak\n", username)
	return
}
