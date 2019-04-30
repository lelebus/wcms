CREATE TABLE wine (
  id integer NOT NULL PRIMARY KEY,
  type text NOT NULL CHECK(wineType = 'red' OR wineType = 'white' OR wineType = 'sparkling'),
  storage ???,
  name text,
  year integer NOT NULL CHECK(year > 1899 AND year < currentYear),
  winery text,
  size float NOT NULL CHECK(size = 0.375 OR size = 0.75 OR size = 1 OR size = 1.5 OR size = 3
							OR size = 4.5 OR size = 6 OR size = 8 OR size = 9 OR size = 12  
							OR size = 15 OR size = 17 OR size = 20 OR size = 30),
  details ???,
  region text,
  country varchar(3),
  price money NOT NULL CHECK(price > 0),
  catalog text,
  isActive ???,
  FOREIGN KEY (winery) REFERENCES wine(name) ON DELETE RESTRICT,
  FOREIGN KEY (region, nation) REFERENCES territory(region, nation) ON DELETE RESTRICT,
  FOREIGN KEY (catalog) REFERENCES catalog(name) ON DELETE RESTRICT
);


CREATE TABLE winery (
  name text NOT NULL PRIMARY KEY,
  FOREIGN KEY (name) REFERENCES wine(winery) ON DELETE RESTRICT
);


CREATE TABLE territory (
  region text NOT NULL PRIMARY KEY,
  nation text NOT NULL PRIMARY KEY,
  FOREIGN KEY (region, nation) REFERENCES wine( region, nation) ON DELETE RESTRICT


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