create table if not exists info
(
  id                    uuid not null primary key,
  api_id                integer,
  latitude              float,
  longitude             float,
  data_id               uuid not null,
  updated               timestamp,
  time_ran              timestamp default '2020-04-06 20:59:23'::timestamp without time zone
);
