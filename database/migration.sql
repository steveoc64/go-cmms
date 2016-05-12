-- DB Migration Script


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
