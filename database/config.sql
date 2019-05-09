
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
  catalog []text,
  details text,
  internal-notes text,
  is-active boolean,
  FOREIGN KEY (winery) REFERENCES winery(name) ON DELETE RESTRICT,
  FOREIGN KEY (territoy, region, nation) REFERENCES territory(name, region, nation) ON DELETE RESTRICT,
);


CREATE TABLE winery (
  winery text NOT NULL PRIMARY KEY,
);


CREATE TABLE origin (
  territory text FOREIGN KEY,
  region text FOREIGN KEY,
  nation text NOT NULL FOREIGN KEY,
);

CREATE TABLE catalog (
  id int PRIMARY KEY NOT NULL,
  name text NOT NULL,
  level integer NOT NULL CHECK (level >= 0 AND level <= 3),
  parent int FOREIGN KEY REFERENCES catalog(id)
  type text[],
  size float[],
  year int[],
  territory text[],
  region text[],
  country text[],
  winery text[],
  storage text[],
  wine integer[],
  FOREIGN KEY (territory, region, country) REFERENCES territory(name, region, country)
);

CREATE TABLE purchase (
  id int PRIMARY KEY NOT NULL,
  wine int NOT NULL references wine(id),
  date timestamp NOT NULL,
  supplier text,
  quantity int NOT NULL,
  cost money NOT NULL
);
