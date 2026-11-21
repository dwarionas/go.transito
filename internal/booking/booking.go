package booking

import "time"

type Booking struct {
	ID             string
	TripID         string
	PassengerName  string
	PassengerPhone string
	Seats          int
	Token          string
	CreatedAt      time.Time
	Status         BookingStatus
}

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingRejected  BookingStatus = "rejected"
	BookingCancelled BookingStatus = "cancelled"
)
