package group

import (
	"testing"

	"keycloak-bridge/config"
	"keycloak-bridge/keycloak"
)

// mockKeyCloakClient holds some data to check the calls being done by the GroupReconciler.
type mockKeyCloakClient struct {
	groupCreated  string
	groupRemoved  string
	memberAdded   string
	memberRemoved string
}

var mock = mockKeyCloakClient{}

func (m *mockKeyCloakClient) CreateGroup(groupName string, tenant string) {
	m.groupCreated = groupName
}

func (m *mockKeyCloakClient) RemoveGroup(groupID string) {
	m.groupRemoved = groupID
}

func (m *mockKeyCloakClient) GetGroupID(groupName string) string {
	return ""
}

func (m *mockKeyCloakClient) GetGroupMembers(groupID string) []keycloak.GroupMember {
	groupMember1 := keycloak.GroupMember{ID: "J", UserName: "Jan"}
	groupMember2 := keycloak.GroupMember{ID: "P", UserName: "Piet"}
	return []keycloak.GroupMember{groupMember1, groupMember2}
}

func (m *mockKeyCloakClient) RemoveMemberFromGroup(userid string, groupID string) {
	m.memberRemoved = userid
}

func (m *mockKeyCloakClient) AddMemberToGroup(userid string, groupID string) {
	m.memberAdded = userid
}

func (m *mockKeyCloakClient) GetUserID(username string) string {
	return string(username[0])
}

func (m *mockKeyCloakClient) GetUser(username string) keycloak.User {
	return keycloak.User{}
}

func (m *mockKeyCloakClient) GetGroups() []keycloak.Group {
	groupAttributes := keycloak.GroupAttributes{Tenant: []string{"MyCustomer"}}
	groupA := keycloak.Group{ID: "A", Name: "GroupA", Attributes: groupAttributes}
	groupB := keycloak.Group{ID: "B", Name: "GroupB", Attributes: groupAttributes}
	return []keycloak.Group{groupA, groupB}
}

func Test_ReconcileGroups(t *testing.T) {

	tenantConfig := config.TenantConfig{
		Name:   "MyCustomer",
		Groups: []config.GroupConfig{{Name: "GroupA"}, {Name: "GroupC"}},
	}
	groupReconciler := Reconciler{
		TenantConfig: tenantConfig,
		KeyCloakAPI:  &mock,
	}

	groupReconciler.ReconcileGroups()

	if mock.groupRemoved != "B" {
		t.Errorf("Group with ID 'B' should have been removed")
	}
	if mock.groupCreated != "GroupC" {
		t.Errorf("Group with name 'GroupC' should have been created")
	}
}

func Test_ReconcileMembers(t *testing.T) {

	members := []string{"Jan", "Fred"}
	tenantConfig := config.TenantConfig{
		Name:   "MyCustomer",
		Groups: []config.GroupConfig{{Name: "GroupA", Members: members}},
	}
	groupReconciler := Reconciler{
		TenantConfig: tenantConfig,
		KeyCloakAPI:  &mock,
	}

	groupReconciler.ReconcileGroups()

	if mock.memberRemoved != "P" {
		t.Errorf("Member 'Piet' should have been removed")
	}
	if mock.memberAdded != "F" {
		t.Errorf("Member 'Fred' should have been added")
	}
}
