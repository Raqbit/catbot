create table if not exists cat_activities
(
    id            serial    not null primary key,
    description   text      not null,
    only_for_type int       null,
    created_at    timestamp not null,
    updated_at    timestamp not null,
    deleted_at    timestamp null,

    --- On deletion of cat types, delete all related cat activities
    constraint cat_activities_cat_types_fk foreign key (only_for_type) references cat_types (id) on delete cascade
);

