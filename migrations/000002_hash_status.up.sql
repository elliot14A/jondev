-- up.sql

create table if not exists hash_status (
    id text primary key,
    is_generated boolean not null default false,
    generated_at timestamp,
    last_verified_at timestamp,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

insert into updated_at_trigger (table_name) 
values ('hash_status');

insert or ignore into hash_status (
    id, 
    is_generated
) values (
    '67b4ed39-b6b9-4957-9e3a-0938f2ac0ebd',  -- fixed uuid for our singleton record
    false
);
