package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &DevicesCellularSimsResource{}
	_ resource.ResourceWithConfigure   = &DevicesCellularSimsResource{}
	_ resource.ResourceWithImportState = &DevicesCellularSimsResource{}
)

func NewDevicesCellularSimsResource() resource.Resource {
	return &DevicesCellularSimsResource{}
}

// DevicesCellularSimsResource defines the resource implementation.
type DevicesCellularSimsResource struct {
	client *openApiClient.APIClient
}

// DevicesCellularSimsResourceModel describes the resource data model.
type DevicesCellularSimsResourceModel struct {
	Id          jsontypes.String `tfsdk:"id"`
	Serial      jsontypes.String `tfsdk:"serial" json:"serial"`
	Sims        []Sim            `tfsdk:"sims" json:"sims"`
	SimFailOver SimFailover      `tfsdk:"sim_failover" json:"simFailover"`
}

type Sim struct {
	Slot      jsontypes.String `tfsdk:"slot" json:"slot"`
	IsPrimary jsontypes.Bool   `tfsdk:"is_primary" json:"isPrimary"`
	Apns      []Apns           `tfsdk:"apns" json:"apns"`
}

type Apns struct {
	Name           jsontypes.String `tfsdk:"name" json:"name"`
	AllowedIpTypes []string         `tfsdk:"allowed_ip_types" json:"allowedIpTypes"`
	Authentication Authentication   `tfsdk:"authentication" json:"authentication"`
}

type Authentication struct {
	Password jsontypes.String `tfsdk:"password" json:"password"`
	Username jsontypes.String `tfsdk:"username" json:"username"`
	Type     jsontypes.String `tfsdk:"type" json:"type"`
}

type SimFailover struct {
	Enabled jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
}

func (r *DevicesCellularSimsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_cellular_sims"
}

