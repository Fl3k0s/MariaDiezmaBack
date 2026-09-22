-- MariaDiezmaBack Initial Schema Migration

-- Enable UUID extension if available
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==========================================================
-- 1. USERS: Backoffice administrators & operators
-- ==========================================================
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'viewer',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(LOWER(email));

-- ==========================================================
-- 2. REQUESTS: REST requests / Inquiries / Appointments
-- ==========================================================
CREATE TABLE IF NOT EXISTS requests (
    id VARCHAR(36) PRIMARY KEY,
    type VARCHAR(50) NOT NULL DEFAULT 'general',
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    priority VARCHAR(50) NOT NULL DEFAULT 'medium',
    sender_name VARCHAR(255) NOT NULL,
    sender_email VARCHAR(255) NOT NULL,
    sender_phone VARCHAR(50),
    subject VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb,
    internal_note TEXT DEFAULT '',
    assigned_to VARCHAR(36) REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_requests_status ON requests(status);
CREATE INDEX IF NOT EXISTS idx_requests_priority ON requests(priority);
CREATE INDEX IF NOT EXISTS idx_requests_type ON requests(type);
CREATE INDEX IF NOT EXISTS idx_requests_created_at ON requests(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_requests_assigned_to ON requests(assigned_to);

-- ==========================================================
-- 3. COLLECTIONS: Collections catalogue for public web frontend
-- ==========================================================
CREATE TABLE IF NOT EXISTS collections (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    image_path VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_collections_name ON collections(LOWER(name));
CREATE INDEX IF NOT EXISTS idx_collections_created_at ON collections(created_at DESC);

-- ==========================================================
-- 4. DRESSES: Dresses catalogue and details for public web frontend
-- ==========================================================
CREATE TABLE IF NOT EXISTS dresses (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    collection VARCHAR(255) NOT NULL,
    image_path VARCHAR(500) NOT NULL,
    image1_path VARCHAR(500) NOT NULL DEFAULT '',
    image2_path VARCHAR(500) NOT NULL DEFAULT '',
    image3_path VARCHAR(500) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_dresses_collection ON dresses(LOWER(collection));
CREATE INDEX IF NOT EXISTS idx_dresses_name_col ON dresses(LOWER(name), LOWER(collection));
CREATE INDEX IF NOT EXISTS idx_dresses_created_at ON dresses(created_at DESC);

-- ==========================================================
-- 5. PRESS_ARTICLES: Press articles / media appearances for web frontend
-- ==========================================================
CREATE TABLE IF NOT EXISTS press_articles (
    id VARCHAR(36) PRIMARY KEY,
    magazine_name VARCHAR(255) NOT NULL,
    publication_date VARCHAR(50) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    article_url VARCHAR(1000) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_press_articles_magazine ON press_articles(LOWER(magazine_name));
CREATE INDEX IF NOT EXISTS idx_press_articles_created_at ON press_articles(created_at ASC);
