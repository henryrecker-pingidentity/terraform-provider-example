package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/henryrecker-pingidentity/terraform-provider-example/internal/resource/config/pluginconfiguration"
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

type sensitiveResourceModel struct {
	Configuration types.Object `tfsdk:"configuration"`
	ManagerID     types.String `tfsdk:"manager_id"`
	Name          types.String `tfsdk:"name"`
	ID            types.String `tfsdk:"id"`
}

// GetSchema defines the schema for the resource.
func (r *sensitiveResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Sensitive resource.",
		Attributes: map[string]schema.Attribute{
			"tables": pluginconfiguration.ToSchema(),
			"manager_id": schema.StringAttribute{
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
			},
		},
	}
}

// Metadata returns the resource type name.
func (r *sensitiveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sensitive"
}

func (r *sensitiveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model sensitiveResourceModel
	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	model.Configuration, diags = pluginconfiguration.ToState(model.Configuration)
	model.ID = types.StringValue(model.ManagerID.ValueString())
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)

	//resp.State.Raw = req.Plan.Raw
}

func (r *sensitiveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model sensitiveResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	model.Configuration, diags = pluginconfiguration.ToState(model.Configuration)
	model.ID = types.StringValue(model.ManagerID.ValueString())
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)

	//resp.State.Raw = req.State.Raw
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sensitiveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model sensitiveResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	model.Configuration, diags = pluginconfiguration.ToState(model.Configuration)
	model.ID = types.StringValue(model.ManagerID.ValueString())
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)

	//resp.State.Raw = req.Plan.Raw
}

// No backend so no logic needed
func (r *sensitiveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
