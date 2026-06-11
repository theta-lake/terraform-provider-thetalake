package swrvrule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type inputSourcesValidator struct{}

func (v inputSourcesValidator) Description(context.Context) string {
	return "input_sources entries with type set to integration must include id"
}

func (v inputSourcesValidator) MarkdownDescription(context.Context) string {
	return "`input_sources` entries with `type` set to `integration` must include `id`"
}

func (v inputSourcesValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var inputSources []swrvRuleInputSourceModel
	resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &inputSources, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for index, inputSource := range inputSources {
		if inputSource.Type.IsNull() || inputSource.Type.IsUnknown() {
			continue
		}

		if inputSource.Type.ValueString() != "integration" {
			continue
		}

		if !inputSource.Id.IsNull() && !inputSource.Id.IsUnknown() {
			continue
		}

		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			req.Path.AtListIndex(index).AtName("id"),
			"Missing integration input source id",
			fmt.Sprintf("input_sources[%d].id must be set when input_sources[%d].type is \"integration\".", index, index),
		))
	}
}

func validateInputSources() validator.List {
	return inputSourcesValidator{}
}

var _ validator.List = inputSourcesValidator{}
var _ = types.List{}
