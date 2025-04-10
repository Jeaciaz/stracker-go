PRAGMA journal_mode=WAL;

CREATE TABLE IF NOT EXISTS categories (
  id TEXT PRIMARY KEY,
  label TEXT NOT NULL,
  emoji TEXT NOT NULL
);

INSERT OR IGNORE INTO categories (id, label, emoji) VALUES ('groceries', 'Groceries', '🛒'),
('subscriptions', 'Subscriptions', '📦'),
('medical', 'Medical', '🏥'),
('entertainment', 'Entertainment', '🎥'),
('clothing', 'Clothing', '👕'),
('cosmetics', 'Cosmetics', '💄'),
('home', 'Home', '🏠'),
('coffee', 'Coffee & Snacks', '☕'),
('cat', 'Cat', '🐱'),
('misc', 'Misc', '🎲');

CREATE TABLE IF NOT EXISTS spendings (
  id TEXT PRIMARY KEY,
  username TEXT NOT NULL,
  timestamp INTEGER NOT NULL,
  amount REAL NOT NULL,
  category TEXT NOT NULL,
  description TEXT,
  FOREIGN KEY (category) REFERENCES categories(id)
);
