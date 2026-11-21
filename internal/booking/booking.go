package booking

import (
	"errors"
	"time"
)

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

var (
	ErrInvalidSeats   = errors.New("seats must be between 1 and 5")
	ErrTripNotOpen    = errors.New("trip is not open")
	ErrTripDeparted   = errors.New("trip is already departed")
	ErrNotEnoughSeats = errors.New("not enough seats")
)

func Book(trip *Trip, name, phone string, seats int, now time.Time) (*Booking, error) {
	if seats < 1 || seats > 5 {
		return nil, ErrInvalidSeats
	}

	if trip.Status != TripOpen {
		return nil, ErrTripNotOpen
	}

	if !now.Before(trip.DepartureTime) {
		return nil, ErrTripDeparted
	}

	if trip.SeatsAvailable < seats {
		return nil, ErrNotEnoughSeats
	}

	trip.SeatsAvailable -= seats

	return &Booking{
		TripID:         trip.ID,
		PassengerName:  name,
		PassengerPhone: phone,
		Seats:          seats,
		CreatedAt:      now,
		Status:         BookingPending,
	}, nil
}
