// Copyright © 2026 Ping Identity Corporation

package modifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func Float64SemanticEqual() planmodifier.Float64 {
	return float64SemanticEqualModifier{}
}

type float64SemanticEqualModifier struct {
}

func (m float64SemanticEqualModifier) Description(_ context.Context) string {
	return ""
}

func (m float64SemanticEqualModifier) MarkdownDescription(_ context.Context) string {
	return ""
}

func (m float64SemanticEqualModifier) PlanModifyFloat64(ctx context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() ||
		req.StateValue.IsNull() || req.PlanValue.IsNull() ||
		req.StateValue.IsUnknown() || req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.ValueFloat64() == req.PlanValue.ValueFloat64() {
		resp.PlanValue = req.StateValue
	}
}
