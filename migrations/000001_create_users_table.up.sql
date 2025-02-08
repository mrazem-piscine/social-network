CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nickname TEXT UNIQUE,
    email TEXT UNIQUE,
    password TEXT,
    age INTEGER,
    gender TEXT,
    first_name TEXT,
    last_name TEXT
);
