CREATE TABLE users(
  id integer not null primary key,
  email varchar not null unique,
  pass varchar not null
);