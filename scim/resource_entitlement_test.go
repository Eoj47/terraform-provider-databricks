package scim

import (
	"testing"

	"github.com/databricks/terraform-provider-databricks/common"
	"github.com/databricks/terraform-provider-databricks/qa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var oldGroup = Group{
	Schemas:     []URN{"urn:ietf:params:scim:schemas:core:2.0:Group"},
	DisplayName: "Data Scientists",
	ID:          "abc",
	Entitlements: []ComplexValue{
		{
			Value: "allow-cluster-create",
		},
	},
}

var newGroup = Group{
	Schemas:     []URN{"urn:ietf:params:scim:schemas:core:2.0:Group"},
	DisplayName: "Data Scientists",
	ID:          "abc",
	Entitlements: []ComplexValue{
		{
			Value: "allow-cluster-create",
		},
		{
			Value: "allow-instance-pool-create",
		},
		{
			Value: "databricks-sql-access",
		},
	},
}

var addRequest = PatchRequestComplexValue([]patchOperation{
	{
		"add", "entitlements", []ComplexValue{
			{
				Value: "allow-cluster-create",
			},
			{
				Value: "allow-instance-pool-create",
			},
			{
				Value: "databricks-sql-access",
			},
		},
	},
})

var updateRequest = PatchRequestComplexValue([]patchOperation{
	{
		"remove", "entitlements", []ComplexValue{
			{
				Value: "allow-cluster-create",
			},
			{
				Value: "allow-instance-pool-create",
			},
			{
				Value: "databricks-sql-access",
			},
			{
				Value: "workspace-access",
			},
		},
	},
	{
		"add", "entitlements", []ComplexValue{
			{
				Value: "allow-cluster-create",
			},
			{
				Value: "allow-instance-pool-create",
			},
			{
				Value: "databricks-sql-access",
			},
		},
	},
})

var deleteRequest = PatchRequestComplexValue([]patchOperation{{"remove", "entitlements", []ComplexValue{
	{
		Value: "allow-cluster-create",
	},
}}})

func TestResourceEntitlementsGroupCreate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Groups/abc",
				Response: oldGroup,
			},
			{
				Method:          "PATCH",
				Resource:        "/api/2.0/preview/scim/v2/Groups/abc",
				ExpectedRequest: addRequest,
				Response: Group{
					ID: "abc",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Groups/abc",
				Response: newGroup,
			},
		},
		Resource: ResourceEntitlements(),
		HCL: `
		group_id = "abc"
		allow_instance_pool_create = true
		allow_cluster_create = true
		databricks_sql_access = true
		`,
		Create: true,
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "group/abc", d.Id())
	assert.Equal(t, true, d.Get("allow_cluster_create"))
	assert.Equal(t, true, d.Get("allow_instance_pool_create"))
	assert.Equal(t, true, d.Get("databricks_sql_access"))
}

func TestResourceEntitlementsGroupRead(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Groups/abc",
				Response: oldGroup,
			},
		},
		Resource: ResourceEntitlements(),
		HCL:      `group_id = "abc"`,
		New:      true,
		Read:     true,
		ID:       "group/abc",
	}.ApplyAndExpectData(t, map[string]any{
		"group_id":             "abc",
		"allow_cluster_create": true,
	})
}

func TestResourceEntitlementsGroupRead_Error(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Groups/abc",
				Status:   400,
				Response: common.APIErrorBody{
					ScimDetail: "Something",
					ScimStatus: "Else",
				},
			},
		},
		Resource: ResourceEntitlements(),
		New:      true,
		Read:     true,
		ID:       "group/abc",
		HCL:      `group_id = "abc"`,
	}.ExpectError(t, "Something")
}

func TestResourceEntitlementsGroupUpdate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Groups/abc",
				Response: oldGroup,
			},
			{
				Method:          "PATCH",
				Resource:        "/api/2.0/preview/scim/v2/Groups/abc",
				ExpectedRequest: updateRequest,
				Response: Group{
					ID: "abc",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Groups/abc",
				Response: newGroup,
			},
		},
		Resource: ResourceEntitlements(),
		Update:   true,
		ID:       "group/abc",
		InstanceState: map[string]string{
			"group_id":             "abc",
			"allow_cluster_create": "true",
		},
		HCL: `
		group_id    = "abc"
		allow_cluster_create = true
		allow_instance_pool_create = true
		databricks_sql_access = true
		`,
	}.Apply(t)
	require.NoError(t, err, err)
	assert.Equal(t, "group/abc", d.Id(), "Id should not be empty")
	assert.Equal(t, true, d.Get("allow_cluster_create"))
	assert.Equal(t, true, d.Get("allow_instance_pool_create"))
	assert.Equal(t, true, d.Get("databricks_sql_access"))
}

