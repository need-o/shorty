create table visits
(
    shorty_id  text     not null,
    referer    text     not null,
    user_ip    text     not null,
    user_agent text     not null,
    created_at datetime not null,
    updated_at datetime not null
);
