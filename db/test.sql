DROP TABLE IF EXISTS users;
CREATE TABLE users(
  id         INT AUTO_INCREMENT NOT NULL,
  name      VARCHAR(128) NOT NULL,
  email     VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO users
  (name, email)
VALUES
  ('Lexi Gray', 'email1@gmail.com'),
  ('Jeru Steps', 'email2@gmail.com');
