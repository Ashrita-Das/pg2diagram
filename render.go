package main

import (
	"fmt"
	"strings"
)

func generateMermaidDiagram(extractedSchema DatabaseSchema) string {
	var mermaidStringBuilder strings.Builder

	
	mermaidStylingDirective := "%%{init: {'theme': 'dark', 'themeVariables': { 'primaryColor': '#1e1e2e', 'primaryBorderColor': '#ff5e00', 'primaryTextColor': '#cdd6f4', 'lineColor': '#ff5e00'}}}%%\n"
	mermaidStringBuilder.WriteString(mermaidStylingDirective)

	mermaidStringBuilder.WriteString("erDiagram\n")

	for _, currentRelation := range extractedSchema.Relations {
		relationLine := fmt.Sprintf("    %s ||--o{ %s : \"%s references %s\"\n",
			currentRelation.SourceTableName,
			currentRelation.TargetTableName,
			currentRelation.SourceColumnName,
			currentRelation.TargetColumnName,
		)
		mermaidStringBuilder.WriteString(relationLine)
	}

	mermaidStringBuilder.WriteString("\n")

	for _, currentTable := range extractedSchema.Tables {
		tableHeader := fmt.Sprintf("    %s {\n", currentTable.TableName)
		mermaidStringBuilder.WriteString(tableHeader)

		for _, currentColumn := range currentTable.Columns {
			sanitizedDataType := strings.ReplaceAll(currentColumn.DataType, " ", "_")

			columnLine := fmt.Sprintf("        %s %s\n", sanitizedDataType, currentColumn.ColumnName)
			mermaidStringBuilder.WriteString(columnLine)
		}

		mermaidStringBuilder.WriteString("    }\n\n")
	}

	return mermaidStringBuilder.String()
}