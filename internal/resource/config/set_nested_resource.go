package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

var (
	_ resource.Resource = &setNestedExampleResource{}
)

func SetNestedExampleResource() resource.Resource {
	return &setNestedExampleResource{}
}

type setNestedExampleResource struct {
}

func (r *setNestedExampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"set": schema.SetNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"bool_one": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						"bool_two": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
					},
				},
			},
		},
	}
}

func (r *setNestedExampleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_set_nested_example"
}

func (r *setNestedExampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.State.Raw = req.Plan.Raw
}

func (r *setNestedExampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.State.Raw = req.State.Raw
}

func (r *setNestedExampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State.Raw = req.Plan.Raw
}

func (r *setNestedExampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
