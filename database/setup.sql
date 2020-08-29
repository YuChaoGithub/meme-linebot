-- This is a denormalized design for easier insertion.
-- Since I don't have time to develop an admin page for myself (yet),
-- a denormalized design makes it much easier for me to insert new memes
-- directly through SQL.

-- Note that there can be alias names for the same meme (i.e. memes with
-- the same URL), that's why it's denormalized.

CREATE TABLE memes(
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) UNIQUE,
    url VARCHAR(128)
);
