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
    scanned BOOLEAN NOT NULL DEFAULT FALSE,
    porn INTEGER NOT NULL DEFAULT -1,
    CONSTRAINT unique_url UNIQUE(url),
    PRIMARY KEY(id)
);

CREATE TABLE visits (
    id SERIAL,
    user_id INTEGER,
    link_id INTEGER,
    PRIMARY KEY(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT unique_user_link UNIQUE(user_id, link_id),
    FOREIGN KEY (link_id) REFERENCES links(id)
);

INSERT INTO users (name, password) values ('test', 'test', 'testghash');
INSERT INTO users (name, password) values ('testadmin', 'testadmin', 'testadminhash');