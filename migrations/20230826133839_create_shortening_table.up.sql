create table shortenings
(
    id         text     not null constraint shortening_pk primary key,
    url        text     not null,
    visits     integer,
    created_at datetime not null,
    updated_at datetime not null
);

create unique index shortening_id_uindex
    on shortenings (id);
