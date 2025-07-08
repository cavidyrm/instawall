-- Drop everything in reverse order of creation to avoid dependency errors.
DROP TRIGGER IF EXISTS update_pages_updated_at ON pages;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TABLE IF EXISTS page_categories;
DROP TABLE IF EXISTS pages;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS update_updated_at_column();