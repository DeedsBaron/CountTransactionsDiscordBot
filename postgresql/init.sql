CREATE DATABASE transactions;
\c transactions;
CREATE TABLE SellBuyinfo (
    UserID NUMERIC PRIMARY KEY,
    UserNickname TEXT NOT NULL,
    Selled INT,
    Bought INT
);
