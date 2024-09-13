package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &sensitiveResource{}
)

// SensitiveResource is a helper function to simplify the provider implementation.
func SensitiveResource() resource.Resource {
	return &sensitiveResource{}
}

// sensitiveResource is the resource implementation.
type sensitiveResource struct {
}

// Metadata returns the resource type name.
func (r *sensitiveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sensitive"
}

// GetSchema defines the schema for the resource.
func (r *sensitiveResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Sensitive resource.",
		Attributes: map[string]schema.Attribute{
			"tables": schema.SetNestedAttribute{
				Description: "List of configuration tables.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rows": schema.ListNestedAttribute{
							Description: "List of table rows.",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"sensitive_fields": schema.SetNestedAttribute{
										Description: "The sensitive configuration fields in the row.",
										Optional:    true,
										Computed:    true,
										Default: setdefault.StaticValue(types.SetValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
											"name":  types.StringType,
											"value": types.StringType,
										}}, nil)),
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Description: "The name of the configuration field.",
													Required:    true,
												},
												"value": schema.StringAttribute{
													Description: "The sensitive value for the configuration field.",
													Required:    true,
													Sensitive:   true,
												},
											},
										},
									},
									"default_row": schema.BoolAttribute{
										Description: "Whether this row is the default.",
										Computed:    true,
										Optional:    true,
										Default:     booldefault.StaticBool(false),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *sensitiveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.State.Raw = req.Plan.Raw
}

func (r *sensitiveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.State.Raw = req.State.Raw
}

func (r *sensitiveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State.Raw = req.Plan.Raw
}

func (r *sensitiveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
