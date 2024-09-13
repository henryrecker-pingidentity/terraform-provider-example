package pluginconfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	fieldAttrTypes = map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}

	rowsSensitiveFieldsSplitAttrTypes = map[string]attr.Type{
		"fields":           types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"sensitive_fields": types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row":      types.BoolType,
	}
	rowsMergedFieldsAttrTypes = map[string]attr.Type{
		"fields":      types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row": types.BoolType,
	}

	tablesSensitiveFieldsSplitAttrTypes = map[string]attr.Type{
		"name": types.StringType,
		"rows": types.ListType{ElemType: types.ObjectType{AttrTypes: rowsSensitiveFieldsSplitAttrTypes}},
	}
	tablesMergedFieldsAttrTypes = map[string]attr.Type{
		"name": types.StringType,
		"rows": types.ListType{ElemType: types.ObjectType{AttrTypes: rowsMergedFieldsAttrTypes}},
	}

	configurationAttrTypes = map[string]attr.Type{
		/*"fields":           types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"sensitive_fields": types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"fields_all":       types.SetType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},*/
		"tables": types.SetType{ElemType: types.ObjectType{AttrTypes: tablesSensitiveFieldsSplitAttrTypes}},
		//"tables_all": types.SetType{ElemType: types.ObjectType{AttrTypes: tablesMergedFieldsAttrTypes}},
	}
)

func AttrTypes() map[string]attr.Type {
	return configurationAttrTypes
}

type pfConfigurationFieldsResult struct {
	plannedCleartextFields types.Set
	plannedSensitiveFields types.Set
	allFields              types.Set
}

type pfConfigurationRowsResult struct {
	allRowsSensitiveFieldsSplit types.List
	allRowsMergedFields         types.List
}

type pfConfigurationTablesResult struct {
	plannedTables         types.Set
	allTablesMergedFields types.Set
}

func ToSchema() schema.SingleNestedAttribute {
	fieldsSetDefault, _ := types.SetValue(types.ObjectType{AttrTypes: fieldAttrTypes}, nil)
	sensitiveFieldsSetDefault, _ := types.SetValue(types.ObjectType{AttrTypes: fieldAttrTypes}, nil)
	tablesSetDefault, _ := types.SetValue(types.ObjectType{AttrTypes: tablesSensitiveFieldsSplitAttrTypes}, nil)
	fieldsNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the configuration field.",
				Required:    true,
			},
			"value": schema.StringAttribute{
				Description: "The value for the configuration field.",
				Required:    true,
			},
		},
	}
	sensitiveFieldsNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the configuration field.",
				Required:    true,
			},
			"value": schema.StringAttribute{
				Description: "The sensitive value for the configuration field.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
	tablesNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the table.",
				Required:    true,
			},
			"rows": schema.ListNestedAttribute{
				Description: "List of table rows.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"fields": schema.SetNestedAttribute{
							Description:  "The configuration fields in the row.",
							Optional:     true,
							Computed:     true,
							Default:      setdefault.StaticValue(fieldsSetDefault),
							NestedObject: fieldsNestedObject,
						},
						"sensitive_fields": schema.SetNestedAttribute{
							Description:  "The sensitive configuration fields in the row.",
							Optional:     true,
							Computed:     true,
							NestedObject: sensitiveFieldsNestedObject,
							Default:      setdefault.StaticValue(sensitiveFieldsSetDefault),
						},
						"default_row": schema.BoolAttribute{
							Description: "Whether this row is the default.",
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
	/*tablesAllNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the table.",
				Required:    true,
			},
			"rows": schema.ListNestedAttribute{
				Description: "List of table rows.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"fields": schema.SetNestedAttribute{
							Description:  "The configuration fields in the row.",
							Optional:     true,
							NestedObject: fieldsNestedObject,
						},
						"default_row": schema.BoolAttribute{
							Description: "Whether this row is the default.",
							Computed:    true,
							Optional:    true,
							Default:     booldefault.StaticBool(false),
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}*/
	return schema.SingleNestedAttribute{
		Description: "Plugin instance configuration.",
		Required:    true,
		Validators: []validator.Object{
			noDuplicateFields(),
		},
		Attributes: map[string]schema.Attribute{
			"tables": schema.SetNestedAttribute{
				Description:  "List of configuration tables.",
				Computed:     true,
				Optional:     true,
				Default:      setdefault.StaticValue(tablesSetDefault),
				NestedObject: tablesNestedObject,
			},
			/*"tables_all": schema.SetNestedAttribute{
				Description:  "List of configuration tables. This attribute will include any values set by default by PingFederate.",
				Computed:     true,
				Optional:     false,
				NestedObject: tablesAllNestedObject,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},*/
			/*"fields": schema.SetNestedAttribute{
				Description:  "List of configuration fields.",
				Computed:     true,
				Optional:     true,
				Default:      setdefault.StaticValue(fieldsSetDefault),
				NestedObject: fieldsNestedObject,
			},
			"sensitive_fields": schema.SetNestedAttribute{
				Description:  "List of sensitive configuration fields.",
				Computed:     true,
				Optional:     true,
				Default:      setdefault.StaticValue(sensitiveFieldsSetDefault),
				NestedObject: sensitiveFieldsNestedObject,
			},
			"fields_all": schema.SetNestedAttribute{
				Description:  "List of configuration fields. This attribute will include any values set by default by PingFederate.",
				Computed:     true,
				Optional:     false,
				NestedObject: fieldsNestedObject,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},*/
		},
	}
}

