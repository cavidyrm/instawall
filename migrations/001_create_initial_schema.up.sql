CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';


-- Step 2: Create the 'users' table and its related objects.
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       mobile_number VARCHAR(20) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       name VARCHAR(100) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       role VARCHAR(50) NOT NULL DEFAULT 'user',
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_users_mobile_number ON users(mobile_number);
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();


-- Step 3: Create the 'categories' table.
CREATE TABLE categories (
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            title VARCHAR(100) UNIQUE NOT NULL,
                            description TEXT,
                            image_url VARCHAR(255),
                            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Step 4: Create the 'pages' table and its related objects.
CREATE TABLE pages (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       image_url VARCHAR(255),
                       link VARCHAR(255),
                       has_issue BOOLEAN NOT NULL DEFAULT FALSE,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_pages_user_id ON pages(user_id);
CREATE TRIGGER update_pages_updated_at
    BEFORE UPDATE ON pages
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();


-- Step 5: Create the 'page_categories' join table and its indexes.
CREATE TABLE page_categories (
                                 page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
                                 category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
                                 PRIMARY KEY (page_id, category_id)
);
CREATE INDEX idx_page_categories_page_id ON page_categories(page_id);
CREATE INDEX idx_page_categories_category_id ON page_categories(category_id);
