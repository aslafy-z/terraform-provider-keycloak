package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/keycloak/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidOptionalClientScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakOpenidOptionalClientScopeCreate,
		ReadContext:   resourceKeycloakOpenidOptionalClientScopesRead,
		DeleteContext: resourceKeycloakOpenidOptionalClientScopeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakOpenidOptionalClientScopeImport,
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

func resourceKeycloakOpenidOptionalClientScopeCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return diag.FromErr(keycloakClient.PutOpenidRealmOptionalClientScope(ctx, realmId, clientScopeId))
}

func resourceKeycloakOpenidOptionalClientScopesRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	clientScope, err := keycloakClient.GetOpenidRealmOptionalClientScope(ctx, realmId, clientScopeId)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	data.Set("client_scope_id", clientScope.Id)
	data.Set("client_scope_name", clientScope.Name)

	return nil
}

func resourceKeycloakOpenidOptionalClientScopeDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientScopeId := data.Get("client_scope_id").(string)

	return diag.FromErr(keycloakClient.DeleteOpenidRealmOptionalClientScope(ctx, realmId, clientScopeId))
}

func resourceKeycloakOpenidOptionalClientScopeImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{openidClientScopeId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
