create table if not exists tbl_user
(
    id bigint not null
        constraint tbl_user_pk
            primary key,
    first_name varchar(128),
    last_name varchar(128),
    username varchar(128),
    is_bot boolean
);

alter table tbl_user owner to tgbot;

create unique index tbl_user_id_uindex
	on tbl_user (id);
