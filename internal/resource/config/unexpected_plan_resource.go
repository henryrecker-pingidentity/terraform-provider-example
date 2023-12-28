package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource = &unexpectedPlanExampleResource{}
)

func UnexpectedPlanExampleResource() resource.Resource {
	return &unexpectedPlanExampleResource{}
}

type unexpectedPlanExampleResource struct {
}

type unexpectedPlanExampleResourceModel struct {
	Readonly types.String `tfsdk:"readonly"`
	Set      types.Set    `tfsdk:"set"`
}

func (r *unexpectedPlanExampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"readonly": schema.StringAttribute{
				Optional: false,
				Computed: true,
			},
			"set": schema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (r *unexpectedPlanExampleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_unexpected_plan_example"
}

func (r *unexpectedPlanExampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model unexpectedPlanExampleResourceModel
	req.Plan.Get(ctx, &model)

	// Set the same value every time for the string
	model.Readonly = types.StringValue("examplereadonly")

	// Set the same value every time for the set. {{"name": "examplename"}}
	setResult := []attr.Value{}
	setVal, _ := types.ObjectValue(map[string]attr.Type{
		"name": types.StringType,
	},
		map[string]attr.Value{
			"name": types.StringValue("examplename"),
		})
	setResult = append(setResult, setVal)
	model.Set, _ = types.SetValue(types.ObjectType{AttrTypes: map[string]attr.Type{
		"name": types.StringType,
	}}, setResult)

	resp.State.Set(ctx, model)
}

func (r *unexpectedPlanExampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.State.Raw = req.State.Raw
}

func (r *unexpectedPlanExampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// not implemented since I'm just testing plans after the initial creation
	resp.State.Raw = req.Plan.Raw
}

func (r *unexpectedPlanExampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
