create table if not exists users (
    id int auto_increment,
    username varchar(40) not null,
    password varchar(255) not null,
    email varchar(255) not null,
    unique(username, email),
    primary key (id)
);

create table if not exists role_details (
    id int auto_increment,
    detail varchar(40) not null,
    unique(detail),
    primary key (id)
);

insert into role_details(detail) values("user");
insert into role_details(detail) values("admin");



create table if not exists user_roles (
    id int auto_increment,
    user_id int not null,
    role_id int not null,
    primary key (id),
    foreign key (role_id) references role_details(id),
    foreign key (user_id) references users(id)
);

insert into users(username, password, email)
values ("foo", "fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9", "foo@example.com");

insert into users(username, password, email)
values ("bao", "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3", "bao@example.com");

insert into user_roles(user_id, role_id) value (1, 1);
insert into user_roles(user_id, role_id) value (2, 1);
insert into user_roles(user_id, role_id) value (2, 2);