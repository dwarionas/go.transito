package booking

import (
	"errors"
	"testing"
	"time"
)

func TestBook(t *testing.T) {
	now := time.Date(2026, 12, 1, 10, 0, 0, 0, time.UTC)
	departure := now.Add(2 * time.Hour)

	newTrip := func(seatsAvailable int, status TripStatus) *Trip {
		return &Trip{
			ID:             "trip-1",
			From:           "Berlin",
			To:             "Lviv",
			DepartureTime:  departure,
			PriceCents:     2500,
			SeatsTotal:     24,
			SeatsAvailable: seatsAvailable,
			Status:         status,
		}
	}

	tests := []struct {
		name          string
		trip          *Trip
		seats         int
		now           time.Time
		wantErr       error
		wantSeatsLeft int
	}{
		// seccess
		{name: "one seat", trip: newTrip(20, TripOpen), seats: 1, now: now, wantSeatsLeft: 19},
		{name: "max 5 seats", trip: newTrip(20, TripOpen), seats: 5, now: now, wantSeatsLeft: 15},
		{name: "all remaining seats", trip: newTrip(3, TripOpen), seats: 3, now: now, wantSeatsLeft: 0},
		{name: "one minute before departure", trip: newTrip(20, TripOpen), seats: 1, now: departure.Add(-time.Minute), wantSeatsLeft: 19},

		// invalid amount of seats
		{name: "zero seats", trip: newTrip(20, TripOpen), seats: 0, now: now, wantErr: ErrInvalidSeats, wantSeatsLeft: 20},
		{name: "six seats", trip: newTrip(20, TripOpen), seats: 6, now: now, wantErr: ErrInvalidSeats, wantSeatsLeft: 20},
		{name: "negative seats", trip: newTrip(20, TripOpen), seats: -1, now: now, wantErr: ErrInvalidSeats, wantSeatsLeft: 20},

		// trip is not open
		{name: "trip closed", trip: newTrip(20, TripClosed), seats: 1, now: now, wantErr: ErrTripNotOpen, wantSeatsLeft: 20},
		{name: "trip cancelled", trip: newTrip(20, TripCancelled), seats: 1, now: now, wantErr: ErrTripNotOpen, wantSeatsLeft: 20},

		// trip is aready departed
		{name: "exactly at departure", trip: newTrip(20, TripOpen), seats: 1, now: departure, wantErr: ErrTripDeparted, wantSeatsLeft: 20},
		{name: "after departure", trip: newTrip(20, TripOpen), seats: 1, now: departure.Add(time.Hour), wantErr: ErrTripDeparted, wantSeatsLeft: 20},

		// not enough seats
		{name: "not enough seats", trip: newTrip(2, TripOpen), seats: 3, now: now, wantErr: ErrNotEnoughSeats, wantSeatsLeft: 2},
		{name: "no seats at all", trip: newTrip(0, TripOpen), seats: 1, now: now, wantErr: ErrNotEnoughSeats, wantSeatsLeft: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Book(tt.trip, "Ivan Davidov", "+380993131310", tt.seats, tt.now)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}

			if tt.trip.SeatsAvailable != tt.wantSeatsLeft {
				t.Errorf("seats left = %d, want %d", tt.trip.SeatsAvailable, tt.wantSeatsLeft)
			}

			if tt.wantErr != nil {
				if got != nil {
					t.Errorf("expected nil booking, got %+v", got)
				}
				return
			}

			if got.TripID != tt.trip.ID {
				t.Errorf("TripID = %q, want %q", got.TripID, tt.trip.ID)
			}
			if got.Seats != tt.seats {
				t.Errorf("Seats = %d, want %d", got.Seats, tt.seats)
			}
			if got.Status != BookingPending {
				t.Errorf("Status = %v, want %v", got.Status, BookingPending)
			}
			if !got.CreatedAt.Equal(tt.now) {
				t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, tt.now)
			}
		})
	}
}
