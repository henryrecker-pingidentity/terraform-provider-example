package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &exampleResource{}

	resourceLinkAttrTypes = map[string]attr.Type{
		"id":       types.StringType,
		"location": types.StringType,
	}

	fieldAttrTypes = map[string]attr.Type{
		"name":      types.StringType,
		"value":     types.StringType,
		"inherited": types.BoolType,
	}

	sourceAttrTypes = map[string]attr.Type{
		"type": types.StringType,
		"id":   types.StringType,
	}

	attributeContractFulfillmentAttrTypes = map[string]attr.Type{
		"source": types.ObjectType{
			AttrTypes: sourceAttrTypes,
		},
		"value": types.StringType,
	}

	customAttrSourceAttrTypes = map[string]attr.Type{
		"type": types.StringType,
		"data_store_ref": types.ObjectType{
			AttrTypes: resourceLinkAttrTypes,
		},
		"id":          types.StringType,
		"description": types.StringType,
		"attribute_contract_fulfillment": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: attributeContractFulfillmentAttrTypes,
			},
		},
		"filter_fields": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"value": types.StringType,
					"name":  types.StringType,
				},
			},
		},
	}

	jdbcAttrSourceAttrTypes = map[string]attr.Type{
		"type": types.StringType,
		"data_store_ref": types.ObjectType{
			AttrTypes: resourceLinkAttrTypes,
		},
		"id":          types.StringType,
		"description": types.StringType,
		"attribute_contract_fulfillment": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: attributeContractFulfillmentAttrTypes,
			},
		},
		"schema": types.StringType,
		"table":  types.StringType,
		"column_names": types.ListType{
			ElemType: types.StringType,
		},
		"filter": types.StringType,
	}

	ldapAttrSourceAttrTypes = map[string]attr.Type{
		"type": types.StringType,
		"data_store_ref": types.ObjectType{
			AttrTypes: resourceLinkAttrTypes,
		},
		"id":          types.StringType,
		"description": types.StringType,
		"attribute_contract_fulfillment": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: attributeContractFulfillmentAttrTypes,
			},
		},
		"search_filter":          types.StringType,
		"search_scope":           types.StringType,
		"member_of_nested_group": types.BoolType,
		"base_dn":                types.StringType,
		"binary_attribute_settings": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"binary_encoding": types.StringType,
				},
			},
		},
		"search_attributes": types.ListType{
			ElemType: types.StringType,
		},
	}

	attributeSourcesElementAttrTypes = map[string]attr.Type{
		"custom_attribute_source": types.ObjectType{
			AttrTypes: customAttrSourceAttrTypes,
		},
		"jdbc_attribute_source": types.ObjectType{
			AttrTypes: jdbcAttrSourceAttrTypes,
		},
		"ldap_attribute_source": types.ObjectType{
			AttrTypes: ldapAttrSourceAttrTypes,
		},
	}

	conditionalCriteriaAttrTypes = map[string]attr.Type{
		"source": types.ObjectType{
			AttrTypes: sourceAttrTypes,
		},
		"attribute_name": types.StringType,
		"condition":      types.StringType,
		"value":          types.StringType,
		"error_result":   types.StringType,
	}

	expressionCriteriaAttrTypes = map[string]attr.Type{
		"expression":   types.StringType,
		"error_result": types.StringType,
	}

	issuanceCriteriaAttrTypes = map[string]attr.Type{
		"conditional_criteria": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: conditionalCriteriaAttrTypes,
			},
		},
		"expression_criteria": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: expressionCriteriaAttrTypes,
			},
		},
	}

	attributeMappingAttrTypes = map[string]attr.Type{
		/*"attribute_sources": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: attributeSourcesElementAttrTypes,
			},
		},
		"attribute_contract_fulfillment": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: attributeContractFulfillmentAttrTypes,
			},
		},*/
		"issuance_criteria": types.ObjectType{
			AttrTypes: issuanceCriteriaAttrTypes,
		},
		"inherited": types.BoolType,
	}

	attributeSourcesEmptyList, _    = types.ListValue(types.ObjectType{AttrTypes: attributeSourcesElementAttrTypes}, []attr.Value{})
	conditionalCriteriaEmptyList, _ = types.ListValue(types.ObjectType{AttrTypes: conditionalCriteriaAttrTypes}, []attr.Value{})
	expressionCriteriaEmtpyList     = types.ListNull(types.ObjectType{AttrTypes: expressionCriteriaAttrTypes})
	issuanceCriteriaEmptyObject, _  = types.ObjectValue(issuanceCriteriaAttrTypes, map[string]attr.Value{
		"conditional_criteria": conditionalCriteriaEmptyList,
		"expression_criteria":  expressionCriteriaEmtpyList,
	})
)

