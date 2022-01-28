package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/keycloak/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidDefaultClientScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakOpenidDefaultClientScopeCreate,
		ReadContext:   resourceKeycloakOpenidDefaultClientScopesRead,
		DeleteContext: resourceKeycloakOpenidDefaultClientScopeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakOpenidDefaultClientScopeImport,
		},

		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_scope_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKeycloakOpenidDefaultClientScopeCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return diag.FromErr(keycloakClient.PutOpenidRealmDefaultClientScope(ctx, realmId, clientScopeId))
}

func resourceKeycloakOpenidDefaultClientScopesRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	clientScope, err := keycloakClient.GetOpenidRealmDefaultClientScope(ctx, realmId, clientScopeId)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	data.Set("client_scope_id", clientScope.Id)
	data.Set("client_scope_name", clientScope.Name)

	return nil
}

func resourceKeycloakOpenidDefaultClientScopeDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return diag.FromErr(keycloakClient.DeleteOpenidRealmDefaultClientScope(ctx, realmId, clientScopeId))
}

func resourceKeycloakOpenidDefaultClientScopeImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{openidClientScopeId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
