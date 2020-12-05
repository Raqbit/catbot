create table if not exists users
(
    id         serial                    not null primary key,
    discord_id text                      not null,
    money      bigint    default 0       not null,
    last_daily timestamp default 'epoch' not null,
    created_at timestamp                 not null,
    updated_at timestamp                 not null,
    deleted_at timestamp                 null
);

-- Prevent duplicate discord users
create unique index if not exists users_discord_id_uindex
    on users (discord_id);
