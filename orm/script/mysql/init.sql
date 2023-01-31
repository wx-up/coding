create database if not exists `integration_test`;
create table if not exists `integration_test`.`simple_structs`(
                                                                 `id` bigint auto_increment,
                                                                 `name` varchar(50),
    primary key (`id`)
    );
