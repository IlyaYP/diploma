--- DROP DATABASE IF EXISTS gmart;
--- CREATE DATABASE gmart;

CREATE TABLE IF NOT EXISTS users
(
    login varchar(40) not null,
    password varchar(64) not null,
    primary key (login)
);


CREATE TABLE IF NOT EXISTS orders
(
    num int not null,
	status int not null,
	accrual int,
	uploaded_at timestamp    not null default now(),
	login varchar(64) not null,
    primary key (num),
    foreign key (login) references users (login)
);


---- create above / drop below ----

DROP TABLE orders;
DROP TABLE users;
