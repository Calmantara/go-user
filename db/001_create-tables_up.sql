CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name varchar(255),
  address varchar(255),
  email varchar(255) NOT NULL,
  password varchar(255) 
);

CREATE TABLE photos(
  id SERIAL PRIMARY KEY,
  name varchar(255),
  user_id int,
  CONSTRAINT fk_user_photo FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE creditcards(
  user_id int,
  number varchar(12) PRIMARY KEY,
  cvv varchar(3),
  expired varchar(4),
  type varchar(100),
  name varchar(255),
  CONSTRAINT fk_user_creditcard FOREIGN KEY(user_id) REFERENCES users(id)
);
