create table if not exists shortener
(
    id  varchar(64) not null,
    url varchar     not null
    );

create index if not exists index_id
    on shortener (id);

alter table shortener
    add constraint key_id
        primary key (id);

