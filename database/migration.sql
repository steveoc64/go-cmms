-- DB Migration Script


-- 2016-05-11  
-- Modify task part records

alter table task_part add qty_used numeric(12,2) not null default 0;