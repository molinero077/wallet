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
    operation_type_id UUID NOT NULL,
    operation_date TIMESTAMP DEFAULT now(),
    amount FLOAT NOT NULL
);

CREATE TABLE operation_type (
    operation_type_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    operation_name VARCHAR NOT NULL
);

-- тестовые данные
WITH new_user AS (
    INSERT INTO users(user_name) VALUES('user1') RETURNING user_id
)
INSERT INTO wallets(user_id) SELECT user_id FROM new_user;

WITH new_user AS (
    INSERT INTO users(user_name) VALUES('user2') RETURNING user_id
)
INSERT INTO wallets(user_id) SELECT user_id FROM new_user;

WITH new_user AS (
    INSERT INTO users(user_name) VALUES('user3') RETURNING user_id
)
INSERT INTO wallets(user_id) SELECT user_id FROM new_user;

WITH new_user AS (
    INSERT INTO users(user_name) VALUES('user4') RETURNING user_id
)
INSERT INTO wallets(user_id) SELECT user_id FROM new_user;


INSERT INTO operation_type(operation_name) VALUES
    ('DEPOSIT'),
    ('WITHDRAW');


INSERT INTO operations(wallet_id, operation_type_id, amount) 
    SELECT wallets.wallet_id, operation_type.operation_type_id, 1000 
    FROM wallets, operation_type 
    WHERE operation_type.operation_name='DEPOSIT';

