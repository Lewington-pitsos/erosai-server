DROP TABLE users;
DROP TABLE links;
DROP TABLE visits;

CREATE TABLE users (
    id SERIAL,
    name VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    token VARCHAR(100) NOT NULL,
    CONSTRAINT unique_name UNIQUE(name),
    PRIMARY KEY(id)
);

CREATE TABLE links (
    id SERIAL,
    url VARCHAR(2000) NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE visits (
    id SERIAL,
    user_id INTEGER,
    link_id INTEGER,
    PRIMARY KEY(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (link_id) REFERENCES links(id)
);

INSERT INTO users (name, password) values ('test', 'test', 'testghash');
INSERT INTO users (name, password) values ('testadmin', 'testadmin', 'testadminhash');