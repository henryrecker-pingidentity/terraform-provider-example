package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	_ resource.Resource = &unexpectedPlanExampleResource{}
)

func UnexpectedPlanExampleResource() resource.Resource {
	return &unexpectedPlanExampleResource{}
}

type unexpectedPlanExampleResource struct {
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
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"bool_one": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							//Default:  booldefault.StaticBool(false),
						},
						"bool_two": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							//Default:  booldefault.StaticBool(false),
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
	resp.State.Raw = req.Plan.Raw
}

func (r *unexpectedPlanExampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.State.Raw = req.State.Raw
}

func (r *unexpectedPlanExampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State.Raw = req.Plan.Raw
}

func (r *unexpectedPlanExampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
