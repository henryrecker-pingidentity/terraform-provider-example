package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &exampleResource{}
)

// ExampleResource is a helper function to simplify the provider implementation.
func ExampleResource() resource.Resource {
	return &exampleResource{}
}

// exampleResource is the resource implementation.
type exampleResource struct {
}

type exampleResourceModel struct {
	StringVal types.String `tfsdk:"string_val"`
	MapVal    types.Map    `tfsdk:"map_val"`
}

// GetSchema defines the schema for the resource.
func (r *exampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Example resource.",
		Attributes: map[string]schema.Attribute{
			"string_val": schema.StringAttribute{
				Description: "Optional string attribute",
				Optional:    true,
			},
			"map_val": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"values": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "A Set of values",
						},
					},
				},
				Optional:    true,
				Description: "Extended Properties allows to store additional information for IdP/SP Connections. The names of these extended properties should be defined in /extendedProperties.",
			},
		},
	}
}

// Metadata returns the resource type name.
func (r *exampleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (r *exampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := exampleResourceModel{
		StringVal: types.StringValue("str"),
		MapVal: types.MapValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"values": types.SetType{
					ElemType: types.StringType,
				},
			}}, map[string]attr.Value{
			"Use Case": types.ObjectValueMust(map[string]attr.Type{
				"values": types.SetType{
					ElemType: types.StringType,
				},
			}, map[string]attr.Value{
				"values": types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("CIAM"),
				}),
			}),
		}),
	}
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *exampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	plan := exampleResourceModel{
		StringVal: types.StringValue("str"),
		MapVal: types.MapValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"values": types.SetType{
					ElemType: types.StringType,
				},
			}}, map[string]attr.Value{
			"Use Case": types.ObjectValueMust(map[string]attr.Type{
				"values": types.SetType{
					ElemType: types.StringType,
				},
			}, map[string]attr.Value{
				"values": types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("CIAM"),
				}),
			}),
		}),
	}
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *exampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := exampleResourceModel{
		StringVal: types.StringValue("str"),
		MapVal: types.MapValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"values": types.SetType{
					ElemType: types.StringType,
				},
			}}, map[string]attr.Value{
			"Use Case": types.ObjectValueMust(map[string]attr.Type{
				"values": types.SetType{
					ElemType: types.StringType,
				},
			}, map[string]attr.Value{
				"values": types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("CIAM"),
				}),
			}),
		}),
	}
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// No backend so no logic needed
func (r *exampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *exampleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("string_val"), req, resp)
}
