package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &awsEventPatternProvider{}
)

type awsEventPatternProviderModel struct {
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &awsEventPatternProvider{
			version: version,
		}
	}
}

// awsEventPatternProvider is the provider implementation.
type awsEventPatternProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (a *awsEventPatternProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "aws-event-pattern"
	resp.Version = a.version
}

// Schema defines the provider-level schema for configuration data.
func (a *awsEventPatternProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

// Configure prepares a testEventPattern API client for data sources and resources.
func (a *awsEventPatternProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var conf awsEventPatternProviderModel
	diags := req.Config.Get(ctx, &conf)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Using the SDK's default configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	client := eventbridge.NewFromConfig(cfg)

	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (a *awsEventPatternProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTestDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (a *awsEventPatternProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
