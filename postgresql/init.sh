#!/bin/bash
while :
do
	sleep 1
done
CREATE DATABASE urlshortener;
\c urlshortener
create user deedsbaron with encrypted password '0809';
CREATE TABLE urls (
	id bigserial PRIMARY KEY,
	longURL VARCHAR ( 50 ) UNIQUE NOT NULL,
	shortURL VARCHAR ( 50 ) NOT NULL
);
INSERT INTO urls (id, longURL, shortURL) VALUES (1238129389123, 'asdsadasd', '666666');
INSERT INTO urls(id, longURL, shortURL)

SELECT longurl, shorturl
FROM urls
WHERE shortURL = 'http://localhost:8080/_h2kndDy6J';