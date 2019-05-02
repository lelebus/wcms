CREATE TABLE wine (
  id integer NOT NULL PRIMARY KEY SERIAL UNIQUE,
  storage-area text NOT NULL,
  type text NOT NULL,
  size float NOT NULL,
  name text NOT NULL,
  winery text NOT NULL,
  year integer NOT NULL,
  territory text,
  region text,
  country text NOT NULL,
  price money NOT NULL,
  catalog text PRIMARY KEY,
  details text,
  internal-notes text,
  is-active boolean,
  FOREIGN KEY (winery) REFERENCES winery(name) ON DELETE RESTRICT,
  FOREIGN KEY (territoy, region, nation) REFERENCES territory(name, region, nation) ON DELETE RESTRICT,
  FOREIGN KEY (catalog) REFERENCES catalog(id) ON DELETE RESTRICT
);


CREATE TABLE winery (
  name text NOT NULL PRIMARY KEY,
);


CREATE TABLE origin (
  territory text FOREIGN KEY,
  region text FOREIGN KEY,
  nation text NOT NULL FOREIGN KEY,
);

CREATE TABLE catalog (
  id int NOT NULL PRIMARY KEY,
  name text NOT NULL,
  level integer NOT NULL CHECK (level >= 0 AND level <= 3),
  parent int FOREIGN KEY REFERENCES catalog(id)
  type text,
  size float,
  year int,
  territory text,
  region text,
  country text,
  winery text,
  storage text,
  wine int,
  FOREIGN KEY (wine) REFERENCES wine(id),
  FOREIGN KEY (territory, region, country) REFERENCES territory(name, region, country)
);

CREATE TABLE purchase (
  id int NOT NULL PRIMARY KEY,
  wine int NOT NULL PRIMARY KEY,
  date timestamp NOT NULL PRIMARY KEY,
  supplier text,
  quantity integer NOT NULL,
  cost money NOT NULL,
  FOREIGN KEY (wine) REFERENCES wine(id) ON DELETE RESTRICT
);
