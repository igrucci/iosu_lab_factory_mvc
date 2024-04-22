insert into Provider(provider_name) values ('MetalWorks Inc');
insert into Provider(provider_name) values ('TechMetal Inc');
insert into Provider(provider_name) values ('SteelForge Inc');

insert into ForgingType(type_name) values ('cold forging');
insert into ForgingType(type_name) values ('open die forging');
insert into ForgingType(type_name) values ('upset forging');

insert into CastingType(type_name) values ('sand casting');
insert into CastingType(type_name) values ('vacuum casting');

insert into EquipmentType(type_name) values ('hydraulic press');
insert into EquipmentType(type_name) values ('lathe');

insert into Place(place_name) values ('casting department');
insert into Place(place_name) values ('forging department');
insert into Place(place_name) values ('office');

insert into Equipment(equipment_type_id, place_id) values (1, 1);
insert into Equipment(equipment_type_id, place_id) values (2, 2);
insert into Equipment(equipment_type_id, place_id) values (1,1);
insert into Equipment(equipment_type_id, place_id) values (2,2);

insert into Position(position_name) values ('engineer');
insert into Position(position_name) values ('manager');

insert into Worker(worker_name, position_id, place_id) values ('Misha I S', 1, 1);
insert into Worker(worker_name, position_id, place_id) values ('Pasha A I', 1, 2);
insert into Worker(worker_name, position_id, place_id) values ('Kostya F D', 1, 1);
insert into Worker(worker_name, position_id, place_id) values ('Anton D D', 2, 3);


insert into Detail(detail_name, material_id, count_material) values ('brake rotor', 1, 100);
insert into Detail(detail_name, material_id, count_material) values ('brake disc', 2, 200);
insert into Detail(detail_name, material_id, count_material) values ('impeller', 3, 150);

insert into MaterialType(type_name) values ('steel');
insert into MaterialType(type_name) values ('aluminum');
insert into MaterialType(type_name) values ('cast iron');

insert into Material(material_type_id, material_count, provider_id) values (1, 10000, 1);
insert into Material(material_type_id, material_count, provider_id) values (2, 15000, 2);
insert into Material(material_type_id, material_count, provider_id) values (3, 20000, 3);






SELECT
    detail_name,
    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 1 THEN count_detail ELSE 0 END) AS "January",
    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 2 THEN count_detail ELSE 0 END) AS "February",
    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 3 THEN count_detail ELSE 0 END) AS "March",
    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 4 THEN count_detail ELSE 0 END) AS "April",
    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 5 THEN count_detail ELSE 0 END) AS "May",
    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 6 THEN count_detail ELSE 0 END) AS "June",
    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 7 THEN count_detail ELSE 0 END) AS "July",
    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 8 THEN count_detail ELSE 0 END) AS "August",
    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 9 THEN count_detail ELSE 0 END) AS "September",
    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 10 THEN count_detail ELSE 0 END) AS "October",
    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 11 THEN count_detail ELSE 0 END) AS "November",
    SUM(CASE WHEN EXTRACT(MONTH FROM date_registration) = 12 THEN count_detail ELSE 0 END) AS "December"
FROM
    detail left join ordering on detail.id = ordering.detail_id
GROUP BY
    detail_name;

SELECT
    o.id AS ordering_id,
    d.detail_name,
    o.count_detail,
    c.customer_name,
    f.count_forging,
    ft.type_name AS forging_type,
    ca.count_casting,
    ct.type_name AS casting_type,
    m.material_count,
    o.date_registration,
    tp.date_end AS date_end,
    o.status
FROM
    Ordering o
        JOIN
    Detail d ON o.detail_id = d.id
        JOIN
    Customer c ON o.customer_id = c.id
        LEFT JOIN
    TechnicalProcess tp ON o.id = tp.ordering_id
        LEFT JOIN
    Forging f ON tp.forging_id = f.id
        LEFT JOIN
    ForgingType ft ON f.forging_type_id = ft.id
        LEFT JOIN
    Casting ca ON f.casting_id = ca.id
        LEFT JOIN
    CastingType ct ON ca.casting_type_id = ct.id
        LEFT JOIN
    Material m ON d.material_id = m.id;


SELECT id, type_name FROM ForgingType
UNION
SELECT id, type_name FROM CastingType;

