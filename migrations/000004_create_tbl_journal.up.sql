create table if not exists tbl_journal
(
    id     serial not null
        constraint tbl_journal_pkey
            primary key,
    record json
);

alter table tbl_journal
    owner to tgbot;
