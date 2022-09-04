create table if not exists test
(
    id  varchar(64) not null,
    url varchar     not null
    );

create unique index if not exists index_id
    on test (id);

alter table test
    add constraint key_id
        primary key (id);
