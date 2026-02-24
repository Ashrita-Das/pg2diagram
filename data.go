package main

type DatabaseSchema struct {
	Tables    []TableMetadata
	Relations []ForeignKeyRelation
}

type TableMetadata struct {
	TableName string
	Columns   []ColumnDefinition
}

type ColumnDefinition struct {
	ColumnName   string
	DataType     string
	IsPrimaryKey bool
}

type ForeignKeyRelation struct {
	SourceTableName  string
	SourceColumnName string
	TargetTableName  string
	TargetColumnName string
}