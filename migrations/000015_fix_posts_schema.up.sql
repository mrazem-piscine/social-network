-- Step 1: Turn off foreign key constraints temporarily
PRAGMA foreign_keys=off;

-- Step 2: Drop `posts_new` if it already exists
DROP TABLE IF EXISTS posts_new;

-- Step 3: Create a new `posts` table with the correct structure
CREATE TABLE posts_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    privacy TEXT CHECK(privacy IN ('public', 'followers-only', 'private')) DEFAULT 'public',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Step 4: Migrate existing data (Use `privacy`, not `category`)
INSERT INTO posts_new (id, user_id, content, privacy, created_at)
SELECT id, user_id, content, privacy, created_at FROM posts;

-- Step 5: Drop the old `posts` table
DROP TABLE posts;

-- Step 6: Rename the new table to `posts`
ALTER TABLE posts_new RENAME TO posts;

-- Step 7: Turn foreign keys back on
PRAGMA foreign_keys=on;
