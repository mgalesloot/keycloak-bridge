package keycloak

import (
	"fmt"
)

// GroupInterface provides access to methods in the KeyCloak REST Client
type GroupInterface interface {
	GetGroups() []Group
	CreateGroup(groupName string, tenant string)
	RemoveGroup(groupID string)
	GetGroupID(groupName string) string
	GetGroupMembers(groupID string) []GroupMember
	RemoveMemberFromGroup(userid string, groupID string)
	AddMemberToGroup(userid string, groupID string)
}

// Group represents a group in KeyCloak
type Group struct {
	ID         string          `json:"id,omitempty"`
	Name       string          `json:"name"`
	Path       string          `json:"path"`
	Attributes GroupAttributes `json:"attributes"`
}

// GroupAttributes represents attributes of a KeyCloak group
type GroupAttributes struct {
	// Tenant attribute indicates the tenant a KeyCloak group belongs to
	Tenant []string `json:"tenant,omitempty"`
}

// GroupMember represents a member of a group in KeyCloak
type GroupMember struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
}

// GetGroups fetches all groups from KeyCloak
func (keyCloakClient Client) GetGroups() []Group {
	keyCloakGroups := make([]Group, 0)
	keyCloakClient.doGetRequest("groups?briefRepresentation=false", &keyCloakGroups)
	return keyCloakGroups
}

// CreateGroup creates a group in KeyCloak
func (keyCloakClient Client) CreateGroup(groupName string, tenant string) {
	requestBody := &Group{
		Name: groupName,
		Path: "/" + groupName,
		Attributes: GroupAttributes{
			Tenant: []string{tenant},
		},
	}
	keyCloakClient.doJSONRequest(MethodPost, "groups", requestBody, nil, 201)
}

// RemoveGroup removes a group from KeyCloak
func (keyCloakClient Client) RemoveGroup(groupID string) {
	keyCloakClient.doDeleteRequest(fmt.Sprintf("groups/%s", groupID))
}

// GetGroupID gets a group from KeyCloak and returns the ID
func (keyCloakClient Client) GetGroupID(groupName string) (groupID string) {
	groups := make([]Group, 0)
	keyCloakClient.doGetRequest(fmt.Sprintf("groups?search=%s", groupName), &groups)
	for _, group := range groups {
		if group.Name == groupName {
			return group.ID
		}
	}
	return
}

// GetGroupMembers gets members of a Group
func (keyCloakClient Client) GetGroupMembers(groupID string) []GroupMember {
	members := make([]GroupMember, 0)
	keyCloakClient.doGetRequest(fmt.Sprintf("groups/%s/members", groupID), &members)
	return members
}

// RemoveMemberFromGroup removes a member from a Group
func (keyCloakClient Client) RemoveMemberFromGroup(userid string, groupID string) {
	keyCloakClient.doDeleteRequest(fmt.Sprintf("users/%s/groups/%s", userid, groupID))
}

// AddMemberToGroup adds a member to a Group
func (keyCloakClient Client) AddMemberToGroup(userid string, groupID string) {
	keyCloakClient.doPutRequest(fmt.Sprintf("users/%s/groups/%s", userid, groupID))
}
