-- Drop the foreign key constraint on blogs.main_post_id
ALTER TABLE blogs DROP CONSTRAINT IF EXISTS fk_main_post_id;

-- Drop the blog_posts table
DROP TABLE IF EXISTS blog_posts;

-- Drop the blogs table
DROP TABLE IF EXISTS blogs;
