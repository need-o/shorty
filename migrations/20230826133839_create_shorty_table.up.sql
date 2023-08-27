create table shorty
(
    id         text     not null constraint shorty_pk primary key,
    url        text     not null,
    created_at datetime not null,
    updated_at datetime not null
);

create unique index shorty_id_uindex
    on shorty (id);
