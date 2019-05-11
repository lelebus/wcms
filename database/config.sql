
CREATE TABLE winery (
  name text NOT NULL PRIMARY KEY
);

CREATE TABLE origin (
  territory text UNIQUE,
  region text UNIQUE,
  nation text NOT NULL UNIQUE
);

CREATE TABLE wine (
  id SERIAL PRIMARY KEY NOT NULL UNIQUE,
  storage_area text NOT NULL,
  type text NOT NULL,
  size float NOT NULL,
  name text NOT NULL,
  winery_cellar text NOT NULL references winery(name),
  year int NOT NULL,
  territory text references origin(territory),
  region text references origin(region),
  nation text references origin(nation) NOT NULL,
  price money NOT NULL,
  catalog text,
  details text,
  internal_notes text,
  is_active boolean
);


CREATE TABLE catalog (
  id int PRIMARY KEY NOT NULL,
  name text NOT NULL,
  level int NOT NULL CHECK (level >= 0 AND level <= 3),
  parent int references catalog(id),
  type text,
  size float,
  year int,
  territory text references origin(territory),
  region text references origin(region),
  nation text references origin(nation),
  winery text,
  storage text,
  wine int references wine(id)
);

CREATE TABLE purchase (
  id int PRIMARY KEY NOT NULL,
  wine int NOT NULL references wine(id),
  date timestamp NOT NULL,
  supplier text,
  quantity int NOT NULL,
  cost money NOT NULL
);
