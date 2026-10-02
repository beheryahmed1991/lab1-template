\connect persons

CREATE TABLE IF NOT EXISTS persons (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age INTEGER,
    address VARCHAR(255),
    work VARCHAR(255)
);

GRANT ALL PRIVILEGES ON TABLE persons TO program;

GRANT USAGE, SELECT ON SEQUENCE persons_id_seq TO program;