package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &sensitiveResource{}
)

var (
	fieldAttrTypes = map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}

	rowsSensitiveFieldsSplitAttrTypes = map[string]attr.Type{
		"fields":           types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"sensitive_fields": types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row":      types.BoolType,
	}
	rowsMergedFieldsAttrTypes = map[string]attr.Type{
		"fields":      types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row": types.BoolType,
	}

	tablesSensitiveFieldsSplitAttrTypes = map[string]attr.Type{
		"name": types.StringType,
		"rows": types.ListType{ElemType: types.ObjectType{AttrTypes: rowsSensitiveFieldsSplitAttrTypes}},
	}
	tablesMergedFieldsAttrTypes = map[string]attr.Type{
		"name": types.StringType,
		"rows": types.ListType{ElemType: types.ObjectType{AttrTypes: rowsMergedFieldsAttrTypes}},
	}

	configurationAttrTypes = map[string]attr.Type{
		/*"fields":           types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"sensitive_fields": types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"fields_all":       types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},*/
		"tables": types.SetType{ElemType: types.ObjectType{AttrTypes: tablesSensitiveFieldsSplitAttrTypes}},
		//"tables_all": types.SetType{ElemType: types.ObjectType{AttrTypes: tablesMergedFieldsAttrTypes}},
	}
)

// SensitiveResource is a helper function to simplify the provider implementation.
func SensitiveResource() resource.Resource {
	return &sensitiveResource{}
}

// sensitiveResource is the resource implementation.
type sensitiveResource struct {
}

type sensitiveResourceModel struct {
	/*Configuration types.Object `tfsdk:"configuration"`
	ManagerID     types.String `tfsdk:"manager_id"`
	Name          types.String `tfsdk:"name"`
	ID            types.String `tfsdk:"id"`*/
	Tables types.Set `tfsdk:"tables"`
}

// GetSchema defines the schema for the resource.
func (r *sensitiveResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fieldsSetDefault, _ := types.SetValue(types.ObjectType{AttrTypes: fieldAttrTypes}, nil)
	sensitiveFieldsSetDefault, _ := types.SetValue(types.ObjectType{AttrTypes: fieldAttrTypes}, nil)
	tablesSetDefault, _ := types.SetValue(types.ObjectType{AttrTypes: tablesSensitiveFieldsSplitAttrTypes}, nil)
	fieldsNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the configuration field.",
				Required:    true,
			},
			"value": schema.StringAttribute{
				Description: "The value for the configuration field.",
				Required:    true,
			},
		},
	}
	sensitiveFieldsNestedObject := schema.NestedAttributeObject{
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
	}
	tablesNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the table.",
				Required:    true,
			},
			"rows": schema.ListNestedAttribute{
				Description: "List of table rows.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"fields": schema.SetNestedAttribute{
							Description:  "The configuration fields in the row.",
							Optional:     true,
							Computed:     true,
							Default:      setdefault.StaticValue(fieldsSetDefault),
							NestedObject: fieldsNestedObject,
						},
						"sensitive_fields": schema.SetNestedAttribute{
							Description:  "The sensitive configuration fields in the row.",
							Optional:     true,
							Computed:     true,
							NestedObject: sensitiveFieldsNestedObject,
							Default:      setdefault.StaticValue(sensitiveFieldsSetDefault),
						},
						"default_row": schema.BoolAttribute{
							Description: "Whether this row is the default.",
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}

	resp.Schema = schema.Schema{
		Description: "Sensitive resource.",
		Attributes: map[string]schema.Attribute{
			"tables": schema.SetNestedAttribute{
				Description:  "List of configuration tables.",
				Computed:     true,
				Optional:     true,
				Default:      setdefault.StaticValue(tablesSetDefault),
				NestedObject: tablesNestedObject,
			},
			/*"manager_id": schema.StringAttribute{
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created. Must be alphanumeric, contain no spaces, and be less than 33 characters.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The plugin instance name. The name can be modified once the instance is created.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created. Must be alphanumeric, contain no spaces, and be less than 33 characters.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},*/
		},
	}
}

// Metadata returns the resource type name.
func (r *sensitiveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sensitive"
}

func (r *sensitiveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	/*var model sensitiveResourceModel
	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	model.Configuration, diags = pluginconfiguration.ToState(model.Configuration)
	model.ID = types.StringValue(model.ManagerID.ValueString())
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)*/

	resp.State.Raw = req.Plan.Raw
}

func (r *sensitiveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	/*var model sensitiveResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	model.Configuration, diags = pluginconfiguration.ToState(model.Configuration)
	model.ID = types.StringValue(model.ManagerID.ValueString())
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)*/

	resp.State.Raw = req.State.Raw
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sensitiveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	/*var model sensitiveResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	model.Configuration, diags = pluginconfiguration.ToState(model.Configuration)
	model.ID = types.StringValue(model.ManagerID.ValueString())
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)*/

	resp.State.Raw = req.Plan.Raw
}

// No backend so no logic needed
func (r *sensitiveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
