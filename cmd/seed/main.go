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

	fmt.Println("Inyectando datos de prueba (usuarios, colecciones, vestidos, citas, artículos de prensa)...")
	if _, err := conn.Exec(ctx, string(seedSQL)); err != nil {
		log.Fatalf("Error inyectando datos de prueba: %v", err)
	}

	// 3. Apply Secretos Collection Migration
	secretosPath := "migrations/000003_add_secretos_collection.sql"
	if secretosSQL, err := os.ReadFile(secretosPath); err == nil {
		fmt.Println("Aplicando migración colección Secretos (000003_add_secretos_collection.sql)...")
		if _, err := conn.Exec(ctx, string(secretosSQL)); err != nil {
			log.Fatalf("Error ejecutando migración de Secretos: %v", err)
		}
		fmt.Println("✓ Colección Secretos inyectada correctamente.")
	}

	// 4. Apply Holidays Table and Initial Madrid Holidays Migration
	holidaysPath := "migrations/000004_create_holidays_table.sql"
	if holidaysSQL, err := os.ReadFile(holidaysPath); err == nil {
		fmt.Println("Aplicando migración festivos de Madrid (000004_create_holidays_table.sql)...")
		if _, err := conn.Exec(ctx, string(holidaysSQL)); err != nil {
			log.Fatalf("Error ejecutando migración de festivos: %v", err)
		}
		fmt.Println("✓ Festivos de Madrid verificados/inyectados correctamente.")
	}

	fmt.Println("✓ ¡Datos de prueba inyectados con éxito en PostgreSQL!")
}
