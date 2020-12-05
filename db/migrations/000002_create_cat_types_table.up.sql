create table if not exists cat_types
(
    id         serial    not null primary key,
    name       text      not null,
    avatar_url text      not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp null
);

