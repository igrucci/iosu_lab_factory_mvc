create table Provider(
                         id bigserial not null ,
                         provider_name varchar not null unique,
   constraint pk_provider primary key (id)
);

create table TechnicalProcess(
                        ordering_id bigint not null ,
                        forging_id bigint not null ,
                        date_end timestamptz default now() not null,
    constraint fk_technical_process_ordering foreign key (ordering_id) references Ordering(id) on delete cascade ,
    constraint fk_technical_process_forging foreign key (forging_id) references Forging(id)

);

create table Forging(
                        id bigserial not null ,
                        forging_type_id bigint not null ,
                        casting_id bigint,
                        count_forging bigint not null ,
                        place_id bigint,
    constraint pk_forging primary key (id),
    constraint fk_forging_forging_type foreign key (forging_type_id) references ForgingType(id),
    constraint fk_forging_casting foreign key (casting_id) references Casting(id),
    constraint fk_forging_place foreign key (place_id) references Place(id)
);

create table ForgingType(
                        id bigserial not null ,
                        type_name varchar not null unique ,
    constraint pk_forging_type primary key (id)

);

create table Casting(
                        id bigserial not null ,
                        casting_type_id bigint not null ,
                        count_casting bigint not null ,
                        material_id bigint not null ,
                        count_material bigint not null ,
                        place_id bigint not null,
    constraint pk_casting primary key (id),
    constraint fk_casting_casting_type foreign key (casting_type_id) references CastingType(id),
    constraint fk_casting_material foreign key (material_id) references Material(id),
    constraint fk_casting_place foreign key (place_id) references Place(id)
);

create table CastingType(
                        id bigserial not null ,
                        type_name varchar not null unique,
    constraint pk_casting_type primary key (id)
);

create table Equipment(
                        id bigserial,
                        equipment_type_id bigint,
                        place_id bigint,
                        is_available varchar default 'available',
    constraint pk_equipment primary key (id),
    constraint fk_equipment_equipment_type foreign key (equipment_type_id) references EquipmentType(id),
    constraint fk_equipment_place foreign key (place_id) references Place(id)
);

create table EquipmentType(
                        id bigserial,
                        type_name varchar unique,
    constraint pk_equipment_type primary key (id)
);

create table Worker (
                        id bigserial not null ,
                        worker_name varchar not null ,
                        position_id bigint not null,
                        place_id bigint not null,
                        is_available varchar not null default 'available',
    constraint pk_worker primary key (id),
    constraint fk_worker_position foreign key (position_id) references Position(id),
    constraint fk_worker_place foreign key (place_id) references Place(id)
);

create table Position (
                        id bigserial not null,
                        position_name varchar not null unique,
    constraint pk_position primary key (id)
);

create table Place (
                        id bigserial not null ,
                        place_name varchar not null unique,
    constraint pk_place primary key (id)
);

create table Detail (
                        id bigserial not null ,
                        detail_name varchar not null unique ,
                        material_id bigint  ,
                        count_material bigint not null,
    constraint pk_detail primary key (id),
    constraint fk_detail_material foreign key (material_id) references Material(id)
);

create table MaterialType(
                        id bigserial not null ,
                        type_name varchar unique not null,
    constraint pk_material_type primary key (id)
);

create table Material (
                        id bigserial not null,
                        material_type_id bigint unique ,
                        material_count bigint not null check ( material_count >= 0) default 0,
                        provider_id bigint not null default 0,
    constraint pk_material primary key (id),
    constraint fk_material_material_type foreign key (material_type_id) references MaterialType(id),
    constraint fk_material_provider foreign key (provider_id) references Provider(id)
);

create table Customer (
                        id     bigserial,
                        customer_name   varchar     not null unique ,
    constraint pk_customer primary key (id)
);


create table Ordering(
                         id    bigserial,
                         detail_id  bigint ,
                         count_detail bigint not null,
                         customer_id  bigint     not null ,
                         date_registration timestamptz default now() not null,
                         status varchar not null default 'available',
    constraint pk_ordering primary key (id),
    constraint fk_ordering_detail foreign key (detail_id) references Detail(id),
    constraint fk_ordering_customer foreign key (customer_id) references Customer(id)
);
create table OrderingArchive(
                          id bigserial,
                          ordering_id    bigint,
                          detail_id  bigint ,
                          count_detail bigint not null,
                          customer_id  bigint     not null ,
                          date_registration timestamptz not null,
                          date_adding timestamptz default now() not null,
                          status varchar not null default 'completed'
);

