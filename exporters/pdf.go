package exporters

import (
	"fmt"
	"reflect"

	"github.com/jung-kurt/gofpdf"
)

func ExportToPDF(data any, filename string, selectedFields []string, title string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 20, 10)
	pdf.AddPage()

	sliceValue := reflect.ValueOf(data)
	if sliceValue.Kind() != reflect.Slice || sliceValue.Len() == 0 {
		return fmt.Errorf("data must be a non-empty slice of structs")
	}

	firstItem := sliceValue.Index(0)

	allFields := getStructFields(firstItem)

	filteredFields := []string{}
	for _, field := range selectedFields {
		if contains(allFields, field) {
			filteredFields = append(filteredFields, field)
		}
	}

	if len(filteredFields) == 0 {
		return fmt.Errorf("no valid fields selected for export")
	}

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 12, title, "0", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 12)
	for _, field := range filteredFields {
		pdf.CellFormat(60, 10, field, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 12)
	for i := range sliceValue.Len() {
		row := sliceValue.Index(i)
		rowValues := getStructValues(row, filteredFields)

		for _, value := range rowValues {
			pdf.CellFormat(60, 10, fmt.Sprintf("%v", value), "1", 0, "C", false, 0, "")
		}

		pdf.Ln(-1)
	}

	if err := pdf.OutputFileAndClose(filename); err != nil {
		return fmt.Errorf("failed to save PDF: %v", err)
	}

	return nil
}

func ExportSingleToPDF(data any, filename string, selectedFields []string, title string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "B", 14)
	pdf.AddPage()

	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Struct {
		return fmt.Errorf("data must be a struct")
	}

	allFields := getStructFields(dataValue)

	filteredFields := []string{}
	for _, field := range selectedFields {
		if contains(allFields, field) {
			filteredFields = append(filteredFields, field)
		}
	}

	if len(filteredFields) == 0 {
		return fmt.Errorf("no valid fields selected for export")
	}

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190.0, 10.0, title, "0", 1, "C", false, 0, "")
	pdf.Ln(5)

	colWidthLeft := 60.0
	colWidthRight := 120.0
	rowHeight := 10.0

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(colWidthLeft, rowHeight, "Field", "1", 0, "C", false, 0, "")
	pdf.CellFormat(colWidthRight, rowHeight, "Value", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 12)
	rowValues := getStructValues(dataValue, filteredFields)

	for i, field := range filteredFields {
		pdf.CellFormat(colWidthLeft, rowHeight, field, "1", 0, "L", false, 0, "")
		pdf.CellFormat(colWidthRight, rowHeight, fmt.Sprintf("%v", rowValues[i]), "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	if err := pdf.OutputFileAndClose(filename); err != nil {
		return fmt.Errorf("failed to save PDF: %v", err)
	}

	return nil
}
