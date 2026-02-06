package config

import (
	"context"
	"fmt"
	"math/big"

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
	Number   types.Number `tfsdk:"number"`
	Computed types.String `tfsdk:"computed"`
}

// const problematicNumber = "242.08120431461208"

// GetSchema defines the schema for the resource.
func (r *exampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Example resource.",
		Attributes: map[string]schema.Attribute{
			"number": schema.NumberAttribute{
				Required: true,
			},
			"computed": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func parseNumber(str string) (types.Number, error) {
	apiRespBigFloat := new(big.Float)
	updatedFloat, ok := apiRespBigFloat.SetString(str)
	if !ok {
		return types.NumberNull(), fmt.Errorf("unable to parse number: %s", str)
	}
	return types.NumberValue(updatedFloat), nil
}

func planNumberToString(number types.Number) string {
	return number.ValueBigFloat().Text('f', 14)
}

// Metadata returns the resource type name.
func (r *exampleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (r *exampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan exampleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	planNumStr := planNumberToString(plan.Number)
	// fmt.Printf("Plan number as string: %s\n", planNumStr)

	numberVal, err := parseNumber(planNumStr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing number",
			err.Error(),
		)
		return
	}
	state := exampleResourceModel{
		Number:   numberVal,
		Computed: types.StringValue("computed value"),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *exampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data exampleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	planNumStr := planNumberToString(data.Number)

	apiRespBigFloat := new(big.Float)
	updatedFloat, _ := apiRespBigFloat.SetString(planNumStr)
	_, acc := updatedFloat.Float64()
	fmt.Printf("Acc: %s\n", acc)

	data.Number = types.NumberValue(updatedFloat)
	data.Computed = types.StringValue("computed value")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *exampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan exampleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	planNumStr := planNumberToString(plan.Number)
	// fmt.Printf("Plan number as string: %s\n", planNumStr)

	numberVal, err := parseNumber(planNumStr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing number",
			err.Error(),
		)
		return
	}
	state := exampleResourceModel{
		Number:   numberVal,
		Computed: types.StringValue("computed value"),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// No backend so no logic needed
func (r *exampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
