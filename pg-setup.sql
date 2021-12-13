CREATE TABLE IF NOT EXISTS clients
 (
     id         INT GENERATED ALWAYS AS IDENTITY,
     name       VARCHAR(255) NOT NULL,
     created_at TIMESTAMP    NOT NULL,
     PRIMARY KEY (id)
 );

DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'currency') THEN
            CREATE TYPE currency AS ENUM ('MXN', 'USD', 'COP');
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS accounts
(
    id         INT GENERATED ALWAYS AS IDENTITY,
    client_id  INT       NOT NULL,
    currency   CURRENCY  NOT NULL,
    balance    FLOAT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_client
        FOREIGN KEY (client_id)
            REFERENCES clients (id)
);

DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transaction_type') THEN
            CREATE TYPE TRANSACTION_TYPE AS ENUM ('deposit','withdraw','transfer');
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS transactions
(
    id          INT GENERATED ALWAYS AS IDENTITY,
    receiver_id INT,
    sender_id   INT,
    amount      FLOAT            NOT NULL,
    type        TRANSACTION_TYPE NOT NULL,
    created_at  TIMESTAMP        NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_receiver
        FOREIGN KEY (receiver_id)
            REFERENCES accounts (id),
    CONSTRAINT fk_sender
        FOREIGN KEY (sender_id)
            REFERENCES accounts (id)
);