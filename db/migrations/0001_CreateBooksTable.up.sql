create table books (
    id serial not null primary key,
    title character varying(100) not null,
    author character varying(100) not null,
    pages integer
);
