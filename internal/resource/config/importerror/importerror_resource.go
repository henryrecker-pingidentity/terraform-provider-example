package importerror

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
	_ resource.Resource = &importErrorResource{}
)

// ImportErrorResource is a helper function to simplify the provider implementation.
func ImportErrorResource() resource.Resource {
	return &importErrorResource{}
}

// importErrorResource is the resource implementation.
type importErrorResource struct {
}

type importErrorResourceModel struct {
	Id      types.String `tfsdk:"test_id"`
	ListVal types.List   `tfsdk:"list_val"`
}

// GetSchema defines the schema for the resource.
func (r *importErrorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Import error resource.",
		Attributes: map[string]schema.Attribute{
			"test_id": schema.StringAttribute{
				Description: "Main id",
				Required:    true,
			},
			"list_val": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "List id",
							Optional:    true,
							Computed:    false,
						},
					},
				},
			},
		},
	}
}

// Metadata returns the resource type name.
func (r *importErrorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_import_error"
}

func (m *importErrorResourceModel) Populate() {
	m.Id = types.StringValue("id")
	listElement, _ := types.ObjectValue(map[string]attr.Type{
		"id": types.StringType,
	}, map[string]attr.Value{
		"id": types.StringValue("list_id"),
	})
	m.ListVal, _ = types.ListValue(types.ObjectType{AttrTypes: map[string]attr.Type{
		"id": types.StringType,
	}}, []attr.Value{listElement})
}

func (r *importErrorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state importErrorResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Populate()
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *importErrorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan importErrorResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Populate()
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *importErrorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan importErrorResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Populate()
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// No backend so no logic needed
func (r *importErrorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *importErrorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("test_id"), req, resp)
}
