CREATE TABLE users (
    user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    creation_date TIMESTAMP DEFAULT now(),
    user_name VARCHAR NOT NULL
);

CREATE TABLE wallets (
    wallet_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID NOT NULL,
    creation_date TIMESTAMP DEFAULT now()
);

CREATE TABLE operations (
    operation_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    wallet_id UUID NOT NULL,
    operation_date TIMESTAMP DEFAULT now(),
    amount FLOAT NOT NULL
);

-- тестовые данные
WITH new_user AS (
    INSERT INTO users(user_name) VALUES('user1') RETURNING user_id
), new_wallet AS (
	INSERT INTO wallets(user_id) SELECT user_id FROM new_user RETURNING wallet_id
)
INSERT INTO operations(wallet_id, amount)
SELECT (SELECT wallet_id FROM new_wallet) AS wallet_id,(random() * 200 - 100)::INT AS amount
FROM generate_series(1, 1000);

WITH new_user AS (
    INSERT INTO users(user_name) VALUES('user2') RETURNING user_id
), new_wallet AS (
	INSERT INTO wallets(user_id) SELECT user_id FROM new_user RETURNING wallet_id
)
INSERT INTO operations(wallet_id, amount)
SELECT (SELECT wallet_id FROM new_wallet) AS wallet_id,(random() * 200 - 100)::INT AS amount
FROM generate_series(1, 1000);

WITH new_user AS (
    INSERT INTO users(user_name) VALUES('user3') RETURNING user_id
), new_wallet AS (
	INSERT INTO wallets(user_id) SELECT user_id FROM new_user RETURNING wallet_id
)
INSERT INTO operations(wallet_id, amount)
SELECT (SELECT wallet_id FROM new_wallet) AS wallet_id,(random() * 200 - 100)::INT AS amount
FROM generate_series(1, 1000);

WITH new_user AS (
    INSERT INTO users(user_name) VALUES('user4') RETURNING user_id
), new_wallet AS (
	INSERT INTO wallets(user_id) SELECT user_id FROM new_user RETURNING wallet_id
)
INSERT INTO operations(wallet_id, amount)
SELECT (SELECT wallet_id FROM new_wallet) AS wallet_id,(random() * 200 - 100)::INT AS amount
FROM generate_series(1, 1000);