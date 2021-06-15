create table if not exists tbl_invite
(
    channel_id bigint not null
        constraint tbl_invite___ch_fk
            references tbl_channel,
    from_id bigint not null
        constraint tbl_invite___fr_fk
            references tbl_user,
    member_id bigint not null
        constraint tbl_invite___mm_fk
            references tbl_user,
    constraint tbl_invite_pk
        primary key (channel_id, from_id, member_id)
);

alter table tbl_invite owner to tgbot;

