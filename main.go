package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	databaseConnectionString := flag.String("db", "", "PostgreSQL connection string (e.g., postgres://user:pass@localhost:5432/mydb)")
	flag.Parse()

	if *databaseConnectionString == "" {
		fmt.Println("Error: Missing database connection string.")
		fmt.Println("Usage: go run main.go -db=\"postgres://user:pass@localhost:5432/mydb\"")
		os.Exit(1)
	}

	backgroundContext := context.Background()

	databaseConnection, connectionError := pgx.Connect(backgroundContext, *databaseConnectionString)
	
	if connectionError != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to connect to database: %v\n", connectionError)
		os.Exit(1)
	}
	
	defer databaseConnection.Close(backgroundContext)

	pingError := databaseConnection.Ping(backgroundContext)
	
	if pingError != nil {
		fmt.Fprintf(os.Stderr, "Error: Database connected, but ping failed: %v\n", pingError)
		os.Exit(1)
	}

	extractedDatabaseSchema, extractionError := extractDatabaseSchema(databaseConnection, backgroundContext)
	if extractionError != nil {
		fmt.Fprintf(os.Stderr, "Error extracting schema: %v\n", extractionError)
		os.Exit(1)
	}

	finalMermaidOutput := generateMermaidDiagram(extractedDatabaseSchema)

	fmt.Println(finalMermaidOutput)
}