package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)


func extractDatabaseSchema(databaseConnection *pgx.Conn, backgroundContext context.Context) (DatabaseSchema, error) {
	var finalExtractedSchema DatabaseSchema

	tablesAndColumnsQuery := `
		SELECT table_name, column_name, data_type 
		FROM information_schema.columns 
		WHERE table_schema = 'public' 
		ORDER BY table_name, ordinal_position;
	`
	
	columnRows, queryError := databaseConnection.Query(backgroundContext, tablesAndColumnsQuery)
	if queryError != nil {
		return finalExtractedSchema, fmt.Errorf("failed to query columns: %v", queryError)
	}
	defer columnRows.Close()

	tableColumnsGroupedMap := make(map[string][]ColumnDefinition)

	for columnRows.Next() {
		var currentTableName string
		var currentColumnName string
		var currentDataType string

		scanError := columnRows.Scan(&currentTableName, &currentColumnName, &currentDataType)
		if scanError != nil {
			return finalExtractedSchema, fmt.Errorf("failed to scan column row: %v", scanError)
		}

		newColumnDefinition := ColumnDefinition{
			ColumnName:   currentColumnName,
			DataType:     currentDataType,
			IsPrimaryKey: false,
		}

		tableColumnsGroupedMap[currentTableName] = append(tableColumnsGroupedMap[currentTableName], newColumnDefinition)
	}

	for mapTableName, mapColumnsList := range tableColumnsGroupedMap {
		newTableMetadata := TableMetadata{
			TableName: mapTableName,
			Columns:   mapColumnsList,
		}
		finalExtractedSchema.Tables = append(finalExtractedSchema.Tables, newTableMetadata)
	}

	
	foreignKeysQuery := `
		SELECT
			tc.table_name AS source_table_name,
			kcu.column_name AS source_column_name,
			ccu.table_name AS target_table_name,
			ccu.column_name AS target_column_name
		FROM 
			information_schema.table_constraints AS tc 
			JOIN information_schema.key_column_usage AS kcu
			  ON tc.constraint_name = kcu.constraint_name
			  AND tc.table_schema = kcu.table_schema
			JOIN information_schema.constraint_column_usage AS ccu
			  ON ccu.constraint_name = tc.constraint_name
			  AND ccu.table_schema = tc.table_schema
		WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_schema='public';
	`
	
	foreignKeyRows, foreignKeyQueryError := databaseConnection.Query(backgroundContext, foreignKeysQuery)
	if foreignKeyQueryError != nil {
		return finalExtractedSchema, fmt.Errorf("failed to query foreign keys: %v", foreignKeyQueryError)
	}
	defer foreignKeyRows.Close()

	for foreignKeyRows.Next() {
		var currentSourceTableName string
		var currentSourceColumnName string
		var currentTargetTableName string
		var currentTargetColumnName string

		scanError := foreignKeyRows.Scan(&currentSourceTableName, &currentSourceColumnName, &currentTargetTableName, &currentTargetColumnName)
		if scanError != nil {
			return finalExtractedSchema, fmt.Errorf("failed to scan foreign key row: %v", scanError)
		}

		newForeignKeyRelation := ForeignKeyRelation{
			SourceTableName:  currentSourceTableName,
			SourceColumnName: currentSourceColumnName,
			TargetTableName:  currentTargetTableName,
			TargetColumnName: currentTargetColumnName,
		}
		finalExtractedSchema.Relations = append(finalExtractedSchema.Relations, newForeignKeyRelation)
	}

	return finalExtractedSchema, nil
}