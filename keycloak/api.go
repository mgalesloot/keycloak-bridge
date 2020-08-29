package keycloak

// API is an interface composed of all KeyCloak interfaces
type API interface {
	GroupInterface
	UserInterface
}
