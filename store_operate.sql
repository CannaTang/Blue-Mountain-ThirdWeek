#创建数据库及内部表
drop database store;
create database store;
use store;
create table depository(
    depository_id varchar(20) key,
    depository_capacity integer
);
create table clothes(
    clothing_id varchar(20) key,
    clothing_size varchar(20),
    clothing_price integer,
    clothing_type varchar(20)
);
create table supplier(
    supplier_id varchar(20) key,
    supplier_name varchar(20)
);
create table supply_condition(
    clothing_id varchar(20),
    supplier_id varchar(20),
    clothing_level varchar(20),
    primary key (clothing_id, supplier_id)
);

#插入
insert into depository values ('001', 10), ('002', 20);
insert into clothes values ('Blue01', 'S', 99, 'dress'), ('Red01', 'S', 100, 'dress');
insert into supplier values ('001', 'tencent'), ('002', 'wangyi');
insert into supply_condition values ('Blue01', '001', 'S+'), ('Red01', '002', 'D');


#查询
#(1)查询服装尺码为'S'且销售价格在100以下的服装信息。
select *
from clothes
where clothing_size = 'S' && clothes.clothing_price < 110;

#(2)查询仓库容量最大的仓库信息。
select *
from depository
where depository_capacity = (
    select max(depository_capacity)
    from depository
    );

#（3）查询A类服装的库存总量。
#假如仓库编号和供应商编号一样， 否则按规定的4个表无法查询
select depository_capacity
from depository
where depository_id = any(
    select supplier_id
    from supply_condition
    where clothing_id = any(
        select clothing_id
        from clothes
        where clothing_type = 'dress'
        )
);

#(4) 查询服装编码以‘A’开始开头的服装。
#这里我只放了Blue01和Red01, 所以搜一下以B开头的
select *
from clothes
where clothing_id like 'B%';

#（5）查询服装质量等级有不合格的供应商信息。
select *
from supplier
where supplier_id = (
    select supplier_id
    from supply_condition
    where clothing_level <= 'D'
    );

#更新
#(1)把服装尺寸为'S'的服装的销售价格均在原来基础上提高10%。
update clothes
set clothing_price = clothing_price*1.1;

#删除
#(1)删除所有服装质量等级不合格的供应情况。
delete from supply_condition
where clothing_level <= 'D';

