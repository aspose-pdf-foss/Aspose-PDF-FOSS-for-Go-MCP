// SPDX-License-Identifier: MIT

package tools

import (
	"context"
	"encoding/json"
	"fmt"

	pdf "github.com/aspose-pdf-foss/aspose-pdf-foss-for-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type formFieldsParams struct {
	InputPath string `json:"input_path" jsonschema:"path to the input PDF file"`
	Password  string `json:"password,omitempty" jsonschema:"password for an encrypted PDF (optional)"`
}

type formFieldJSON struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Page     int    `json:"page,omitempty"`
	ReadOnly bool   `json:"read_only,omitempty"`
	Required bool   `json:"required,omitempty"`
}

type formFieldsResult struct {
	Count  int             `json:"count"`
	Fields []formFieldJSON `json:"fields"`
}

func handleFormFields(ctx context.Context, req *mcp.CallToolRequest, args formFieldsParams) (*mcp.CallToolResult, any, error) {
	doc, err := openDocument(args.InputPath, args.Password)
	if err != nil {
		return nil, nil, err
	}
	fields := doc.Form().Fields()
	out := formFieldsResult{Count: len(fields), Fields: []formFieldJSON{}}
	for _, f := range fields {
		out.Fields = append(out.Fields, formFieldJSON{
			Name:     f.FullName(),
			Type:     fieldTypeName(pdf.FieldType(f)),
			Value:    f.Value(),
			Page:     f.PageIndex(),
			ReadOnly: f.IsReadOnly(),
			Required: f.IsRequired(),
		})
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal fields: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}, nil, nil
}

func fieldTypeName(t pdf.FormFieldType) string {
	switch t {
	case pdf.FormFieldTypeText:
		return "text"
	case pdf.FormFieldTypeCheckbox:
		return "checkbox"
	case pdf.FormFieldTypeRadioButton:
		return "radio"
	case pdf.FormFieldTypePushButton:
		return "button"
	case pdf.FormFieldTypeComboBox:
		return "combobox"
	case pdf.FormFieldTypeListBox:
		return "listbox"
	case pdf.FormFieldTypePassword:
		return "password"
	case pdf.FormFieldTypeFileSelect:
		return "fileselect"
	case pdf.FormFieldTypeRichText:
		return "richtext"
	case pdf.FormFieldTypeNumber:
		return "number"
	case pdf.FormFieldTypeDate:
		return "date"
	default:
		return "unknown"
	}
}
