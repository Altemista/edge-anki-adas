package main

import (
	"testing"
	 anki "github.com/okoeth/edge-anki-base"
	"time"
)

func TestCanDriveOn(t *testing.T) {
	t.Log("TEST: Starting testCanDriveOn")

	var track = getEmptyStatusArray()

	//Two cars, same lane, different tiles
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 4, PosLocation: 1, LaneNo: 2, MaxTileNo: 6}

	if !canDriveOn(0, track, nil) {
		t.Failed()
	}

	//Two cars, same lane, same tiles
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 3, PosLocation: 1, LaneNo: 2, MaxTileNo: 6}

	if canDriveOn(0, track, nil) {
		t.Failed()
	}
}

func TestGetAvailableLane(t *testing.T) {
	t.Log("TEST: Starting testGetAvailableLane")

	var track = getEmptyStatusArray()

	//Two cars, different lane, different tiles, should work
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 4, PosLocation: 1, LaneNo: 3, MaxTileNo: 6}

	//Should change to lane 3
	if getAvailableLane(0, track, nil) != 3 {
		t.Failed()
	}

	//Two cars, car 2 ahead 1 tile, car 2 left of car 1
	//Car would change to lane 1 as 3 is blocked
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 3, PosLocation: 1, LaneNo: 3, MaxTileNo: 6}

	if getAvailableLane(0, track, nil) != 1 {
		t.Failed()
	}

	//Two cars, same tile, car 2 left of car 1
	//Car would change to lane 1 as 3 is blocked
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 3, MaxTileNo: 6}

	if getAvailableLane(0, track, nil) != 1 {
		t.Failed()
	}

	//Two cars, car 0 is at most left lane
	//lane 3 is blocked
	//should give lane 2
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 4, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 3, MaxTileNo: 6}

	if getAvailableLane(0, track, nil) != 2 {
		t.Failed()
	}

	//Two cars, car 0 is at most right lane
	//lane 3 is blocked
	//should give lane 2
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 0, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 3, MaxTileNo: 6}

	if getAvailableLane(0, track, nil) != 2 {
		t.Failed()
	}

	//Block all lanes
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 4, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 3, MaxTileNo: 6}
	track[2] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, MaxTileNo: 6}
	track[3] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 1, MaxTileNo: 6}

	if getAvailableLane(0, track, nil) != -1 {
		t.Failed()
	}
}

func TestAdjustSpeed(t *testing.T) {
	t.Log("TEST: Starting adjustSpeed")

	var track = getEmptyStatusArray()

	//Two cars, same lane, same tile, should reduce speed to 300 (slower car)
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, CarSpeed: 400, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, CarSpeed: 300, MaxTileNo: 6}

	if calculateSpeed(0, track) != 300 {
		t.Failed()
	}

	//Two cars, different lane, same tile, should give speed 0
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 3, CarSpeed: 400, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, CarSpeed: 300, MaxTileNo: 6}

	if calculateSpeed(0, track) != 0 {
		t.Failed()
	}
}

func getEmptyStatusArray() []anki.Status {
	return []anki.Status { getEmptyStatus(), getEmptyStatus(), getEmptyStatus(), getEmptyStatus()};
}

func getEmptyStatus() anki.Status {
	return anki.Status{}
}
