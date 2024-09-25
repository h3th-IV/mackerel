CREATE TABLE user (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    user_name VARCHAR(100),
    password VARCHAR(255) NOT NULL, -- 
    ip_address VARCHAR(45) UNIQUE
);

CREATE TABLE geoLocation (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    city VARCHAR(100),
    country VARCHAR(100),
    ip_address VARCHAR(45) NOT NULL UNIQUE,
    region VARCHAR(100),
    lat_long VARCHAR(50), -- latitude and longitude
    organization VARCHAR(100), -- Network Provider Organization
    timezone VARCHAR(100),
    FOREIGN KEY (user_id) REFERENCES User(id) ON DELETE CASCADE
);


export MC_DBNAME=mackerel
export MC_PASSWORD=z3us
export MC_USER=zeus
export MC_HOST=localhost
export MC_PORT=3306
export MAILER_REGION=
export MAILER_ACCESS_ID_KEY=
export MAILER_SECRET_ACCESS_KEY=
export MAILER_SENDER=