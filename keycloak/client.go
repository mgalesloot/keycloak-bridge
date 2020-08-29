package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"keycloak-bridge/config"
)

// Http Methods
const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)

// Content types
const (
	jsonContentType = "application/json"
)

// Client struct with api methods and accesstoken
type Client struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	RenewTokenAt time.Time
	FQDN         string
	Realm        string
	Namespace    string
}

// NewKeyCloakClient initalizes a new Client struct and connects to KeyCloak
func NewKeyCloakClient(fileName string) *Client {
	var keyCloakClient = &Client{}

	log.Printf("Init Keycloak Client")
	platform := config.LoadPlatformSettings(fileName)

	keyCloakClient.FQDN = platform.Config.KeyCloak.FQDN
	keyCloakClient.Realm = platform.Config.KeyCloak.Realm
	keyCloakClient.Namespace = platform.Config.KeyCloak.Namespace

	secretData := getSecretData("keycloak-http", keyCloakClient.Namespace)

	formData := url.Values{}
	formData.Set("grant_type", "password")
	formData.Set("client_id", "admin-cli")
	formData.Set("username", "keycloak")
	formData.Set("password", string(secretData["password"]))

	res, err := http.PostForm(fmt.Sprintf("https://%s/auth/realms/master/protocol/openid-connect/token", keyCloakClient.FQDN), formData)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		panic(err)
	}

	json.Unmarshal(b, keyCloakClient)
	setExpiresAt(keyCloakClient)

	return keyCloakClient
}

// setExpiresAt sets the time when the access token needs renewal
// RenewTokenAt is set 5 seconds before actual expiration time
func setExpiresAt(keyCloakClient *Client) {
	const secondsBefore = 5
	keyCloakClient.RenewTokenAt = time.Now().Add(time.Duration(keyCloakClient.ExpiresIn-secondsBefore) * time.Second)
}

func (keyCloakClient Client) getNewToken() {

	if keyCloakClient.RefreshToken == "" {
		panic("Refresh Token should not be empty")
	}

	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("client_id", "admin-cli")
	formData.Set("refresh_token", keyCloakClient.RefreshToken)

	res, err := http.PostForm(fmt.Sprintf("https://%s/auth/realms/master/protocol/openid-connect/token", keyCloakClient.FQDN), formData)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		panic(err)
	}

	json.Unmarshal(b, &keyCloakClient)
	setExpiresAt(&keyCloakClient)
}

func (keyCloakClient Client) doGetRequest(requestPath string, response interface{}) {
	keyCloakClient.doJSONRequest(MethodGet, requestPath, nil, response, 200)
}

func (keyCloakClient Client) doDeleteRequest(requestPath string) {
	keyCloakClient.doJSONRequest(MethodDelete, requestPath, nil, nil, 204)
}

func (keyCloakClient Client) doPutRequest(requestPath string) {
	keyCloakClient.doJSONRequest(MethodPut, requestPath, nil, nil, 204)
}

func (keyCloakClient Client) doJSONRequest(method string, requestPath string, body interface{}, response interface{}, allowedStatuscode int) {
	if keyCloakClient.AccessToken == "" {
		panic("Access Token should not be empty")
	}

	if time.Now().After(keyCloakClient.RenewTokenAt) {
		keyCloakClient.getNewToken()
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)

	url := fmt.Sprintf("https://%s/auth/admin/realms/%s/%s", keyCloakClient.FQDN, keyCloakClient.Realm, requestPath)
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+keyCloakClient.AccessToken)
	req.Header.Set("Content-Type", jsonContentType)

	var client = &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != allowedStatuscode {
		panic(fmt.Sprintf("URL: %s, Method %s, Body: %s, Error status code: %d", url, method, body, res.StatusCode))
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if len(responseData) > 0 {
		err = json.Unmarshal(responseData, response)
		if err != nil {
			panic(fmt.Sprintf("URL: %s, Method %s, Body: %s, Response: %s", url, method, body, responseData))
		}
	}
}
