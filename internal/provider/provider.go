package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/henryrecker-pingidentity/terraform-provider-example/internal/resource/config"
)

// Ensure the implementation satisfies the expected interfacesß
var (
	_ provider.Provider = &exampleProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func NewFactory(version string) func() provider.Provider {
	return func() provider.Provider {
		return &exampleProvider{}
	}
}

// NewTestProvider is a helper function to simplify testing implementation.
func NewTestProvider() provider.Provider {
	return NewFactory("test")()
}

// exampleProvider is the provider implementation.
type exampleProvider struct {
}

// Metadata returns the provider type name.
func (p *exampleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "example"
}

// GetSchema defines the provider-level schema for configuration data.
// Schema defines the provider-level schema for configuration data.
func (p *exampleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
}

func (p *exampleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// DataSources defines the data sources implemented in the provider.
func (p *exampleProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *exampleProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		config.ExampleResource,
		config.SensitiveResource,
	}
}
