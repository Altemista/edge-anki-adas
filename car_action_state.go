package main

import "time"

type (
	/**
	This state is set if a car is in an action
	like waiting for the crossing to become free
	or it is on the crossing
	 */
	CarActionState struct {
		Timestamp time.Time
		PosTileNo int
		CarNo int
		Lane int
		Speed int
	}
)