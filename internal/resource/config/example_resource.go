package config

import (
	"context"

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
	Id        types.String `tfsdk:"id"`
	StringVal types.String `tfsdk:"string_val"`
}

// GetSchema defines the schema for the resource.
func (r *exampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Example resource.",
		Attributes: map[string]schema.Attribute{
			"level_one": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"level_two": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"level_three": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"level_four_primary": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"level_four_primary_string": schema.StringAttribute{
												Optional:    true,
												Description: "Parent should be level_one.level_two.level_three.level_four_primary.",
											},
											"level_five": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"level_five_string": schema.StringAttribute{
														Optional:    true,
														Description: "Parent should be level_one.level_two.level_three.level_four_primary.level_five.",
													},
												},
												Optional:    true,
												Description: "Parent should be level_one.level_two.level_three.level_four_primary.",
											},
										},
										Optional: true,
									},
									"level_four_secondary": schema.StringAttribute{
										Optional: true,
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

// Metadata returns the resource type name.
func (r *exampleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (m *exampleResourceModel) SetId() {
	m.Id = types.StringValue("id")
}

func (r *exampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan exampleResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set computed id
	plan.SetId()
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
	state.SetId()
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
	plan.SetId()
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// No backend so no logic needed
func (r *exampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
