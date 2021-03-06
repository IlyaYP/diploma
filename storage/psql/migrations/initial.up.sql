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
    num bigint not null,
	status int not null,
	accrual int not null default 0,
	uploaded_at timestamp with time zone   not null default now(),
	login varchar(64) not null,
    primary key (num),
    foreign key (login) references users (login)
);


---- create above / drop below ----

DROP TABLE orders;
DROP TABLE users;
