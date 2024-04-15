CREATE TABLE IF NOT EXISTS cars
(
    id      UUID PRIMARY KEY,
    regNum TEXT  NOT NULL,
    mark    TEXT  NOT NULL,
    model   TEXT  NOT NULL,
    year    INT DEFAULT 0,
    owner   JSONB NOT NULL
);

ALTER TABLE cars
    ADD CONSTRAINT regNum_unique UNIQUE (regNum);
