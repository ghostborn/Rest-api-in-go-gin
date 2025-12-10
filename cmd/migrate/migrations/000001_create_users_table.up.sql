create table if not exists users
(
    id       integer primary key autoincrement,
    email    text not null unique,
    name     text not null,
    password text not null
)