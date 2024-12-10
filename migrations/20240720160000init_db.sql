create table image (
    id integer not null auto_increment,
    image_url text,
    file_name text,
    source varchar(128),
    nsfw boolean not null default false,
    image_type_id integer,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,

    primary key (id)
);

create table image_type (
    id integer not null auto_increment,
    name text,

    primary key (id)
);

insert into image_type (name) values ('normal');
