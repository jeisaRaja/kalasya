-- Step 1: Create blogs table without the foreign key on main_post_id
CREATE TABLE IF NOT EXISTS blogs (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  subdomain TEXT NOT NULL UNIQUE,
  nav TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Step 2: Create blog_posts table with foreign key on blog_id
CREATE TABLE IF NOT EXISTS blog_posts (
  id SERIAL PRIMARY KEY,
  slug TEXT UNIQUE,
  blog_id INT NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  published BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Step 3: Add main_post_id column to blogs
ALTER TABLE blogs ADD COLUMN main_post_id INT UNIQUE;

-- Step 4: Add the foreign key constraint on main_post_id in blogs table
ALTER TABLE blogs
ADD CONSTRAINT fk_main_post_id FOREIGN KEY (main_post_id) REFERENCES blog_posts(id) ON DELETE SET NULL;
