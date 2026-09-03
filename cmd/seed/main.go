package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env if present
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/mariadiezma?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("Conectando a PostgreSQL (%s)...\n", dbURL)
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Error al conectar con PostgreSQL: %v\n\nTIP: Si usas Docker, asegúrate de haber iniciado los contenedores con 'make docker-up' o 'docker compose up -d'.", err)
	}
	defer conn.Close(ctx)

	// 1. Ensure Schema is applied
	schemaPath := "migrations/000001_init_schema.up.sql"
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Error leyendo archivo de esquema %s: %v", schemaPath, err)
	}

	fmt.Println("Aplicando/verificando esquema de base de datos...")
	if _, err := conn.Exec(ctx, string(schemaSQL)); err != nil {
		log.Fatalf("Error ejecutando esquema SQL: %v", err)
	}
	fmt.Println("✓ Esquema base verificado correctamente.")

	// 2. Apply Seed Data
	seedPath := "migrations/000002_seed_data.sql"
	seedSQL, err := os.ReadFile(seedPath)
	if err != nil {
		log.Fatalf("Error leyendo archivo de datos de prueba %s: %v", seedPath, err)
	}

	fmt.Println("Inyectando datos de prueba (usuarios, colecciones, vestidos, citas)...")
	if _, err := conn.Exec(ctx, string(seedSQL)); err != nil {
		log.Fatalf("Error inyectando datos de prueba: %v", err)
	}

	fmt.Println("✓ ¡Datos de prueba inyectados con éxito en PostgreSQL!")
}
