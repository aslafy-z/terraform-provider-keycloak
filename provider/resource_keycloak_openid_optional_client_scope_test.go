package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccKeycloakDataSourceOpenidOptionalClientScope_basic(t *testing.T) {
	t.Parallel()
	clientId := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakOpenidOptionalClientScope_basic(clientId),
				Check:  testAccCheckKeycloakOpenidClientHasOptionalScope("keycloak_openid_optional_client_scope"),
			},
		},
	})
}

func testAccCheckKeycloakOpenidClientHasOptionalScope(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		realm := rs.Primary.Attributes["realm_id"]
		clientScopeId := rs.Primary.Attributes["client_scope_id"]

		var client string
		if strings.HasPrefix(resourceName, "keycloak_openid_client.") {
			client = rs.Primary.Attributes["client_id"]
		} else {
			client = rs.Primary.ID
		}

		keycloakOptionalClientScopes, err := keycloakClient.GetOpenidClientOptionalScopes(testCtx, realm, client)

		if err != nil {
			return err
		}

		var found = false
		for _, keycloakOptionalScope := range keycloakOptionalClientScopes {
			if keycloakOptionalScope.Id == clientScopeId {
				found = true

				break
			}
		}

		if !found {
			return fmt.Errorf("default scope %s is not assigned to client", clientScopeId)
		}

		return nil
	}
}

func testAccKeycloakOpenidOptionalClientScope_basic(clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "openid_client_scope" {
  realm_id = data.keycloak_realm.realm.id
  name     = "groups"
}

resource "keycloak_openid_default_client_scope" "openid_default_client_scope" {
	realm_id        = data.keycloak_realm.realm.id
	client_scope_id = keycloak_openid_client_scope.openid_client_scope.id
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = data.keycloak_realm.realm.id
	client_id   = "%s"
	access_type = "CONFIDENTIAL"
}
`, testAccRealm.Realm, clientId)
}
