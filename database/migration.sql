-- DB Migration Script


-- 2016-05-11  
-- Modify task part records

alter table task_part add qty_used numeric(12,2) not null default 0;
alter table part_stock add descr text not null default '';
alter table part_price add descr text not null default '';
