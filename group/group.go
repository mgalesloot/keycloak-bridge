package group

import (
	"log"
	"sync"

	"keycloak-bridge/config"
	"keycloak-bridge/keycloak"
)

// Reconciler configures groups and members in KeyCloak according to a YAML config
type Reconciler struct {
	TenantConfig config.TenantConfig
	KeyCloakAPI  keycloak.API
}

// ReconcileGroups read TenantConfig and create and deletes KeyCloak groups accordingly.
func (groupReconciler Reconciler) ReconcileGroups() {
	var wg sync.WaitGroup

	// find all KeyCloak groups
	keyCloakGroups := groupReconciler.KeyCloakAPI.GetGroups()

	// filter groups for current tenant
	tenantKeyCloakGroups := filterTenantGroups(groupReconciler.TenantConfig.Name, keyCloakGroups)

	// handle current groups in tenant config
	for _, groupConfig := range groupReconciler.TenantConfig.Groups {
		wg.Add(1)
		go func(groupConfig config.GroupConfig) {
			defer wg.Done()
			groupReconciler.reconcileGroup(groupConfig, tenantKeyCloakGroups)
		}(groupConfig)
	}

	// handle groups removed from tenant config
	for _, keyCloakGroup := range tenantKeyCloakGroups {
		if !groupReconciler.isKeyCloakGroupRequired(keyCloakGroup) {
			wg.Add(1)
			go func(keyCloakGroup keycloak.Group) {
				defer wg.Done()
				groupReconciler.removeGroup(keyCloakGroup)
			}(keyCloakGroup)
		}
	}
	wg.Wait()
}

func (groupReconciler Reconciler) removeGroup(keyCloakGroup keycloak.Group) {
	log.Printf("Remove KeyCloak group %s\n", keyCloakGroup.Name)
	groupReconciler.KeyCloakAPI.RemoveGroup(keyCloakGroup.ID)
}

func (groupReconciler Reconciler) isKeyCloakGroupRequired(keyCloakGroup keycloak.Group) bool {
	for _, groupConfig := range groupReconciler.TenantConfig.Groups {
		if groupConfig.Name == keyCloakGroup.Name {
			return true
		}
	}
	return false
}

func (groupReconciler Reconciler) reconcileGroup(groupConfig config.GroupConfig, tenantKeyCloakGroups []keycloak.Group) {
	log.Printf("Reconcile KeyCloak group %q\n", groupConfig.Name)

	keyCloakGroupID, exists := getGroupID(tenantKeyCloakGroups, groupConfig.Name)

	if !exists {
		log.Printf("Create KeyCloak group %s\n", groupConfig.Name)
		groupReconciler.KeyCloakAPI.CreateGroup(groupConfig.Name, groupReconciler.TenantConfig.Name)

		keyCloakGroupID = groupReconciler.KeyCloakAPI.GetGroupID(groupConfig.Name)
	}

	groupReconciler.groupMembers(groupConfig, keyCloakGroupID)
}

func (groupReconciler Reconciler) groupMembers(groupConfig config.GroupConfig, keyCloakGroupID string) {
	members := groupReconciler.KeyCloakAPI.GetGroupMembers(keyCloakGroupID)

	memberUserNames := getUserNames(members)
	log.Printf("Current members in Keycloak group %q: %q\n", groupConfig.Name, memberUserNames)

	membersToBeRemoved := getMembersToBeRemoved(groupConfig, members)
	for _, memberToBeRemoved := range membersToBeRemoved {
		log.Printf("Remove user %q from KeyCloak group %q\n", memberToBeRemoved.UserName, groupConfig.Name)
		groupReconciler.KeyCloakAPI.RemoveMemberFromGroup(memberToBeRemoved.ID, keyCloakGroupID)
	}

	membersToBeAdded := getMembersToBeAdded(groupConfig, memberUserNames)
	for _, memberToBeAdded := range membersToBeAdded {
		groupReconciler.addMemberToGroup(groupConfig, keyCloakGroupID, memberToBeAdded)
	}
}

func (groupReconciler Reconciler) addMemberToGroup(groupConfig config.GroupConfig, keyCloakGroupID string, member string) {
	log.Printf("Add user %q to KeyCloak group %q\n", member, groupConfig.Name)
	if keyCloakUserID := groupReconciler.KeyCloakAPI.GetUserID(member); keyCloakUserID != "" {
		groupReconciler.KeyCloakAPI.AddMemberToGroup(keyCloakUserID, keyCloakGroupID)
	}
}

func filterTenantGroups(tenant string, keyCloakGroups []keycloak.Group) []keycloak.Group {
	tenantKeyCloakGroups := make([]keycloak.Group, 0)
	for _, keyCloakGroup := range keyCloakGroups {
		if len(keyCloakGroup.Attributes.Tenant) > 0 && keyCloakGroup.Attributes.Tenant[0] == tenant {
			tenantKeyCloakGroups = append(tenantKeyCloakGroups, keyCloakGroup)
		}
	}
	return tenantKeyCloakGroups
}

func getUserNames(currentMembers []keycloak.GroupMember) []string {
	usernames := make([]string, len(currentMembers))
	for i, member := range currentMembers {
		usernames[i] = member.UserName
	}
	return usernames
}

func getMembersToBeRemoved(groupConfig config.GroupConfig, currentMembers []keycloak.GroupMember) []keycloak.GroupMember {
	var membersToBeRemoved []keycloak.GroupMember
	for _, member := range currentMembers {
		if !stringInSlice(member.UserName, groupConfig.Members) {
			membersToBeRemoved = append(membersToBeRemoved, member)
		}
	}
	return membersToBeRemoved
}

func getMembersToBeAdded(groupConfig config.GroupConfig, currentMembers []string) []string {
	var membersToBeAdded []string
	for _, member := range groupConfig.Members {
		if !stringInSlice(member, currentMembers) {
			membersToBeAdded = append(membersToBeAdded, member)
		}
	}
	return membersToBeAdded
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getGroupID(groups []keycloak.Group, groupName string) (keyCloakGroupID string, exists bool) {
	for _, group := range groups {
		if group.Name == groupName {
			keyCloakGroupID = group.ID
			exists = true
			return
		}
	}
	return
}