func readFieldsResponse(planFields, planSensitiveFields *types.Set, diags *diag.Diagnostics) pfConfigurationFieldsResult {
	plannedCleartextFields := []attr.Value{}
	plannedSensitiveFields := []attr.Value{}
	allFields := []attr.Value{}
	plannedFieldsValues := map[string]*string{}
	plannedSensitiveFieldsValues := map[string]*string{}
	// Build up a map of all the values from the plan
	if planFields != nil {
		for _, planField := range planFields.Elements() {
			planFieldObj := planField.(types.Object)
			plannedFieldsValues[planFieldObj.Attributes()["name"].(types.String).ValueString()] =
				planFieldObj.Attributes()["value"].(types.String).ValueStringPointer()
		}
	}
	if planSensitiveFields != nil {
		for _, planField := range planSensitiveFields.Elements() {
			planFieldObj := planField.(types.Object)
			plannedSensitiveFieldsValues[planFieldObj.Attributes()["name"].(types.String).ValueString()] =
				planFieldObj.Attributes()["value"].(types.String).ValueStringPointer()
		}
	}

	for _, field := range planFields.Elements() {
		plannedCleartextFields = append(plannedCleartextFields, field)
		allFields = append(allFields, field)
	}
	for _, field := range planSensitiveFields.Elements() {
		plannedSensitiveFields = append(plannedSensitiveFields, field)
		allFields = append(allFields, field)
	}

	plannedCleartextFieldsSet, respDiags := types.SetValue(types.ObjectType{
		AttrTypes: fieldAttrTypes,
	}, plannedCleartextFields)
	diags.Append(respDiags...)
	plannedSensitiveFieldsSet, respDiags := types.SetValue(types.ObjectType{
		AttrTypes: fieldAttrTypes,
	}, plannedSensitiveFields)
	diags.Append(respDiags...)

	allFieldsSet, respDiags := types.SetValue(types.ObjectType{
		AttrTypes: fieldAttrTypes,
	}, allFields)
	diags.Append(respDiags...)

	return pfConfigurationFieldsResult{
		plannedCleartextFields: plannedCleartextFieldsSet,
		plannedSensitiveFields: plannedSensitiveFieldsSet,
		allFields:              allFieldsSet,
	}
}

