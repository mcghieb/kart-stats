CREATE TABLE races (
  id smallint PRIMARY KEY, -- this is the heat number
  time timestamp,
)

CREATE TABLE karts (
  id smallint PRIMARY KEY, -- this is the kart number
)

CREATE TABLE drivers (
  id varchar(15) PRIMARY KEY, -- this is the CustID that is listed on the website
  name varchar(30),
  alias varchar(30), 
  proskill_rating smallint,
)

CREATE TABLE laps (
  id uuid PRIMARY KEY, -- find out how postgres deals with uuid's
  race_id smallint, -- FK to races.id
  kart_id smallint, -- FK to karts.id
  driver_id varchar(15), -- FK to drivers.id
  lap_time number(10, 3),
)

CREATE TABLE race_results (
  race_id smallint, -- FK to races.id
  kart_id smallint, -- FK to karts.id
  driver_id varchar(15), -- FK to drivers.id
  position smallint,
  penalties smallint,
  best_laptime number(10, 3),
  avg_laptime number(10, 3),
  num_laps smallint,
  gap_from_leader number(10, 3),

  -- PK is composite key of (race_id, kart_id, driver_id, position)
)


