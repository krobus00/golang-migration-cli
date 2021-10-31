-- +goose Up
-- +goose StatementBegin
CREATE TABLE Persons (
    PersonID int,
    LastName varchar(255),
    FirstName varchar(255),
    Address varchar(255),
    City varchar(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP Table Persons;
-- +goose StatementEnd
