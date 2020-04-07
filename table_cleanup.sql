delete from data a using data b where a.time_ran < b.time_ran and a.cases = b.cases and a.country = b.country;

delete from info a using info b where a.time_ran < b.time_ran and a.api_id = b.api_id and a.latitude = b.latitude and a.longitude = b.longitude;
