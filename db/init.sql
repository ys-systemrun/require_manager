-- Initial database setup for the Requirement Manager.
-- The database itself is created by the postgres image via POSTGRES_DB.
-- Table schemas are managed by the backend through GORM AutoMigrate,
-- so this file only ensures required extensions and encoding.

SET client_encoding = 'UTF8';

-- Case-insensitive text and trigram search helpers (optional but useful
-- for the terminology dictionary / full-text-ish lookups).
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
