create table if not exists cats
(
    id                  serial primary key        not null,
    type_id             int                       not null,
    current_activity_id int                       not null,
    owner_id            int                       not null,
    hunger              int       default 100     not null,
    wants_ham           bool      default false   not null,
    last_fed            timestamp default 'epoch' not null,
    last_ham            timestamp default 'epoch' not null,
    created_at          timestamp                 not null,
    updated_at          timestamp                 not null,
    deleted_at          timestamp                 null,

    --- Cascading deletion to prevent stray cats
    constraint cats_users_fk foreign key (owner_id) references users (id) on delete cascade,
    --- Restricting deletion of a cat type which still has cats
    constraint cats_cat_types_fk foreign key (type_id) references cat_types (id) on delete restrict,
    --- Restricting deletion of a cat activity which still has cats
    constraint cats_cat_activities_fk foreign key (current_activity_id) references cat_activities (id) on delete restrict
);

-- Users can only have a single cat for now
create unique index if not exists cats_owner_id_uindex
    on cats (owner_id);

