CREATE TABLE wine (
  id integer NOT NULL PRIMARY KEY SERIAL UNIQUE,
  type text NOT NULL,
  storage text NOT NULL,
  name text NOT NULL,
  year integer NOT NULL,
  winery text NOT NULL,
  size float NOT NULL,
  details text,
  territory text,
  region text,
  country text NOT NULL,
  price money NOT NULL,
  catalog text PRIMARY KEY,
  isActive boolean,
  FOREIGN KEY (winery) REFERENCES winery(name) ON DELETE RESTRICT,
  FOREIGN KEY (territoy, region, nation) REFERENCES territory(name, region, nation) ON DELETE RESTRICT,
  FOREIGN KEY (catalog) REFERENCES catalog(id) ON DELETE RESTRICT
);


CREATE TABLE winery (
  name text NOT NULL PRIMARY KEY,
  FOREIGN KEY (name) REFERENCES wine(winery) ON DELETE RESTRICT
);


CREATE TABLE territory (
  name text NOT NULL PRIMARY KEY,
  region text NOT NULL PRIMARY KEY,
  nation text NOT NULL PRIMARY KEY,
  FOREIGN KEY (name, region, nation) REFERENCES wine(territory, region, nation) ON DELETE RESTRICT


CREATE TABLE catalog (
  name text NOT NULL PRIMARY KEY,
  level integer NOT NULL CHECK(level >= 0 AND level <= 3),
  FOREIGN KEY (name) REFERENCES wine(catalog)
);


CREATE TABLE purchase (
  wine text NOT NULL PRIMARY KEY,
  date timestamp NOT NULL PRIMARY KEY,
  supplier text,
  quantity integer NOT NULL,
  cost money NOT NULL,
  FOREIGN KEY (wine) REFERENCES wine(id) ON DELETE RESTRICT
);
