CREATE TABLE IF NOT EXISTS Cars (
    ID SERIAL PRIMARY KEY,
    CarIdentifier VARCHAR (127) NOT NULL UNIQUE,
    Name VARCHAR (127) NULL,
    DateOfManufacture DATE NULL,
    Total INT NOT NULL,
    TotalInUse INT NOT NULL
)