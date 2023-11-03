DROP TABLE IF EXISTS users;
CREATE TABLE users(
  id         INT UNIQUE,
  name       VARCHAR(128) NOT NULL,
  email      VARCHAR(128) UNIQUE,
  password   VARCHAR(128) NOT NULL,
  tier       INT NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO users
  (id, name, email, password, tier)
VALUES
  (1, "Vivian", "vivian@vivian.com", "vivian123", 5);