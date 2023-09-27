package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &testDataSource{}
	_ datasource.DataSourceWithConfigure = &testDataSource{}
)

// coffeesModel maps coffees schema data.
type testDataSourceModel struct {
	EventPattern types.String `tfsdk:"event_pattern"`
	Event        types.String `tfsdk:"event"`
	Match        types.Bool   `tfsdk:"match"`
}

// NewTestDataSource is a helper function to simplify the provider implementation.
func NewTestDataSource() datasource.DataSource {
	return &testDataSource{}
}

// testDataSource is the data source implementation.
type testDataSource struct {
	client *eventbridge.Client
}

// Metadata returns the data source type name.
func (d *testDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_test"
}

// Schema defines the schema for the data source.
func (d *testDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"event_pattern": schema.StringAttribute{
				Required: true,
			},
			"event": schema.StringAttribute{
				Required: true,
			},
			"match": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *testDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state testDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	eventPattern, err := strconv.Unquote(state.EventPattern.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to unquote `event_pattern`",
			fmt.Sprintf("error: %T\nevent_pattern: %v", err, eventPattern),
		)
		return
	}
	event, err := strconv.Unquote(state.Event.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to unquote `event`",
			fmt.Sprintf("error: %T\nevent: %v", err, event),
		)
		return
	}

	result, err := d.client.TestEventPattern(ctx, &eventbridge.TestEventPatternInput{
		EventPattern: &eventPattern,
		Event:        &event,
	}, func(opts *eventbridge.Options) {})

	// https://aws.github.io/aws-sdk-go-v2/docs/handling-errors/
	if err != nil {
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			resp.Diagnostics.AddError(
				"Failed to test the event_pattern",
				fmt.Sprintf("failed to call service: %s, operation: %s, error: %v\nevent_pattern: \n  %v\nevent: \n  %v", oe.Service(), oe.Operation(), oe.Unwrap(), eventPattern, event),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to test the event_pattern",
			fmt.Sprintf("error: %T\nevent_pattern: \n  %v\nevent: \n  %v", err, eventPattern, event),
		)
		return
	}

	if result.Result {
		state.Match = types.BoolValue(true)
	} else {
		state.Match = types.BoolValue(false)

	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *testDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*eventbridge.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type Configure",
			fmt.Sprintf("Expected *eventbridge.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
