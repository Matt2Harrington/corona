create table if not exists data
(
  id                     uuid not null primary key,
  country                varchar,
	cases                  integer,
	cases_today            integer,
	deaths                 integer,
	deaths_today           integer,
	recovered              integer,
	active                 integer,
	critical               integer,
	cases_per_one_million  float,
	deaths_per_one_million float,
	updated                timestamp,
	  time_ran             timestamp default '2020-04-06 20:59:23'::timestamp without time zone
);
