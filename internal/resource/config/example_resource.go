package config

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &exampleResource{}
	_ resource.ResourceWithImportState = &exampleResource{}

	// API always returns same response for this example
	apiResponseJson = `{
		"float64": 242.1,
		"computed": "computed value"
	}`
)

// ExampleResource is a helper function to simplify the provider implementation.
func ExampleResource() resource.Resource {
	return &exampleResource{}
}

// exampleResource is the resource implementation.
type exampleResource struct {
}

type apiResponseModel struct {
	Float64  float64 `json:"float64"`
	Computed string  `json:"computed"`
}

type exampleResourceModel struct {
	Float64  types.Float64 `tfsdk:"float64"`
	Computed types.String  `tfsdk:"computed"`
}

// GetSchema defines the schema for the resource.
func (r *exampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Example resource.",
		Attributes: map[string]schema.Attribute{
			"float64": schema.Float64Attribute{
				Required: true,
			},
			"computed": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Metadata returns the resource type name.
func (r *exampleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func parseApiResponse() (*exampleResourceModel, diag.Diagnostics) {
	var apiResp apiResponseModel
	var diags diag.Diagnostics

	err := json.Unmarshal([]byte(apiResponseJson), &apiResp)
	if err != nil {
		diags.AddError(
			"Error parsing API response",
			"Could not parse API response: "+err.Error(),
		)
		return nil, diags
	}

	result := &exampleResourceModel{
		Float64:  types.Float64Value(apiResp.Float64),
		Computed: types.StringValue(apiResp.Computed),
	}
	return result, diags
}

func (r *exampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	respModel, diags := parseApiResponse()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, respModel)...)
}

func (r *exampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	respModel, diags := parseApiResponse()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, respModel)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *exampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	respModel, diags := parseApiResponse()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, respModel)...)
}

// No backend so no logic needed
func (r *exampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// ImportState implements resource.ResourceWithImportState.
func (r *exampleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// No special import logic needed
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("computed"), types.StringValue("computed value"))...)
}
