CREATE TABLE Employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    coins INTEGER
);

CREATE TABLE Merch (
    id SERIAL PRIMARY KEY,
    name VARCHAR(10) NOT NULL,
    price INTEGER
);

CREATE TABLE Buy (
    id_employees INTEGER,
    id_merch INTEGER,
    FOREIGN KEY (id_employees) REFERENCES Employees (id),
    FOREIGN KEY (id_merch) REFERENCES Merch (id)
);

CREATE TABLE Transaction (
    id_src INTEGER, 
    id_dest INTEGER, 
    coins INTEGER, 
    FOREIGN KEY (id_src) REFERENCES Employees(id), 
    FOREIGN KEY (id_dest) REFERENCES Employees(id)
);

INSERT INTO Merch (name, price)
VALUES
('t-shirt', 80),
('cup', 20),
('book', 50),
('pen', 10),
('powerbank', 200),
('hoody', 300),
('umbrella', 200),
('socks', 10),
('wallet', 50),
('pink-hoody', 500);

