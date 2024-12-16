create table restaurants (
    uuid binary(16) primary key,
    name varchar(255),
    contact_phone varchar(255),
    description text,
    cover_image_url varchar(255),
    address varchar(255),
    primary key (uuid)
);

create table restaurant_tables (
    uuid binary(16),
    restaurant_uuid binary(16),
    number int,
    primary key (uuid),
    foreign key (restaurant_uuid) references restaurants (uuid)
);

create table table_reservations (
    uuid binary(16),
    client_phone varchar(255),
    start_date date,
    end_date date,
    restaurant_uuid binary(16),
    table_uuid binary(16),
    primary key (uuid),
    foreign key (restaurant_uuid) references restaurants (uuid),
    foreign key (table_uuid) references restaurant_tables (uuid)
);
