create table if not exists `versions` (
    `id` int(10) unsigned not null PRIMARY KEY auto_increment,
    `name` varchar(128) not null default '',
    `created` timestamp default current_timestamp
)
engine = innodb
auto_increment = 1
character set utf8
collate utf8_general_ci;