func readRowsResponse(planRows *types.List, diags *diag.Diagnostics) pfConfigurationRowsResult {
	var rowsMergedFields, rowsSensitiveFieldsSplit []attr.Value
	// This is assuming there are never any rows added by the PF API. If there
	// are ever rows added, this will cause a nil pointer exception trying to read
	// index i of planRowsElements.
	for _, row := range planRows.Elements() {
		rowAttrs := row.(types.Object).Attributes()
		attrValues := map[string]attr.Value{}
		attrValuesSensitiveSplit := map[string]attr.Value{}
		attrValues["default_row"] = rowAttrs["default_row"]
		attrValuesSensitiveSplit["default_row"] = rowAttrs["default_row"]
		var planRowFields, planRowSensitiveFields *types.Set
		planRowFieldsVal, ok := rowAttrs["fields"]
		if ok {
			setVal := planRowFieldsVal.(types.Set)
			planRowFields = &setVal
		}
		planRowSensitiveFieldsVal, ok := rowAttrs["sensitive_fields"]
		if ok {
			setVal := planRowSensitiveFieldsVal.(types.Set)
			planRowSensitiveFields = &setVal
		}

		rowFields := readFieldsResponse(planRowFields, planRowSensitiveFields, diags)
		attrValues["fields"] = rowFields.allFields
		attrValuesSensitiveSplit["fields"] = rowFields.plannedCleartextFields
		attrValuesSensitiveSplit["sensitive_fields"] = rowFields.plannedSensitiveFields

		rowMergedFields, respDiags := types.ObjectValue(rowsMergedFieldsAttrTypes, attrValues)
		diags.Append(respDiags...)
		rowsMergedFields = append(rowsMergedFields, rowMergedFields)
		rowSensitiveFieldsSplit, respDiags := types.ObjectValue(rowsSensitiveFieldsSplitAttrTypes, attrValuesSensitiveSplit)
		diags.Append(respDiags...)
		rowsSensitiveFieldsSplit = append(rowsSensitiveFieldsSplit, rowSensitiveFieldsSplit)
	}

	rowsMergedFieldsList, respDiags := types.ListValue(types.ObjectType{
		AttrTypes: rowsMergedFieldsAttrTypes,
	}, rowsMergedFields)
	diags.Append(respDiags...)
	rowsSensitiveFieldsSplitList, respDiags := types.ListValue(types.ObjectType{
		AttrTypes: rowsSensitiveFieldsSplitAttrTypes,
	}, rowsSensitiveFieldsSplit)
	diags.Append(respDiags...)
	return pfConfigurationRowsResult{
		allRowsSensitiveFieldsSplit: rowsSensitiveFieldsSplitList,
		allRowsMergedFields:         rowsMergedFieldsList,
	}
}

func toTablesSetValue(planTables *types.Set, diags *diag.Diagnostics) pfConfigurationTablesResult {
	// List of *all* tables values to return
	allTablesMergedFields := []attr.Value{}
	// List of tables values to return that were expected based on the plan
	plannedTables := []attr.Value{}

	for _, table := range planTables.Elements() {
		tableAttrs := table.(types.Object).Attributes()
		attrValues := map[string]attr.Value{}
		attrValuesSensitiveSplit := map[string]attr.Value{}
		attrValues["name"] = tableAttrs["name"]
		attrValuesSensitiveSplit["name"] = tableAttrs["name"]
		// If this table was in the plan, pass in the planned rows when getting the 'rows' values in case there are some encrypted values
		// that aren't returned by the PF API
		var planTableRows *types.List
		planTableRowsVal, ok := tableAttrs["rows"]
		if ok {
			listValue := planTableRowsVal.(types.List)
			planTableRows = &listValue
		}

		tableRows := readRowsResponse(planTableRows, diags)
		attrValues["rows"] = tableRows.allRowsMergedFields
		attrValuesSensitiveSplit["rows"] = tableRows.allRowsSensitiveFieldsSplit

		tableMergedFields, respDiags := types.ObjectValue(tablesMergedFieldsAttrTypes, attrValues)
		diags.Append(respDiags...)
		allTablesMergedFields = append(allTablesMergedFields, tableMergedFields)
		tableSensitiveFieldsSplit, respDiags := types.ObjectValue(tablesSensitiveFieldsSplitAttrTypes, attrValuesSensitiveSplit)
		diags.Append(respDiags...)
		plannedTables = append(plannedTables, tableSensitiveFieldsSplit)

	}

	allTablesMergedFieldsSet, respDiags := types.SetValue(types.ObjectType{
		AttrTypes: tablesMergedFieldsAttrTypes,
	}, allTablesMergedFields)
	diags.Append(respDiags...)
	plannedTablesSet, respDiags := types.SetValue(types.ObjectType{
		AttrTypes: tablesSensitiveFieldsSplitAttrTypes,
	}, plannedTables)
	diags.Append(respDiags...)

	return pfConfigurationTablesResult{
		plannedTables:         plannedTablesSet,
		allTablesMergedFields: allTablesMergedFieldsSet,
	}
}

