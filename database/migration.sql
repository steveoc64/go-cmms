-- DB Migration Script

-- To Reset all the events and tasks, and clean out all the schedules
TRUNCATE event RESTART IDENTITY;
TRUNCATE task RESTART IDENTITY;
TRUNCATE task_check RESTART IDENTITY;
TRUNCATE task_part RESTART IDENTITY;
TRUNCATE sched_task RESTART IDENTITY;
TRUNCATE sched_task_part RESTART IDENTITY;


-- 2016-05-11  
-- Modify task part records

alter table task add labour_hrs numeric(12,2) not null default 0;
alter table task_part add qty_used numeric(12,2) not null default 0;
alter table part_stock add descr text not null default '';
alter table part_price add descr text not null default '';

-- Capture SMS transmissions

drop table if exists sms_trans;
create table sms_trans (
	id serial not null primary key,
	number_to text not null default '',
	number_used text not null default '',
	user_id int not null default 0,
	message text not null default '',
	date_sent timestamptz not null default localtimestamp,
	ref text not null default '',
	status text not null default '',
	error text not null default ''
);

-- 2016-05-12
-- Modify user to have hourly rate, and seq task IDs by site

alter table users add hourly_rate numeric(12,2) not null default 0;
alter table users add address text not null default '';
alter table users add site_id int not null default 0;
alter table users add notes text not null default '';

-- 2016-05-16
-- Syslog has a more useful fields

alter table user_log 
add channel int not null default 0,
add user_id int not null default 0,
add entity text not null default '',
add entity_id int not null default 0,
add error text not null default '',
add is_update bool not null default false;

-- Parts tree
alter table part add category int not null default 0;
create table category (
	id serial not null primary key,
	parent_id int not null default 0,
	name text not null default '',
	descr text not null default ''
);

create table site_category (
	site_id int not null,
	cat_id int not null
);
create index site_category_idx on site_category (site_id, cat_id);