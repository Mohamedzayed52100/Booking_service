package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Event represents an event with its details
type Event struct {
	EventID        int    `json:"eventID"`
	EventName      string `json:"eventName"`
	TotalSeats     int    `json:"totalSeats"`
	AvailableSeats int    `json:"availableSeats"`
}

// Booking represents a user booking for an event
type Booking struct {
	BookingID   int    `json:"bookingID"`
	EventID     int    `json:"eventID"`
	UserID      int    `json:"userID"`
	BookingTime string `json:"bookingTime"`
}

var db *sql.DB // Database connection

// Connect to the database (replace connection details)
func connectDB() error {
	var err error
	db, err = sql.Open("postgres", "user=your_user dbname=your_db password=your_password sslmode=disable")
	if err != nil {
		return err
	}
	return db.Ping()
}

// bookSeat attempts to book a seat for a user in an event
func bookSeat(c echo.Context) error {
	var (
		err     error
		eventID int
		userID  int
	)

	if err := c.Bind(&eventID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid event ID")
	}
	if err := c.Bind(&userID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// Use transaction to ensure atomicity
	tx, err := db.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Check if user already has a booking for this event
	var hasBooking bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM Bookings WHERE event_id = $1 AND user_id = $2)", eventID, userID).Scan(&hasBooking)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking existing booking")
	}
	if hasBooking {
		return echo.NewHTTPError(http.StatusConflict, "User already has a booking for this event")
	}

	// Update available seats with a lock to prevent race conditions
	result, err := tx.Exec("UPDATE Events SET available_seats = available_seats - 1 WHERE event_id = $1 AND available_seats > 0", eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error updating available seats")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking rows affected")
	}
	if rowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "No seats available for this event")
	}

	// Insert booking record
	_, err = tx.Exec("INSERT INTO Bookings (event_id, user_id) VALUES ($1, $2)", eventID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating booking")
	}

	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}

func main() {
	err := connectDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	e := echo.New()

	// POST /bookseat endpoint to handle booking requests
	e.POST("/bookseat", bookSeat)

	e.Logger.Fatal(e.Start(":1323"))
}