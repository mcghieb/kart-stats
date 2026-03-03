CREATE TABLE races (
  id smallint NOT NULL PRIMARY KEY, -- this is the heat number
  time timestamp NOT NULL
)

CREATE TABLE karts (
  id smallint NOT NULL PRIMARY KEY -- this is the kart number
)

CREATE TABLE drivers (
  id varchar(15) NOT NULL PRIMARY KEY, -- this is the CustID that is listed on the website
  alias varchar(30) NOT NULL,
  name varchar(30), 
  proskill_rating smallint NOT NULL
)

CREATE TABLE laps (
  id SERIAL PRIMARY KEY, -- find out how postgres deals with uuid's
  race_id smallint NOT NULL REFERENCES races(id), -- FK to races.id
  kart_id smallint NOT NULL REFERENCES karts(id), -- FK to karts.id
  driver_id varchar(15) NOT NULL REFERENCES drivers(id), -- FK to drivers.id
  lap_time numeric(10, 3) NOT NULL
)

CREATE INDEX idx_laps_race ON laps (race_id);
CREATE INDEX idx_laps_kart ON laps (kart_id);
CREATE INDEX idx_laps_driver ON laps (driver_id);

CREATE TABLE race_results (
  race_id smallint NOT NULL REFERENCES races(id), -- FK to races.id
  kart_id smallint NOT NULL REFERENCES karts(id), -- FK to karts.id
  driver_id varchar(15) NOT NULL REFERENCES drivers(id), -- FK to drivers.id
  position smallint NOT NULL,
  penalties smallint NOT NULL,
  best_laptime numeric(10, 3) NOT NULL,
  avg_laptime numeric(10, 3) NOT NULL,
  num_laps smallint NOT NULL,
  gap_from_leader numeric(10, 3) NOT NULL,

  PRIMARY KEY (race_id, kart_id, driver_id),
  UNIQUE (race_id, position),
  UNIQUE (race_id, kart_id)
)

CREATE INDEX idx_race_results_race ON race_results (race_id);
CREATE INDEX idx_race_results_kart ON race_results (kart_id);
CREATE INDEX idx_race_results_driver ON race_results (driver_id);

