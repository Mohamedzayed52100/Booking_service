
CREATE TABLE Events (
  eventID SERIAL PRIMARY KEY,
  eventName VARCHAR(255) NOT NULL,
  totalSeats INTEGER NOT NULL,
  availableSeats INTEGER NOT NULL CHECK (availableSeats >= 0)
);