// ExampleResource is a helper function to simplify the provider implementation.
func ExampleResource() resource.Resource {
	return &exampleResource{}
}

// exampleResource is the resource implementation.
type exampleResource struct {
}

type exampleResourceModel struct {
	Id               types.String `tfsdk:"id"`
	StringVal        types.String `tfsdk:"string_val"`
	Fields           types.List   `tfsdk:"fields"`
	AttributeMapping types.Object `tfsdk:"attribute_mapping"`
}

// GetSchema defines the schema for the resource.
func (r *exampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Example resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Computed id",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"string_val": schema.StringAttribute{
				Description: "Optional string attribute",
				Optional:    true,
			},
			"fields": schema.ListNestedAttribute{
				Description: "List of configuration fields. This attribute will include any values set by default by PingFederate.",
				Computed:    true,
				Optional:    false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the configuration field.",
							Required:    true,
							Computed:    false,
						},
						"value": schema.StringAttribute{
							Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
							Required:    true,
							Computed:    false,
						},
						"inherited": schema.BoolAttribute{
							Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
							Required:    true,
							Computed:    false,
						},
					},
				},
			},
			"attribute_mapping": schema.SingleNestedAttribute{
				Description: "The attributes mapping from attribute sources to attribute targets.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					/*"attribute_sources": schema.ListNestedAttribute{
						Optional:    true,
						Computed:    true,
						Default:     listdefault.StaticValue(attributeSourcesEmptyList),
						Description: "A list of configured data stores to look up attributes from.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"custom_attribute_source": schema.SingleNestedAttribute{
									Optional:    true,
									Description: "The configured settings to look up attributes from an associated data store.",
									Attributes: map[string]schema.Attribute{
										//TODO only need type on ldap dat source, others are implicit. Make the others readonly
										"type": schema.StringAttribute{
											Description: "The data store type of this attribute source.",
											Required:    true,
											//TODO is this type attribute really required? Why are there 4 possible types and only 3 attribute source implementations
											Validators: []validator.String{
												stringvalidator.OneOf("LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"),
											},
										},
										//TODO use shared schema
										"data_store_ref": schema.SingleNestedAttribute{
											Description: "Reference to the associated data store.",
											Required:    true,
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Description: "The ID of the resource.",
													Required:    true,
												},
												"location": schema.StringAttribute{
													Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
													Optional:    false,
													Computed:    true,
												},
											},
										},
										"id": schema.StringAttribute{
											Description: "The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.",
											Optional:    true,
										},
										"description": schema.StringAttribute{
											Description: "The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping. Note: Required for APC-to-SP Adapter Mappings",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												stringplanmodifier.UseStateForUnknown(),
											},
										},
										"attribute_contract_fulfillment": schema.MapNestedAttribute{
											Description: "A list of mappings from attribute names to their fulfillment values. This field is only valid for the SP Connection's Browser SSO mappings",
											Optional:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"source": schema.SingleNestedAttribute{
														Description: "The attribute value source.",
														Required:    true,
														Attributes: map[string]schema.Attribute{
															"type": schema.StringAttribute{
																Required:    true,
																Description: "The source type of this key.",
																//TODO enum validator
															},
															"id": schema.StringAttribute{
																Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																Optional:    true,
															},
														},
													},
													"value": schema.StringAttribute{
														Description: "The value for this attribute.",
														Required:    true,
													},
												},
											},
										},
										"filter_fields": schema.ListNestedAttribute{
											Description: "The list of fields that can be used to filter a request to the custom data store.",
											Optional:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"value": schema.StringAttribute{
														Description: "The value of this field. Whether or not the value is required will be determined by plugin validation checks.",
														Optional:    true,
													},
													"name": schema.StringAttribute{
														Description: "The name of this field.",
														Required:    true,
													},
												},
											},
										},
									},
								},
								"jdbc_attribute_source": schema.SingleNestedAttribute{
									Optional:    true,
									Description: "The configured settings to look up attributes from a JDBC data store.",
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description: "The data store type of this attribute source.",
											Required:    true,
											//TODO is this type attribute really required? Why are there 4 possible types and only 3 attribute source implementations
											Validators: []validator.String{
												stringvalidator.OneOf("LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"),
											},
										},
										//TODO use shared schema
										"data_store_ref": schema.SingleNestedAttribute{
											Description: "Reference to the associated data store.",
											Required:    true,
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Description: "The ID of the resource.",
													Required:    true,
												},
												"location": schema.StringAttribute{
													Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
													Optional:    false,
													Computed:    true,
												},
											},
										},
										"id": schema.StringAttribute{
											Description: "The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.",
											Optional:    true,
										},
										"description": schema.StringAttribute{
											Description: "The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping. Note: Required for APC-to-SP Adapter Mappings",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												stringplanmodifier.UseStateForUnknown(),
											},
										},
										"attribute_contract_fulfillment": schema.MapNestedAttribute{
											Description: "A list of mappings from attribute names to their fulfillment values. This field is only valid for the SP Connection's Browser SSO mappings",
											Optional:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"source": schema.SingleNestedAttribute{
														Description: "The attribute value source.",
														Required:    true,
														Attributes: map[string]schema.Attribute{
															"type": schema.StringAttribute{
																Required:    true,
																Description: "The source type of this key.",
																//TODO enum validator
															},
															"id": schema.StringAttribute{
																Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																Optional:    true,
															},
														},
													},
													"value": schema.StringAttribute{
														Description: "The value for this attribute.",
														Required:    true,
													},
												},
											},
										},
										"schema": schema.StringAttribute{
											Description: "Lists the table structure that stores information within a database. Some databases, such as Oracle, require a schema for a JDBC query. Other databases, such as MySQL, do not require a schema.",
											Optional:    true,
										},
										"table": schema.StringAttribute{
											Description: "The name of the database table. The name is used to construct the SQL query to retrieve data from the data store.",
											Required:    true,
										},
										"column_names": schema.ListAttribute{
											Description: "A list of column names used to construct the SQL query to retrieve data from the specified table in the datastore.",
											ElementType: types.StringType,
											Optional:    true,
										},
										"filter": schema.StringAttribute{
											Description: "The JDBC WHERE clause used to query your data store to locate a user record.",
											Required:    true,
										},
									},
								},
								"ldap_attribute_source": schema.SingleNestedAttribute{
									Optional:    true,
									Description: "The configured settings to look up attributes from a LDAP data store.",
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description: "The data store type of this attribute source.",
											Required:    true,
											//TODO is this type attribute really required? Why are there 4 possible types and only 3 attribute source implementations
											Validators: []validator.String{
												stringvalidator.OneOf("LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"),
											},
										},
										//TODO use shared schema
										"data_store_ref": schema.SingleNestedAttribute{
											Description: "Reference to the associated data store.",
											Required:    true,
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Description: "The ID of the resource.",
													Required:    true,
												},
												"location": schema.StringAttribute{
													Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
													Optional:    false,
													Computed:    true,
												},
											},
										},
										"id": schema.StringAttribute{
											Description: "The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.",
											Optional:    true,
										},
										"description": schema.StringAttribute{
											Description: "The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping. Note: Required for APC-to-SP Adapter Mappings",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												stringplanmodifier.UseStateForUnknown(),
											},
										},
										"attribute_contract_fulfillment": schema.MapNestedAttribute{
											Description: "A list of mappings from attribute names to their fulfillment values. This field is only valid for the SP Connection's Browser SSO mappings",
											Optional:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"source": schema.SingleNestedAttribute{
														Description: "The attribute value source.",
														Required:    true,
														Attributes: map[string]schema.Attribute{
															"type": schema.StringAttribute{
																Required:    true,
																Description: "The source type of this key.",
																//TODO enum validator
															},
															"id": schema.StringAttribute{
																Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
																Optional:    true,
															},
														},
													},
													"value": schema.StringAttribute{
														Description: "The value for this attribute.",
														Required:    true,
													},
												},
											},
										},
										"search_filter": schema.StringAttribute{
											Description: "The LDAP filter that will be used to lookup the objects from the directory.",
											Required:    true,
										},
										"search_scope": schema.StringAttribute{
											Description: "Determines the node depth of the query.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOf("OBJECT", "ONE_LEVEL", "SUBTREE"),
											},
										},
										"member_of_nested_group": schema.BoolAttribute{
											Description: "Set this to true to return transitive group memberships for the 'memberOf' attribute. This only applies for Active Directory data sources. All other data sources will be set to false.",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.Bool{
												boolplanmodifier.UseStateForUnknown(),
											},
											Default: booldefault.StaticBool(false),
										},
										"base_dn": schema.StringAttribute{
											Description: "The base DN to search from. If not specified, the search will start at the LDAP's root.",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												stringplanmodifier.UseStateForUnknown(),
											},
										},
										"binary_attribute_settings": schema.MapNestedAttribute{
											Description: "The advanced settings for binary LDAP attributes.",
											Optional:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"binary_encoding": schema.StringAttribute{
														Optional:    true,
														Description: "Get the encoding type for this attribute. If not specified, the default is BASE64.",
														Validators: []validator.String{
															stringvalidator.OneOf("OBJECT", "ONE_LEVEL", "SUBTREE"),
														},
													},
												},
											},
										},
										"search_attributes": schema.ListAttribute{
											Description: "A list of LDAP attributes returned from search and available for mapping.",
											Optional:    true,
											ElementType: types.StringType,
										},
									},
								},
							},
						},
					},
					"attribute_contract_fulfillment": schema.MapNestedAttribute{
						Description: "A list of mappings from attribute names to their fulfillment values.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"source": schema.SingleNestedAttribute{
									Description: "The attribute value source.",
									Required:    true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description: "The source type of this key.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOf([]string{"TOKEN_EXCHANGE_PROCESSOR_POLICY", "ACCOUNT_LINK", "ADAPTER", "ASSERTION", "CONTEXT", "CUSTOM_DATA_STORE", "EXPRESSION", "JDBC_DATA_STORE", "LDAP_DATA_STORE", "PING_ONE_LDAP_GATEWAY_DATA_STORE", "MAPPED_ATTRIBUTES", "NO_MAPPING", "TEXT", "TOKEN", "REQUEST", "OAUTH_PERSISTENT_GRANT", "SUBJECT_TOKEN", "ACTOR_TOKEN", "PASSWORD_CREDENTIAL_VALIDATOR", "IDP_CONNECTION", "AUTHENTICATION_POLICY_CONTRACT", "CLAIMS", "LOCAL_IDENTITY_PROFILE", "EXTENDED_CLIENT_METADATA", "EXTENDED_PROPERTIES", "TRACKED_HTTP_PARAMS", "FRAGMENT", "INPUTS", "ATTRIBUTE_QUERY", "IDENTITY_STORE_USER", "IDENTITY_STORE_GROUP", "SCIM_USER", "SCIM_GROUP"}...),
											},
										},
										"id": schema.StringAttribute{
											Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
											Optional:    true,
										},
									},
								},
								"value": schema.StringAttribute{
									Description: "The value for this attribute.",
									Required:    true,
								},
							},
						},
					},*/
					"issuance_criteria": schema.SingleNestedAttribute{
						Description: "The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled.",
						Optional:    true,
						Computed:    true,
						Default:     objectdefault.StaticValue(issuanceCriteriaEmptyObject),
						Attributes: map[string]schema.Attribute{
							"conditional_criteria": schema.ListNestedAttribute{
								Description: "An issuance criterion that checks a source attribute against a particular condition and the expected value. If the condition is true then this issuance criterion passes, otherwise the criterion fails.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										//TODO share these definitions
										"source": schema.SingleNestedAttribute{
											Description: "The attribute value source.",
											Required:    true,
											Attributes: map[string]schema.Attribute{
												"type": schema.StringAttribute{
													Description: "The source type of this key.",
													Required:    true,
													Validators: []validator.String{
														stringvalidator.OneOf([]string{"TOKEN_EXCHANGE_PROCESSOR_POLICY", "ACCOUNT_LINK", "ADAPTER", "ASSERTION", "CONTEXT", "CUSTOM_DATA_STORE", "EXPRESSION", "JDBC_DATA_STORE", "LDAP_DATA_STORE", "PING_ONE_LDAP_GATEWAY_DATA_STORE", "MAPPED_ATTRIBUTES", "NO_MAPPING", "TEXT", "TOKEN", "REQUEST", "OAUTH_PERSISTENT_GRANT", "SUBJECT_TOKEN", "ACTOR_TOKEN", "PASSWORD_CREDENTIAL_VALIDATOR", "IDP_CONNECTION", "AUTHENTICATION_POLICY_CONTRACT", "CLAIMS", "LOCAL_IDENTITY_PROFILE", "EXTENDED_CLIENT_METADATA", "EXTENDED_PROPERTIES", "TRACKED_HTTP_PARAMS", "FRAGMENT", "INPUTS", "ATTRIBUTE_QUERY", "IDENTITY_STORE_USER", "IDENTITY_STORE_GROUP", "SCIM_USER", "SCIM_GROUP"}...),
													},
												},
												"id": schema.StringAttribute{
													Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
													Optional:    true,
												},
											},
										},
										"attribute_name": schema.StringAttribute{
											Description: "The name of the attribute to use in this issuance criterion.",
											Required:    true,
										},
										"condition": schema.StringAttribute{
											Description: "The condition that will be applied to the source attribute's value and the expected value.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOf([]string{"EQUALS", "EQUALS_CASE_INSENSITIVE", "EQUALS_DN", "NOT_EQUAL", "NOT_EQUAL_CASE_INSENSITIVE", "NOT_EQUAL_DN", "MULTIVALUE_CONTAINS", "MULTIVALUE_CONTAINS_CASE_INSENSITIVE", "MULTIVALUE_CONTAINS_DN", "MULTIVALUE_DOES_NOT_CONTAIN", "MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE", "MULTIVALUE_DOES_NOT_CONTAIN_DN"}...),
											},
										},
										"value": schema.StringAttribute{
											Description: "The expected value of this issuance criterion.",
											Required:    true,
										},
										"error_result": schema.StringAttribute{
											Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
											Optional:    true,
										},
									},
								},
							},
							"expression_criteria": schema.ListNestedAttribute{
								Description: "An issuance criterion that uses a Boolean return value from an OGNL expression to determine whether or not it passes.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"expression": schema.StringAttribute{
											Required:    true,
											Description: "The OGNL expression to evaluate.",
										},
										"error_result": schema.StringAttribute{
											Optional:    true,
											Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
										},
									},
								},
							},
						},
					},
					"inherited": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Whether this attribute mapping is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
					},
				},
			},
		},
	}
}

