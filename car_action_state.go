package main

type (
	/**
	This state is set if a car is in an action
	like waiting for the crossing to become free
	or it is on the crossing
	 */
	CarActionState struct {
		PosTileNo int
		CarNo int
		Lane int
		Speed int
	}
)