func TestResourceEntitlementsGroupDelete(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Groups/abc",
				Response: oldGroup,
			},
			{
				Method:          "PATCH",
				Resource:        "/api/2.0/preview/scim/v2/Groups/abc",
				ExpectedRequest: deleteRequest,
				Response: Group{
					ID: "abc",
				},
			},
		},
		Resource: ResourceEntitlements(),
		Delete:   true,
		ID:       "group/abc",
		InstanceState: map[string]string{
			"group_id":             "abc",
			"allow_cluster_create": "true",
		},
		HCL: `
		group_id    = "abc"
		allow_cluster_create = true
		`,
	}.Apply(t)
}

var oldUser = User{
	DisplayName: "Example user",
	Active:      true,
	UserName:    "me@example.com",
	ID:          "abc",
	Entitlements: []ComplexValue{
		{
			Value: "allow-cluster-create",
		},
	},
	Groups: []ComplexValue{
		{
			Display: "admins",
			Value:   "4567",
		},
		{
			Display: "ds",
			Value:   "9877",
		},
	},
	Roles: []ComplexValue{
		{
			Value: "a",
		},
		{
			Value: "b",
		},
	},
}

var newUser = User{
	DisplayName: "Example user",
	Active:      true,
	UserName:    "me@example.com",
	ID:          "abc",
	Entitlements: []ComplexValue{
		{
			Value: "allow-cluster-create",
		},
		{
			Value: "allow-instance-pool-create",
		},
		{
			Value: "databricks-sql-access",
		},
	},
	Groups: []ComplexValue{
		{
			Display: "admins",
			Value:   "4567",
		},
		{
			Display: "ds",
			Value:   "9877",
		},
	},
	Roles: []ComplexValue{
		{
			Value: "a",
		},
		{
			Value: "b",
		},
	},
}

func TestResourceEntitlementsUserCreate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Users/abc",
				Response: oldUser,
			},
			{
				Method:          "PATCH",
				Resource:        "/api/2.0/preview/scim/v2/Users/abc",
				ExpectedRequest: addRequest,
				Response: User{
					ID: "abc",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Users/abc",
				Response: newUser,
			},
		},
		Resource: ResourceEntitlements(),
		HCL: `
		user_id = "abc"
		allow_instance_pool_create = true
		allow_cluster_create = true
		databricks_sql_access = true
		`,
		Create: true,
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "user/abc", d.Id())
	assert.Equal(t, true, d.Get("allow_cluster_create"))
	assert.Equal(t, true, d.Get("allow_instance_pool_create"))
	assert.Equal(t, true, d.Get("databricks_sql_access"))
}

func TestResourceEntitlementsUserRead(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Users/abc",
				Response: oldUser,
			},
		},
		Resource: ResourceEntitlements(),
		HCL:      `user_id = "abc"`,
		New:      true,
		Read:     true,
		ID:       "user/abc",
	}.ApplyAndExpectData(t, map[string]any{
		"user_id":              "abc",
		"allow_cluster_create": true,
	})
}

func TestResourceEntitlementsUserRead_Error(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Users/abc",
				Status:   400,
				Response: common.APIErrorBody{
					ScimDetail: "Something",
					ScimStatus: "Else",
				},
			},
		},
		Resource: ResourceEntitlements(),
		New:      true,
		Read:     true,
		ID:       "user/abc",
		HCL:      `user_id = "abc"`,
	}.ExpectError(t, "Something")
}

func TestResourceEntitlementsUserUpdate_Error(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Users/abc",
				Status:   400,
				Response: common.APIErrorBody{
					ScimDetail: "Something",
					ScimStatus: "Else",
				},
			},
			{
				Method:          "PATCH",
				Resource:        "/api/2.0/preview/scim/v2/Users/abc",
				ExpectedRequest: updateRequest,
				Status:          400,
				Response: common.APIErrorBody{
					ScimDetail: "Something",
					ScimStatus: "Else",
				},
			},
		},
		Resource: ResourceEntitlements(),
		Update:   true,
		ID:       "user/abc",
		InstanceState: map[string]string{
			"user_id":              "abc",
			"allow_cluster_create": "true",
		},
		HCL: `
		user_id    = "abc"
		allow_cluster_create = true
		allow_instance_pool_create = true
		databricks_sql_access = true
		`,
	}.ExpectError(t, "Something")
}

func TestResourceEntitlementsUserUpdate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Users/abc",
				Response: oldUser,
			},
			{
				Method:          "PATCH",
				Resource:        "/api/2.0/preview/scim/v2/Users/abc",
				ExpectedRequest: updateRequest,
				Response: User{
					ID: "abc",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Users/abc",
				Response: newUser,
			},
		},
		Resource: ResourceEntitlements(),
		Update:   true,
		ID:       "user/abc",
		InstanceState: map[string]string{
			"user_id":              "abc",
			"allow_cluster_create": "true",
		},
		HCL: `
		user_id    = "abc"
		allow_cluster_create = true
		allow_instance_pool_create = true
		databricks_sql_access = true
		`,
	}.Apply(t)
	require.NoError(t, err, err)
	assert.Equal(t, "user/abc", d.Id(), "Id should not be empty")
	assert.Equal(t, true, d.Get("allow_cluster_create"))
	assert.Equal(t, true, d.Get("allow_instance_pool_create"))
	assert.Equal(t, true, d.Get("databricks_sql_access"))
}