// Metadata returns the resource type name.
func (r *exampleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func AttributeContractFulfillmentValue(value string) types.Object {
	sourceVal, _ := types.ObjectValue(sourceAttrTypes, map[string]attr.Value{
		"type": types.StringValue("ADAPTER"),
		"id":   types.StringNull(),
	})
	objVal, _ := types.ObjectValue(attributeContractFulfillmentAttrTypes, map[string]attr.Value{
		"source": sourceVal,
		"value":  types.StringValue(value),
	})
	return objVal
}

func (m *exampleResourceModel) Populate(ctx context.Context) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	m.Id = types.StringValue("id")
	fields := []attr.Value{
		CreateField("field1", "val1", false),
		CreateField("field2", "val2", true),
		CreateField("field3", "val3", false),
	}
	m.Fields, diags = types.ListValue(types.ObjectType{AttrTypes: fieldAttrTypes}, fields)
	respDiags.Append(diags...)

	attributeMappingValues := map[string]attr.Value{
		"inherited": types.BoolValue(false),
	}

	// Build attribute_contract_fulfillment value
	/*attributeContractFulfillmentElementAttrTypes := attributeMappingAttrTypes["attribute_contract_fulfillment"].(types.MapType).ElemType.(types.ObjectType).AttrTypes
	attributeMappingValues["attribute_contract_fulfillment"], _ =
		types.MapValue(types.ObjectType{AttrTypes: attributeContractFulfillmentElementAttrTypes}, map[string]attr.Value{
			"entryUUID":     AttributeContractFulfillmentValue("entryUUID"),
			"policy.action": AttributeContractFulfillmentValue("policy.action"),
			"username":      AttributeContractFulfillmentValue("username"),
		})*/

	// Build issuance_criteria value
	conditional, diags := types.ListValue(types.ObjectType{AttrTypes: conditionalCriteriaAttrTypes}, []attr.Value{})
	respDiags.Append(diags...)
	attributeMappingValues["issuance_criteria"], diags = types.ObjectValue(
		issuanceCriteriaAttrTypes, map[string]attr.Value{
			"conditional_criteria": conditional,
			"expression_criteria":  types.ListNull(types.ObjectType{AttrTypes: expressionCriteriaAttrTypes}),
		})
	respDiags.Append(diags...)

	// Build attribute_sources value
	//TODO
	/*attrSourceElements := []attr.Value{}
	/*for _, attrSource := range r.AttributeMapping.AttributeSources {
		attrSourceValues := map[string]attr.Value{}
		if attrSource.CustomAttributeSource != nil {
			customAttrSourceValues := map[string]attr.Value{}
			customAttrSourceValues["filter_fields"], diags = types.ListValueFrom(ctx,
				customAttrSourceAttrTypes["filter_fields"].(types.ListType).ElemType, attrSource.CustomAttributeSource.FilterFields)
			respDiags.Append(diags...)

			customAttrSourceValues["type"] = types.StringValue(attrSource.CustomAttributeSource.Type)
			customAttrSourceValues["data_store_ref"], diags = types.ObjectValueFrom(ctx, internaltypes.ResourceLinkStateAttrType(), attrSource.CustomAttributeSource.DataStoreRef)
			respDiags.Append(diags...)
			customAttrSourceValues["id"] = types.StringPointerValue(attrSource.CustomAttributeSource.Id)
			customAttrSourceValues["description"] = types.StringPointerValue(attrSource.CustomAttributeSource.Description)
			customAttrSourceValues["attribute_contract_fulfillment"], diags = types.MapValueFrom(ctx,
				types.ObjectType{AttrTypes: attributeContractFulfillmentAttrTypes}, attrSource.CustomAttributeSource.AttributeContractFulfillment)
			respDiags.Append(diags...)
			attrSourceValues["custom_attribute_source"], diags = types.ObjectValue(customAttrSourceAttrTypes, customAttrSourceValues)
			respDiags.Append(diags...)
		} else {
			attrSourceValues["custom_attribute_source"] = types.ObjectNull(customAttrSourceAttrTypes)
		}
		if attrSource.JdbcAttributeSource != nil {
			jdbcAttrSourceValues := map[string]attr.Value{}
			jdbcAttrSourceValues["schema"] = types.StringPointerValue(attrSource.JdbcAttributeSource.Schema)
			jdbcAttrSourceValues["table"] = types.StringValue(attrSource.JdbcAttributeSource.Table)
			jdbcAttrSourceValues["column_names"], diags = types.ListValueFrom(ctx, types.StringType, attrSource.JdbcAttributeSource.ColumnNames)
			respDiags.Append(diags...)
			jdbcAttrSourceValues["filter"] = types.StringValue(attrSource.JdbcAttributeSource.Filter)

			jdbcAttrSourceValues["type"] = types.StringValue(attrSource.JdbcAttributeSource.Type)
			jdbcAttrSourceValues["data_store_ref"], diags = types.ObjectValueFrom(ctx, internaltypes.ResourceLinkStateAttrType(), attrSource.JdbcAttributeSource.DataStoreRef)
			respDiags.Append(diags...)
			jdbcAttrSourceValues["id"] = types.StringPointerValue(attrSource.JdbcAttributeSource.Id)
			jdbcAttrSourceValues["description"] = types.StringPointerValue(attrSource.JdbcAttributeSource.Description)
			jdbcAttrSourceValues["attribute_contract_fulfillment"], diags = types.MapValueFrom(ctx,
				types.ObjectType{AttrTypes: attributeContractFulfillmentAttrTypes}, attrSource.JdbcAttributeSource.AttributeContractFulfillment)
			respDiags.Append(diags...)
			attrSourceValues["jdbc_attribute_source"], diags = types.ObjectValue(jdbcAttrSourceAttrTypes, jdbcAttrSourceValues)
			respDiags.Append(diags...)
		} else {
			attrSourceValues["jdbc_attribute_source"] = types.ObjectNull(jdbcAttrSourceAttrTypes)
		}
		if attrSource.LdapAttributeSource != nil {
			ldapAttrSourceValues := map[string]attr.Value{}
			ldapAttrSourceValues["base_dn"] = types.StringPointerValue(attrSource.LdapAttributeSource.BaseDn)
			ldapAttrSourceValues["search_scope"] = types.StringValue(attrSource.LdapAttributeSource.SearchScope)
			ldapAttrSourceValues["search_filter"] = types.StringValue(attrSource.LdapAttributeSource.SearchFilter)
			ldapAttrSourceValues["search_attributes"], diags = types.ListValueFrom(ctx, types.StringType, attrSource.LdapAttributeSource.SearchAttributes)
			respDiags.Append(diags...)
			ldapAttrSourceValues["binary_attribute_settings"], diags = types.MapValueFrom(ctx,
				ldapAttrSourceAttrTypes["binary_attribute_settings"].(types.MapType).ElemType, attrSource.LdapAttributeSource.BinaryAttributeSettings)
			respDiags.Append(diags...)
			ldapAttrSourceValues["member_of_nested_group"] = types.BoolPointerValue(attrSource.LdapAttributeSource.MemberOfNestedGroup)

			ldapAttrSourceValues["type"] = types.StringValue(attrSource.LdapAttributeSource.Type)
			ldapAttrSourceValues["data_store_ref"], diags = types.ObjectValueFrom(ctx, internaltypes.ResourceLinkStateAttrType(), attrSource.LdapAttributeSource.DataStoreRef)
			respDiags.Append(diags...)
			ldapAttrSourceValues["id"] = types.StringPointerValue(attrSource.LdapAttributeSource.Id)
			ldapAttrSourceValues["description"] = types.StringPointerValue(attrSource.LdapAttributeSource.Description)
			ldapAttrSourceValues["attribute_contract_fulfillment"], diags = types.MapValueFrom(ctx,
				types.ObjectType{AttrTypes: attributeContractFulfillmentAttrTypes}, attrSource.LdapAttributeSource.AttributeContractFulfillment)
			respDiags.Append(diags...)
			attrSourceValues["ldap_attribute_source"], diags = types.ObjectValue(ldapAttrSourceAttrTypes, ldapAttrSourceValues)
			respDiags.Append(diags...)
		} else {
			attrSourceValues["ldap_attribute_source"] = types.ObjectNull(ldapAttrSourceAttrTypes)
		}
		attrSourceElement, objectValueFromDiags := types.ObjectValue(attributeSourcesElementAttrTypes, attrSourceValues)
		respDiags.Append(objectValueFromDiags...)
		attrSourceElements = append(attrSourceElements, attrSourceElement)
	}*/
	/*attributeMappingValues["attribute_sources"], diags = types.ListValue(types.ObjectType{AttrTypes: attributeSourcesElementAttrTypes}, attrSourceElements)
	respDiags.Append(diags...)*/

	// Build complete attribute mapping value
	m.AttributeMapping, diags = types.ObjectValue(attributeMappingAttrTypes, attributeMappingValues)
	respDiags.Append(diags...)
	return respDiags
}

func CreateField(name, value string, inherited bool) types.Object {
	val, _ := types.ObjectValue(fieldAttrTypes, map[string]attr.Value{
		"name":      types.StringValue(name),
		"value":     types.StringValue(value),
		"inherited": types.BoolValue(inherited),
	})
	return val
}

func (r *exampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan exampleResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set computed id
	resp.Diagnostics.Append(plan.Populate(ctx)...)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *exampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state exampleResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set computed id
	resp.Diagnostics.Append(state.Populate(ctx)...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *exampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan exampleResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set computed id
	resp.Diagnostics.Append(plan.Populate(ctx)...)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// No backend so no logic needed
func (r *exampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
