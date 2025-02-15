PRAGMA foreign_keys=off;

-- Rename the old table
ALTER TABLE group_members RENAME TO group_members_old;

-- Create the new table with the correct status options
CREATE TABLE group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status TEXT CHECK(status IN ('pending', 'member', 'admin')) DEFAULT 'pending',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(group_id, user_id)
);

-- Copy data from the old table to the new one
INSERT INTO group_members (id, group_id, user_id, status, joined_at)
SELECT id, group_id, user_id, 
       CASE status 
           WHEN 'approved' THEN 'member' 
           ELSE status 
       END, 
       joined_at FROM group_members_old;

-- Drop the old table
DROP TABLE group_members_old;

PRAGMA foreign_keys=on;
