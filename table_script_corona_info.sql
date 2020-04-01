create table if not exists info
(
  id                    uuid not null primary key,
  api_id                integer,
  latitude              float,
  longitude             float,
  data_id               uuid not null,
);