func TestResourceEntitlementsUserDelete(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/Users/abc",
				Response: oldUser,
			},
			{
				Method:          "PATCH",
				Resource:        "/api/2.0/preview/scim/v2/Users/abc",
				ExpectedRequest: deleteRequest,
				Response: User{
					ID: "abc",
				},
			},
		},
		Resource: ResourceEntitlements(),
		Delete:   true,
		ID:       "user/abc",
		InstanceState: map[string]string{
			"user_id":              "abc",
			"allow_cluster_create": "true",
		},
		HCL: `
		user_id    = "abc"
		allow_cluster_create = true
		`,
	}.ApplyNoError(t)
}

func TestResourceEntitlementsSPNCreate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/ServicePrincipals/abc",
				Response: oldUser,
			},
			{
				Method:          "PATCH",
				Resource:        "/api/2.0/preview/scim/v2/ServicePrincipals/abc",
				ExpectedRequest: addRequest,
				Response: User{
					ID: "abc",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/ServicePrincipals/abc",
				Response: newUser,
			},
		},
		Resource: ResourceEntitlements(),
		HCL: `
		service_principal_id = "abc"
		allow_cluster_create = true
		allow_instance_pool_create = true
		databricks_sql_access = true
		`,
		Create: true,
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "spn/abc", d.Id())
	assert.Equal(t, true, d.Get("allow_cluster_create"))
	assert.Equal(t, true, d.Get("allow_instance_pool_create"))
	assert.Equal(t, true, d.Get("databricks_sql_access"))
}

func TestResourceEntitlementsSPNRead(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/ServicePrincipals/abc",
				Response: User{
					ID:            "abc",
					ApplicationID: "bcd",
					DisplayName:   "Example Service Principal",
					Active:        true,
					Entitlements: []ComplexValue{
						{
							Value: "allow-cluster-create",
						},
					},
				},
			},
		},
		Resource: ResourceEntitlements(),
		HCL:      `service_principal_id = "abc"`,
		New:      true,
		Read:     true,
		ID:       "spn/abc",
	}.ApplyAndExpectData(t, map[string]any{
		"service_principal_id": "abc",
		"allow_cluster_create": true,
	})
}

func TestResourceEntitlementsSPNRead_NotFound(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/ServicePrincipals/abc",
				Status:   404,
			},
		},
		Resource: ResourceEntitlements(),
		New:      true,
		Read:     true,
		Removed:  true,
		ID:       "spn/abc",
		HCL:      `service_principal_id = "abc"`,
	}.ApplyNoError(t)
}

func TestResourceEntitlementsSPNRead_Error(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/ServicePrincipals/abc",
				Status:   400,
				Response: common.APIErrorBody{
					ScimDetail: "Something",
					ScimStatus: "Else",
				},
			},
		},
		Resource: ResourceEntitlements(),
		New:      true,
		Read:     true,
		ID:       "spn/abc",
		HCL:      `service_principal_id = "abc"`,
	}.ExpectError(t, "Something")
}

func TestResourceEntitlementsSPNUpdate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/ServicePrincipals/abc",
				Response: oldUser,
			},
			{
				Method:          "PATCH",
				Resource:        "/api/2.0/preview/scim/v2/ServicePrincipals/abc",
				ExpectedRequest: updateRequest,
				Response: Group{
					ID: "abc",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/ServicePrincipals/abc",
				Response: newUser,
			},
		},
		Resource: ResourceEntitlements(),
		Update:   true,
		ID:       "spn/abc",
		InstanceState: map[string]string{
			"service_principal_id": "abc",
			"allow_cluster_create": "true",
		},
		HCL: `
		service_principal_id       = "abc"
		allow_cluster_create       = true
		allow_instance_pool_create = true
		databricks_sql_access      = true
		`,
	}.Apply(t)
	require.NoError(t, err, err)
	assert.Equal(t, "spn/abc", d.Id(), "Id should not be empty")
	assert.Equal(t, true, d.Get("allow_cluster_create"))
	assert.Equal(t, true, d.Get("allow_instance_pool_create"))
	assert.Equal(t, true, d.Get("databricks_sql_access"))
}

func TestResourceEntitlementsSPNDelete(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/preview/scim/v2/ServicePrincipals/abc",
				Response: oldUser,
			},
			{
				Method:          "PATCH",
				Resource:        "/api/2.0/preview/scim/v2/ServicePrincipals/abc",
				ExpectedRequest: deleteRequest,
				Response: User{
					ID: "abc",
				},
			},
		},
		Resource: ResourceEntitlements(),
		Delete:   true,
		ID:       "spn/abc",
		InstanceState: map[string]string{
			"service_principal_id": "abc",
			"allow_cluster_create": "true",
		},
		HCL: `
		service_principal_id = "abc"
		allow_cluster_create = true
		`,
	}.ApplyNoError(t)
}
