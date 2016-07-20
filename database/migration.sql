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

-- 2016-06-03 
-- Fix up machine layout for Chinderrah and Connecticut
delete from site_layout where site_id=8;
delete from site_layout where site_id=9;

insert into site_layout (site_id, seq, machine_id, span) values
(8,1,26,12),
(8,2,22,12),
(8,3,23,12),
(8,4,25,12),
(8,5,24,12),
(9,1,40,12),
(9,2,41,12),
(9,3,39,12),
(9,4,38,12),
(9,5,37,12),
(9,6,42,12),
(9,7,43,12);

-- 2016-06-13
-- MachineTypes database

drop table if exists machine_type;
create table machine_type (
	id serial not null primary key,
	name text not null default '',
	electrical bool default true,
	hydraulic bool default true,
	pnuematic bool default true,
	lube bool default true,
	printer bool default true,
	console bool default true,
	uncoiler bool default true,
	rollbed bool  default true
);

insert into machine_type (id,name) 
values (1,'Bracket'),(2,'Stud'),(3,'Chord'),(4,'Plate'),(5,'Web'),(6,'Floor'),(7,'Valley'),
(8,'Top Hat 22'),(9,'Top Hat 40'),(10,'Mill'),(11,'Wall'),(12,'Lathe'),(13,'Surface Grinder'),
(14,'Guillotine'),(15,'Folder');

alter table machine add pnuematic text not null default 'Running';

drop table if exists machine_type_tool;
create table machine_type_tool (
	machine_id int not null,
	position int not null default 0,
	name text not null default ''
);

create index machine_type_tool_idx on machine_type_tool (machine_id, position);

insert into machine_type_tool (machine_id, position, name)
values (1,1,'Guillo'),
(2,1,'Brick Tie'),(2,2,'Service Hole #1'),(2,3,'Quad Dimple'),(2,4,'Service Hole #2'),(2,5,'Single Dimple & Rib'),(2,6,'Curl #1'),(2,7,'Guillo'),(2,8,'Curl #2'),
(3,1,'Down Dimple'),(3,2,'Tie Down Slot'),(3,3,'Up Dimple'),(3,4,'Half Notch'),(3,5,'Full Notch'),(3,6,'Right Angle Guillo'),(3,7,'Straight Guillo'),(3,8,'Left Angle Guillo'),
(4,1,'Single Dimple Square'),(4,2,'Service Hole / Curl'),(4,3,'Tie Down Slot'),(4,4,'Notch'),(4,5,'Nogging'),(4,6,'Guillo'),
(5,1,'Pierce Location'),(5,2,'Crush'),(5,3,'Guillo'),
(6,1,'H-Cut'),(6,2,'Notch'),(6,3,'Pier Slot'),(6,4,'Tie Down Slot'),(6,5,'Fold Rib'),(6,6,'Dimple Bearer'),(6,7,'Dimple Joist'),(6,8,'Service Hole'),(6,9,'Swage'),(6,10,'Up Dimple'),(6,11,'Down Dimple'),(6,12,'Guillo Bearer'),(6,13,'Guillo Joist'),
(7,1,'Guillo Valley'),(7,2,'Guillo Lintel'),
(8,1,'Guillo'),
(9,1,'Guillo'),
(11,1,'Service Hole'),(11,2,'Quad Dimple'),(11,3,'Single Dimple Crush'),(11,4,'Curl'),(11,5,'Guillo'),(11,6,'Single Dimple Square'),(11,7,'Notch'),(11,8,'Tie Down Slot'),(11,9,'Nogging');

alter table machine add machine_type int not null default 0;

update machine set machine_type = 1 where name like 'Bracket%';
update machine set machine_type = 2 where name like 'Stud%';
update machine set machine_type = 3 where name like 'Chord%';
update machine set machine_type = 4 where name like 'Plate%';
update machine set machine_type = 5 where name like 'Web%';
update machine set machine_type = 6 where name like 'Floor%';
update machine set machine_type = 7 where name like 'Valley%';
update machine set machine_type = 8 where name like 'Top Hat 22%';
update machine set machine_type = 9 where name like 'Top Hat 40%';
update machine set machine_type = 10 where name like 'Mill%';
update machine set machine_type = 11 where name like 'Wall%';
update machine set machine_type = 12 where name like 'Lathe%';
update machine set machine_type = 13 where name like 'Surface Grinder%';
update machine set machine_type = 14 where name like 'Guillotine%';
update machine set machine_type = 15 where name like 'Folder%';

-- 2016 06 29
-- Add conveyor to machine

alter table machine add conveyor text not null default 'Running';
alter table machine_type add conveyor bool default true;

-- 2016 06 29
-- Toggle SMS on / off per user

alter table users add use_mobile bool default false;

-- 2016 06 30
-- Extend parts list functionality, and add photo uploader

alter table category add stock_code text not null default '';

-- 2016 07 05
-- Photo upload tester

create table phototest (
	id serial not null primary key,
	name text,
	photo text,
	preview text,
	thumbnail text
);

-- 2016 07 05
-- Extend users table and add migration tracking

create table migration (
	id serial not null primary key,
	name text,
	date timestamptz not null default localtimestamp
);

insert into migration (name) values ('Init Migration Database');

alter table users add is_tech bool not null default false;
update users set is_tech = true where role = 'Technician' or username = 'Shane Voigt';
update users set use_mobile = true where sms <> '';

insert into migration (name) values ('Extend user info');

-- 2016 07 06
-- Update task to track whether its been read or not

alter table task add is_read bool default false;
alter table task add read_date timestamptz;

insert into migration (name) values ('Extend task to track whether user has read it or not');

-- Add photo to stoppage event

alter table event add photo text not null default '';
alter table event add photo_preview text not null default '';
alter table event add photo_thumbnail text not null default '';

insert into migration (name) values ('Add photo on event');

-- User can allocate flag

alter table users add can_allocate bool not null default false;
insert into migration (name) values ('User can allocate flag');

-- 2016 07 13
-- additional photos on event


create table phototest2 (
	id serial not null primary key,
	name text,
	photo text[],
	preview text[],
	thumbnail text[]
);

-- 2016 07 16
-- more photos

alter table machine_type add photo text not null default '';
alter table machine_type add photo_preview text not null default '';
alter table machine_type add photo_thumbnail text not null default '';

insert into migration (name) values ('Add more photo fields');

-- 2016 07 17
-- link category to machine type

alter table category add machine_type int not null default 0;
alter table category add machine_tool int not null default 0;

insert into migration (name) values ('Add machine type link on category');

-- Attach photos to the actual task

alter table task add photo1 text not null default '';
alter table task add photo2 text not null default '';
alter table task add photo3 text not null default '';

alter table task add preview1 text not null default '';
alter table task add preview2 text not null default '';
alter table task add preview3 text not null default '';

alter table task add thumb1 text not null default '';
alter table task add thumb2 text not null default '';
alter table task add thumb3 text not null default '';

insert into migration (name) values ('More photos on the task');

-- 2016 07 20
-- Add encoder and strip guide to the machine and machine type, and all the diagrams

alter table machine add encoder text not null default 'Running';
alter table machine add strip_guide text not null default 'Running';
alter table machine_type add encoder bool default true;
alter table machine_type add strip_guide bool default true;

insert into migration (name) values ('Encoder and strip guide');