func (r *DevicesCellularSimsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Updates the SIM and APN configurations for a cellular device.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "serial.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"sim_failover": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Failover to secondary SIM.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
				},
			},
			"sims": schema.SetNestedAttribute{
				MarkdownDescription: "List of SIMs. If a SIM was previously configured and not specified in this request, it will remain unchanged.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"slot": schema.StringAttribute{
							MarkdownDescription: "SIM slot being configured. Must be 'sim1' on single-sim devices.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"is_primary": schema.BoolAttribute{
							MarkdownDescription: "If true, this SIM is used for boot. Must be true on single-sim devices.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"apns": schema.SetNestedAttribute{
							MarkdownDescription: "APN configurations. If empty, the default APN will be used.",
							Required:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										MarkdownDescription: "serial.",
										Required:            true,
										CustomType:          jsontypes.StringType,
									},
									"allowed_ip_types": schema.SetAttribute{
										MarkdownDescription: "IP versions to support (permitted values include 'ipv4', 'ipv6').",
										Required:            true,
										ElementType:         jsontypes.StringType,
									},
									"authentication": schema.SingleNestedAttribute{
										Optional:            true,
										Computed:            true,
										MarkdownDescription: "APN authentication configurations.",
										Attributes: map[string]schema.Attribute{
											"password": schema.StringAttribute{
												MarkdownDescription: "APN password, if type is set (if APN password is not supplied, the password is left unchanged).",
												Optional:            true,
												Computed:            true,
												CustomType:          jsontypes.StringType,
											},
											"username": schema.StringAttribute{
												MarkdownDescription: "APN username, if type is set.",
												Optional:            true,
												Computed:            true,
												CustomType:          jsontypes.StringType,
											},
											"type": schema.StringAttribute{
												MarkdownDescription: "APN auth type. Valid values are 'chap', 'none', 'pap'.",
												Optional:            true,
												Computed:            true,
												CustomType:          jsontypes.StringType,
											},
										},
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

func (r *DevicesCellularSimsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DevicesCellularSimsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DevicesCellularSimsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateDeviceCellularSims := *openApiClient.NewUpdateDeviceCellularSimsRequest()

	if !data.SimFailOver.Enabled.IsUnknown() {
		var enabled openApiClient.UpdateDeviceCellularSimsRequestSimFailover
		enabled.SetEnabled(data.SimFailOver.Enabled.ValueBool())
		updateDeviceCellularSims.SetSimFailover(enabled)
	}

	if len(data.Sims) > 0 {
		var devicesSerialCellularSims []openApiClient.UpdateDeviceCellularSimsRequestSimsInner
		for _, attribute := range data.Sims {
			var devicesSerialCellularSim openApiClient.UpdateDeviceCellularSimsRequestSimsInner
			devicesSerialCellularSim.SetIsPrimary(attribute.IsPrimary.ValueBool())
			devicesSerialCellularSim.SetSlot(attribute.Slot.ValueString())
			if len(attribute.Apns) > 0 {
				var devicesSerialCellularSimsApns []openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInner
				for _, apn := range attribute.Apns {
					var devicesSerialCellularSimsApn openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInner
					devicesSerialCellularSimsApn.SetName(apn.Name.ValueString())
					devicesSerialCellularSimsApn.SetAllowedIpTypes(apn.AllowedIpTypes)
					var devicesSerialCellularSimsAuthentication openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInnerAuthentication
					devicesSerialCellularSimsAuthentication.SetPassword(apn.Authentication.Password.ValueString())
					devicesSerialCellularSimsAuthentication.SetUsername(apn.Authentication.Username.ValueString())
					devicesSerialCellularSimsAuthentication.SetType(apn.Authentication.Type.ValueString())
					devicesSerialCellularSimsApn.SetAuthentication(devicesSerialCellularSimsAuthentication)
					devicesSerialCellularSimsApns = append(devicesSerialCellularSimsApns, devicesSerialCellularSimsApn)
				}
				devicesSerialCellularSim.SetApns(devicesSerialCellularSimsApns)
			}
			devicesSerialCellularSims = append(devicesSerialCellularSims, devicesSerialCellularSim)
		}
		updateDeviceCellularSims.SetSims(devicesSerialCellularSims)
	}

	_, httpResp, err := r.client.CellularApi.UpdateDeviceCellularSims(context.Background(), data.Serial.ValueString()).UpdateDeviceCellularSimsRequest(updateDeviceCellularSims).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *DevicesCellularSimsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesCellularSimsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.CellularApi.GetDeviceCellularSims(context.Background(), data.Serial.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *DevicesCellularSimsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DevicesCellularSimsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateDeviceCellularSims := *openApiClient.NewUpdateDeviceCellularSimsRequest()

	if !data.SimFailOver.Enabled.IsUnknown() {
		var enabled openApiClient.UpdateDeviceCellularSimsRequestSimFailover
		enabled.SetEnabled(data.SimFailOver.Enabled.ValueBool())
		updateDeviceCellularSims.SetSimFailover(enabled)
	}

	if len(data.Sims) > 0 {
		var devicesSerialCellularSims []openApiClient.UpdateDeviceCellularSimsRequestSimsInner
		for _, attribute := range data.Sims {
			var devicesSerialCellularSim openApiClient.UpdateDeviceCellularSimsRequestSimsInner
			devicesSerialCellularSim.SetIsPrimary(attribute.IsPrimary.ValueBool())
			devicesSerialCellularSim.SetSlot(attribute.Slot.ValueString())
			if len(attribute.Apns) > 0 {
				var devicesSerialCellularSimsApns []openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInner
				for _, apn := range attribute.Apns {
					var devicesSerialCellularSimsApn openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInner
					devicesSerialCellularSimsApn.SetName(apn.Name.ValueString())
					devicesSerialCellularSimsApn.SetAllowedIpTypes(apn.AllowedIpTypes)
					var devicesSerialCellularSimsAuthentication openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInnerAuthentication
					devicesSerialCellularSimsAuthentication.SetPassword(apn.Authentication.Password.ValueString())
					devicesSerialCellularSimsAuthentication.SetUsername(apn.Authentication.Username.ValueString())
					devicesSerialCellularSimsAuthentication.SetType(apn.Authentication.Type.ValueString())
					devicesSerialCellularSimsApn.SetAuthentication(devicesSerialCellularSimsAuthentication)
					devicesSerialCellularSimsApns = append(devicesSerialCellularSimsApns, devicesSerialCellularSimsApn)
				}
				devicesSerialCellularSim.SetApns(devicesSerialCellularSimsApns)
			}
			devicesSerialCellularSims = append(devicesSerialCellularSims, devicesSerialCellularSim)
		}
		updateDeviceCellularSims.SetSims(devicesSerialCellularSims)
	}

	_, httpResp, err := r.client.CellularApi.UpdateDeviceCellularSims(context.Background(), data.Serial.ValueString()).UpdateDeviceCellularSimsRequest(updateDeviceCellularSims).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *DevicesCellularSimsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DevicesCellularSimsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateDeviceCellularSims := *openApiClient.NewUpdateDeviceCellularSimsRequest()

	if !data.SimFailOver.Enabled.IsUnknown() {
		var enabled openApiClient.UpdateDeviceCellularSimsRequestSimFailover
		enabled.SetEnabled(data.SimFailOver.Enabled.ValueBool())
		updateDeviceCellularSims.SetSimFailover(enabled)
	}

	if len(data.Sims) > 0 {
		var devicesSerialCellularSims []openApiClient.UpdateDeviceCellularSimsRequestSimsInner
		for _, attribute := range data.Sims {
			var devicesSerialCellularSim openApiClient.UpdateDeviceCellularSimsRequestSimsInner
			devicesSerialCellularSim.SetIsPrimary(attribute.IsPrimary.ValueBool())
			devicesSerialCellularSim.SetSlot(attribute.Slot.ValueString())
			if len(attribute.Apns) > 0 {
				var devicesSerialCellularSimsApns []openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInner
				for _, apn := range attribute.Apns {
					var devicesSerialCellularSimsApn openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInner
					devicesSerialCellularSimsApn.SetName(apn.Name.ValueString())
					devicesSerialCellularSimsApn.SetAllowedIpTypes(apn.AllowedIpTypes)
					var devicesSerialCellularSimsAuthentication openApiClient.UpdateDeviceCellularSimsRequestSimsInnerApnsInnerAuthentication
					devicesSerialCellularSimsAuthentication.SetPassword(apn.Authentication.Password.ValueString())
					devicesSerialCellularSimsAuthentication.SetUsername(apn.Authentication.Username.ValueString())
					devicesSerialCellularSimsAuthentication.SetType(apn.Authentication.Type.ValueString())
					devicesSerialCellularSimsApn.SetAuthentication(devicesSerialCellularSimsAuthentication)
					devicesSerialCellularSimsApns = append(devicesSerialCellularSimsApns, devicesSerialCellularSimsApn)
				}
				devicesSerialCellularSim.SetApns(devicesSerialCellularSimsApns)
			}
			devicesSerialCellularSims = append(devicesSerialCellularSims, devicesSerialCellularSim)
		}
		updateDeviceCellularSims.SetSims(devicesSerialCellularSims)
	}

	_, httpResp, err := r.client.CellularApi.UpdateDeviceCellularSims(context.Background(), data.Serial.ValueString()).UpdateDeviceCellularSimsRequest(updateDeviceCellularSims).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *DevicesCellularSimsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
