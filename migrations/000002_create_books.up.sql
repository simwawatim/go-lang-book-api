CREATE TABLE
    books (
        id BIGSERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        author TEXT NOT NULL,
        isbn TEXT NOT NULL DEFAULT '',
        published_year INT,
        description TEXT NOT NULL DEFAULT '',
        owner_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_books_owner_id ON books (owner_id);