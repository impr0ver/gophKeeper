CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
                        user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        login VARCHAR(255),
                        password VARCHAR(255)
);

CREATE TABLE data (
                       record_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       user_id VARCHAR(256),
                       record_type INTEGER,
                       keyhint VARCHAR(256),
                       metadata VARCHAR(256),
                       crypted_data VARCHAR(256)
);