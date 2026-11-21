package booking

import "time"

type Trip struct {
	ID             string
	From           string
	To             string
	DepartureTime  time.Time
	PriceCents     int
	SeatsTotal     int
	SeatsAvailable int
	Status         TripStatus
}

type TripStatus string

const (
	TripOpen      TripStatus = "open"
	TripDeparted  TripStatus = "departed"
	TripClosed    TripStatus = "closed"
	TripCancelled TripStatus = "cancelled"
)
