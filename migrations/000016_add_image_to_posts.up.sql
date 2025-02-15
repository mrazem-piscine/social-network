-- Add an `image` column to store image file paths
ALTER TABLE posts ADD COLUMN image TEXT DEFAULT NULL;
