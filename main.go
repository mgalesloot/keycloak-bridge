package main

import (
	"flag"

	"keycloak-bridge/config"
	"keycloak-bridge/group"
	"keycloak-bridge/keycloak"
)

func main() {

	platformConfigFile := flag.String("p", "examples/platform.config.yaml", "Platform config file")
	tenantConfigFile := flag.String("t", "examples/tenant.config.yaml", "Tenant config file")
	flag.Parse()

	keyCloakClient := keycloak.NewKeyCloakClient(*platformConfigFile)

	groupReconciler := group.Reconciler{
		TenantConfig: config.LoadTenantConfig(*tenantConfigFile),
		KeyCloakAPI:  keyCloakClient,
	}

	groupReconciler.ReconcileGroups()
}
