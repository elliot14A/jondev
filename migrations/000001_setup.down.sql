-- down.sql
drop trigger if exists auto_update_timestamp;
drop table if exists updated_at_trigger;
