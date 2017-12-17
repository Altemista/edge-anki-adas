package main

import (
	"testing"
	 anki "github.com/okoeth/edge-anki-base"
	"time"
)

func TestCanDriveOn(t *testing.T) {
	t.Log("TEST: Starting testCanDriveOn")

	var track = getEmptyStatusArray()

	//Two cars, same lane, different tiles, far away
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, LaneLength: 800, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 5,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, LaneLength: 800, CarSpeed: 200}

	if !canDriveOn(0, track, nil) {
		t.Failed()
	}

	//Two cars, same lane, different tiles, enough distance
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
	PosLocation: 1, LaneNo: 2, MaxTileNo: 6, LaneLength: 800, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
	PosLocation: 1, LaneNo: 2, MaxTileNo: 6, LaneLength: 800, CarSpeed: 200}

	if !canDriveOn(0, track, nil) {
		t.Failed()
	}

	//Two cars, same lane, same tiles, enough distance
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
	PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-400 * time.Millisecond),
	PosTileNo: 2, PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 200}

	if !canDriveOn(0, track, nil) {
		t.Failed()
	}

	//Two cars, same lane, same tiles, not enough distance
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-200 * time.Millisecond),
	PosTileNo: 2, PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 200}

	if canDriveOn(0, track, nil) {
		t.Failed()
	}
}

func TestCanDriveLeft(t *testing.T) {
	t.Log("TEST: Starting testCanDriveLeft")

	var track = getEmptyStatusArray()

	//Two cars, different lane, different tiles, should work
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, LaneLength: 800, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
		PosLocation: 1, LaneNo: 3, MaxTileNo: 6, LaneLength: 800, CarSpeed: 200}

	if !canChangeLeft(0, track, nil) {
		t.Failed()
	}

	//Two cars, different lane, car 2 on same tile but far enough, car 2 left of car 1
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-400 * time.Millisecond),
		PosTileNo: 2, PosLocation: 1, LaneNo: 3, MaxTileNo: 6, CarSpeed: 200}

	if canChangeLeft(0, track, nil) {
		t.Failed()
	}

	//Two cars, different lane, same tile, not enough distance
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-200 * time.Millisecond),
		PosTileNo: 2, PosLocation: 1, LaneNo: 3, MaxTileNo: 6, CarSpeed: 200}

	if canChangeLeft(0, track, nil) {
		t.Failed()
	}

	//Two cars, different lane, car 2 on same tile but far enough, car 2 left of car 1, car 0 at most left lane
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 4, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-400 * time.Millisecond),
		PosTileNo: 2, PosLocation: 1, LaneNo: 3, MaxTileNo: 6, CarSpeed: 200}

	if canChangeLeft(0, track, nil) {
		t.Failed()
	}
}

func TestCanDriveRight(t *testing.T) {
	t.Log("TEST: Starting testCanDriveRight")

	var track = getEmptyStatusArray()

	//Two cars, different lane, different tiles, should work
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 3, MaxTileNo: 6, LaneLength: 800, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, LaneLength: 800, CarSpeed: 200}

	if !canChangeRight(0, track, nil) {
		t.Failed()
	}

	//Two cars, different lane, car 2 on same tile but far enough, car 2 left of car 1
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 3, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-400 * time.Millisecond),
		PosTileNo: 2, PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 200}

	if canChangeRight(0, track, nil) {
		t.Failed()
	}

	//Two cars, different lane, same tile, not enough distance
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 3, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-200 * time.Millisecond),
		PosTileNo: 2, PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 200}

	if canChangeRight(0, track, nil) {
		t.Failed()
	}

	//Two cars, different lane, car 2 on same tile but far enough, car 2 left of car 1, car 0 at most left lane
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 1, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-400 * time.Millisecond),
		PosTileNo: 2, PosLocation: 1, LaneNo: 3, MaxTileNo: 6, CarSpeed: 200}
}

func TestAdjustSpeed(t *testing.T) {
	t.Log("TEST: Starting adjustSpeed")

	var track = getEmptyStatusArray()

	//Two cars, same lane, same tile, should reduce speed to 300 (slower car)
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, CarSpeed: 400, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, CarSpeed: 300, MaxTileNo: 6}

	if calculateSpeed(0, track) != 300 {
		t.Failed()
	}

	//Two cars, different lane, same tile, should give speed 0
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 3, CarSpeed: 400, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, CarSpeed: 300, MaxTileNo: 6}

	if calculateSpeed(0, track) != 0 {
		t.Failed()
	}
}

func TestCalculatePosition(t *testing.T) {
	t.Log("TEST: Calculate position")
	var distanceTravelled = CalculateDistanceTravelled(200, 200)

	if distanceTravelled - 37.4 > 0.0001 {
		t.Failed()
	}

	var relativePosition = CalculateRelativePosition(1000, distanceTravelled)
	if relativePosition - 0.0374 > 0.0001 {
		t.Failed()
	}
}

func getEmptyStatusArray() []anki.Status {
	return []anki.Status { getEmptyStatus(), getEmptyStatus(), getEmptyStatus(), getEmptyStatus()};
}

func getEmptyStatus() anki.Status {
	return anki.Status{}
}
