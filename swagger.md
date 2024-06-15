swagger: "2.0"
info:
  title: Event Booking System API
  description: API for booking seats in events
  version: "1.0"

# Paths for your API endpoints
paths:
  /bookseat:
    post:
      summary: Book a seat for an event
      description: Attempts to book a seat for a user in an event.
      tags: [events]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                eventID:
                  type: integer
                  description: The ID of the event.
                userID:
                  type: integer
                  description: The ID of the user booking the seat.
      responses:
        "200":
          description: Booking successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    description: Whether the booking was successful.
        "400":
          description: Bad request (invalid event ID or user ID)
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: string
                    description: Error message.
        "409":
          description: Conflict (user already has a booking for this event)
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: string
                    description: Error message.
        "500":
          description: Internal Server Error
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: string
                    description: Error message.