func ToState(configFromPlan types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	//var planFields, planSensitiveFields *types.Set
	var planTables *types.Set

	/*planFieldsValue, ok := configFromPlan.Attributes()["fields"]
	if ok {
		setVal := planFieldsValue.(types.Set)
		planFields = &setVal
	}
	planSensitiveFieldsValue, ok := configFromPlan.Attributes()["sensitive_fields"]
	if ok {
		setVal := planSensitiveFieldsValue.(types.Set)
		planSensitiveFields = &setVal
	}*/
	planTablesValue, ok := configFromPlan.Attributes()["tables"]
	if ok {
		setVal := planTablesValue.(types.Set)
		planTables = &setVal
	}

	//fields := readFieldsResponse(planFields, planSensitiveFields, &diags)
	tables := toTablesSetValue(planTables, &diags)

	//fieldsAttrValue := fields.plannedCleartextFields
	//sensitiveFieldsAttrValue := fields.plannedSensitiveFields
	tablesAttrValue := tables.plannedTables

	configurationAttrValue := map[string]attr.Value{
		/*"fields":           fieldsAttrValue,
		"sensitive_fields": sensitiveFieldsAttrValue,
		"fields_all":       fields.allFields,*/
		"tables": tablesAttrValue,
		//"tables_all": tables.allTablesMergedFields,
	}
	configObj, valueFromDiags := types.ObjectValue(configurationAttrTypes, configurationAttrValue)
	diags.Append(valueFromDiags...)
	return configObj, diags
}

// Mark fields_all and tables_all configuration as unknown if the fields and tables have changed in the plan
// func MarkComputedAttrsUnknownOnChange(planConfiguration, stateConfiguration types.Object) (types.Object, diag.Diagnostics) {
// 	if planConfiguration.IsNull() || planConfiguration.IsUnknown() || !internaltypes.IsDefined(stateConfiguration) {
// 		return planConfiguration, nil
// 	}
// 	planConfigurationAttrs := planConfiguration.Attributes()
// 	planFields := planConfiguration.Attributes()["fields"]
// 	stateFields := stateConfiguration.Attributes()["fields"]
// 	if !planFields.Equal(stateFields) {
// 		planConfigurationAttrs["fields_all"] = types.SetUnknown(types.ObjectType{AttrTypes: fieldAttrTypes})
// 	}

// 	planTables := planConfiguration.Attributes()["tables"]
// 	stateTables := stateConfiguration.Attributes()["tables"]
// 	if !planTables.Equal(stateTables) {
// 		planConfigurationAttrs["tables_all"] = types.ListUnknown(types.ObjectType{AttrTypes: tablesMergedFieldsAttrTypes})
// 	}

// 	return types.ObjectValue(configurationAttrTypes, planConfigurationAttrs)
// }

// // Mark fields_all and tables_all configuration as unknown
// func MarkComputedAttrsUnknown(planConfiguration types.Object) (types.Object, diag.Diagnostics) {
// 	if !internaltypes.IsDefined(planConfiguration) {
// 		return planConfiguration, nil
// 	}
// 	planConfigurationAttrs := planConfiguration.Attributes()
// 	planConfigurationAttrs["fields_all"] = types.SetUnknown(types.ObjectType{AttrTypes: fieldAttrTypes})
// 	planConfigurationAttrs["tables_all"] = types.ListUnknown(types.ObjectType{AttrTypes: tablesSensitiveFieldsSplitAttrTypes})
// 	return types.ObjectValue(configurationAttrTypes, planConfigurationAttrs)
// }
