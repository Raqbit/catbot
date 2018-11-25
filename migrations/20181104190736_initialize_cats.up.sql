create table if not exists cats
(
  id          serial primary key        not null,
  owner_id    integer                   not null,
  ck_id       integer                   not null,
  name        varchar(100)              not null,
  pronoun     varchar(10)               not null,
  hunger      int default 100           not null,
  last_fed    timestamp default 'epoch' not null,
  away        boolean default false     not null,

  --- Cascading deletion to prevent stray cats
  constraint cats_users_id_fk foreign key (owner_id) references users (id) on delete cascade
);