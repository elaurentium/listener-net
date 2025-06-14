-- migrations/001_init.sql

CREATE TABLE USERS (
    ID UUID PRIMARY KEY,
    IP VARCHAR(50) UNIQUE NOT NULL,
    NAME VARCHAR(50) NOT NULL,
    LAST_SEEN TIMESTAMP,
    DISPOSITIVE VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE BLACKLIST (
    ID UUID PRIMARY KEY,
    IP VARCHAR(50) UNIQUE NOT NULL,
    NAME VARCHAR(50) NOT NULL,
    LAST_SEEN TIMESTAMP,
    DISPOSITIVE VARCHAR(255) UNIQUE NOT NULL
);

