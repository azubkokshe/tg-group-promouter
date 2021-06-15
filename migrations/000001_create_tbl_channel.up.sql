create table if not exists tbl_channel
(
    id bigint not null,
    title varchar not null
);

create unique index tbl_channel_id_uindex
	on tbl_channel (id);

alter table tbl_channel
    add constraint tbl_channel_pk
        primary key (id);