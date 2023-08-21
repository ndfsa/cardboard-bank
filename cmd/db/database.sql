DROP TABLE IF EXISTS user_service CASCADE;
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    id BIGSERIAL,
    fullname VARCHAR(300) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(60) NOT NULL,
    PRIMARY KEY (id)
);

DROP TYPE CURR CASCADE;
CREATE TYPE curr AS ENUM ('USD', 'CAD', 'JPY', 'NOK');

DROP TABLE IF EXISTS services CASCADE;
CREATE TABLE services (
    id BIGSERIAL,
    type SMALLINT,
    state SMALLINT,
    currency CURR,
    init_balance NUMERIC(20, 2),
    balance NUMERIC(20, 2),
    PRIMARY KEY (id)
);

CREATE TABLE user_service (
    user_id BIGINT,
    service_id BIGINT,
    PRIMARY KEY (user_id, service_id),
    FOREIGN KEY (user_id) REFERENCES users ON DELETE CASCADE,
    FOREIGN KEY (service_id) REFERENCES services ON DELETE CASCADE
);

-- create root user with default password
INSERT INTO users (fullname, username, password) VALUES (
    'root user',
    'root',
    '$2a$06$DZxsYD5zF5NI/ugKmMmZw.7/hehCmlCpzDOuPutYFmwIlyT37SDGy'
);


-- function to create a user, checks for null and encrypts the password
DROP FUNCTION IF EXISTS CREATE_USER;
CREATE FUNCTION CREATE_USER(
    _fullname VARCHAR(300),
    _username VARCHAR(100),
    _password VARCHAR(72)
) RETURNS BIGINT AS
$$
DECLARE
    res BIGINT;
BEGIN
    IF _fullname IS NULL OR _fullname = '' THEN
        RAISE EXCEPTION 'fullname cannot be NULL or empty.';
    END IF;

    IF _username IS NULL OR _username = '' THEN
        RAISE EXCEPTION 'username cannot be NULL or empty.';
    END IF;

    IF _password IS NULL OR _password = '' THEN
        RAISE EXCEPTION 'password cannot be NULL or empty.';
    END IF;

    INSERT INTO users(fullname, username, password) VALUES (
        _fullname,
        _username,
        crypt(_password, gen_salt('bf')))
    RETURNING id INTO res;

    RETURN res;
END
$$
LANGUAGE 'plpgsql';


-- function to update user, only username and fullname can be updated
DROP PROCEDURE IF EXISTS UPDATE_USER;
CREATE PROCEDURE UPDATE_USER(
    _id BIGINT,
    _fullname VARCHAR(300),
    _username VARCHAR(100)
) AS
$$
BEGIN
    UPDATE users SET
        fullname = COALESCE(NULLIF(_fullname, ''), fullname),
        username = COALESCE(NULLIF(_username, ''), username)
    WHERE id = _id;
END
$$
LANGUAGE 'plpgsql';


-- authentication function, returns user ID for tokens
DROP FUNCTION IF EXISTS AUTHENTICATE_USER;
CREATE FUNCTION AUTHENTICATE_USER(
    _username VARCHAR,
    _password VARCHAR
) RETURNS BIGINT AS
$$
DECLARE
    _id BIGINT;
    _auth BOOLEAN;
BEGIN
    IF _username IS NULL OR _username = '' THEN
        RAISE EXCEPTION 'username cannot be NULL or empty.';
    END IF;

    IF _password IS NULL OR _password = '' THEN
        RAISE EXCEPTION 'password cannot be NULL or empty.';
    END IF;

    SELECT (password = crypt(_password, password)) , id
    INTO _auth, _id
    FROM users
    WHERE username = _username;

    IF _auth THEN
        RETURN _id;
    END IF;

    RAISE EXCEPTION 'user authentication unsuccessful.';
END
$$
LANGUAGE 'plpgsql';


DROP FUNCTION IF EXISTS GET_USER_SERVICES;
CREATE FUNCTION GET_USER_SERVICES(
    _id BIGINT
) RETURNS SETOF SERVICES AS
$$
BEGIN
    RETURN QUERY SELECT s.id, s.type, s.state, s.currency, s.init_balance, s.balance
    FROM users u
    JOIN user_service us ON u.id = us.user_id
    JOIN services s ON s.id = us.service_id
    WHERE u.id = _id;
END
$$
LANGUAGE 'plpgsql';

DROP FUNCTION IF EXISTS CREATE_SERVICE;
CREATE FUNCTION CREATE_SERVICE(
    _user_id BIGINT,
    _type SMALLINT,
    _currency CURR,
    _init_balacne NUMERIC(20, 2)
) RETURNS BIGINT AS
$$
DECLARE
    res BIGINT;
BEGIN
    INSERT INTO services (type, state, currency, init_balance, balance)
    VALUES (_type, 1, _currency, _init_balacne, 0)
    RETURNING id INTO res;

    INSERT INTO user_service (user_id, service_id)
    VALUES (_user_id, res);

    RETURN res;
END
$$
LANGUAGE 'plpgsql';
