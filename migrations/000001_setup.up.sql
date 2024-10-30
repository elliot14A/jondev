-- up.sql
create table if not exists updated_at_trigger (
    id integer primary key,
    table_name text not null
);

create trigger if not exists auto_update_timestamp
after update on updated_at_trigger
begin
    update updated_at_trigger 
    set updated_at = datetime('now')
    where id = new.id;
end;

