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
		t.FailNow()
	}

	//Two cars, same lane, different tiles, enough distance
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
	PosLocation: 1, LaneNo: 2, MaxTileNo: 6, LaneLength: 800, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
	PosLocation: 1, LaneNo: 2, MaxTileNo: 6, LaneLength: 800, CarSpeed: 200}

	if !canDriveOn(0, track, nil) {
		t.FailNow()
	}

	//Two cars, same lane, same tiles, enough distance
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
	PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 200}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-1100 * time.Millisecond),
	PosTileNo: 2, PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 200}

	if !canDriveOn(0, track, nil) {
		t.FailNow()
	}

	//Two cars, same lane, same tiles, not enough distance
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-200 * time.Millisecond),
	PosTileNo: 2, PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 200}

	if canDriveOn(0, track, nil) {
		t.FailNow()
	}
}

func TestGetAvailableLane(t *testing.T) {

	//TODO: Change always two lanes?
	t.Log("TEST: Starting testGetAvailableLane")

	var track = getEmptyStatusArray()

	//Two cars, different lane, different tiles, should work
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, LaneLength: 800, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
		PosLocation: 1, LaneNo: 3, MaxTileNo: 6, LaneLength: 800, CarSpeed: 200}

	//Should change to lane 3
	if getAvailableLane(0, track, nil) != 3 {
		t.FailNow()
	}

	//Two cars, car 2 ahead 1 tile, car 2 left of car 1
	//Car would change to lane 1 as 3 is blocked
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-200 * time.Millisecond),
		PosTileNo: 2, PosLocation: 1, LaneNo: 3, MaxTileNo: 6, CarSpeed: 200}

	if getAvailableLane(0, track, nil) != 1 {
		t.FailNow()
	}

	//Two cars, same tile, car 2 left of car 1
	//Car would change to lane 1 as 3 is blocked
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-200 * time.Millisecond),
		PosTileNo: 2, PosLocation: 1, LaneNo: 3, MaxTileNo: 6, CarSpeed: 200}

	if getAvailableLane(0, track, nil) != 1 {
		t.FailNow()
	}

	//Two cars, car 0 is at most left lane
	//lane 3 is blocked
	//should give lane 2
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 4, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-200 * time.Millisecond),
		PosTileNo: 2, PosLocation: 1, LaneNo: 3, MaxTileNo: 6, CarSpeed: 200}

	if getAvailableLane(0, track, nil) != 2 {
		t.FailNow()
	}

	//Two cars, car 0 is at most right lane
	//lane 3 is blocked
	//should give lane 2
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 1, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-200 * time.Millisecond),
		PosTileNo: 2, PosLocation: 1, LaneNo: 3, MaxTileNo: 6, CarSpeed: 200}

	if getAvailableLane(0, track, nil) != 2 {
		t.FailNow()
	}

	//Block all lanes
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 4, MaxTileNo: 6, CarSpeed: 250}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 3, MaxTileNo: 6, CarSpeed: 250}
	track[2] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, CarSpeed: 250}
	track[3] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 1, MaxTileNo: 6, CarSpeed: 250}

	if getAvailableLane(0, track, nil) != -1 {
		t.FailNow()
	}
}

func TestAdjustSpeed(t *testing.T) {
	t.Log("TEST: Starting adjustSpeed")

	var track = getEmptyStatusArray()

	//Two cars, same lane, same tile, should reduce speed to 300 (slower car)
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, CarSpeed: 400, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, CarSpeed: 300, MaxTileNo: 6}

	if calculateSpeed(0, track) != 300 {
		t.FailNow()
	}

	//Two cars, different lane, same tile, should give speed 0
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 3, CarSpeed: 400, MaxTileNo: 6}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2, PosLocation: 1, LaneNo: 2, CarSpeed: 300, MaxTileNo: 6}

	if calculateSpeed(0, track) != 0 {
		t.FailNow()
	}
}

func TestCalculatePosition(t *testing.T) {
	t.Log("TEST: Calculate position")
	var distanceTravelled = CalculateDistanceTravelled(200, 200)

	if distanceTravelled - 37.4 > 0.0001 {
		t.FailNow()
	}

	var relativePosition = CalculateRelativePosition(1000, distanceTravelled)
	if relativePosition - 0.0374 > 0.0001 {
		t.FailNow()
	}
}

func TestObstacle(t *testing.T) {
	//If obstacle on road set timestamp, else delete timestamp
	t.Log("TEST: Starting testObstacle")

	//TODO: Calculation of position for obstacle needed x,y -> lane, tile is fixed

	var track = getEmptyStatusArray()

	//One car, obstacle on lane
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, LaneLength: 800, CarSpeed: 250}
	track[4] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), CarNo: -1, TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 6, LaneLength: 800, CarSpeed: 0}

	if canDriveOn(0, track, nil) {
		t.FailNow()
	}
}

func TestCrossing(t *testing.T) {
	t.Log("TEST: Starting testCrossing")
	crossing := NewCrossing(3, 6)

	//Crossing has tileNo 3
	var track = getEmptyStatusArray()

	//No car is on crossing, both could pass crossing
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 2,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 250, CarNo: 0}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 1,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 1}

	if !canDriveCrossing(0, track, &crossing) {
		t.FailNow()
	}

	if !canDriveCrossing(1, track, &crossing) {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 0 || len(crossing.CarsOnCrossing) != 0 {
		t.FailNow()
	}

	//Car one is now on crossing, car two is going in the same direction
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 250, CarNo: 0}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 1}

	if !canDriveCrossing(0, track, &crossing) {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 0 || len(crossing.CarsOnCrossing) != 1 {
		t.FailNow()
	}

	if !canDriveCrossing(1, track, &crossing) {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 0 || len(crossing.CarsOnCrossing) != 2 {
		t.FailNow()
	}

	//Car one is still on crossing, car two should also be on crossing
	//Car three wants to pass crossing from other direction -> wait
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 250, CarNo: 0}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 1}
	track[2] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 6,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 2}

	if !canDriveCrossing(0, track, &crossing) {
		t.FailNow()
	}

	if !canDriveCrossing(1, track, &crossing) {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 0 || len(crossing.CarsOnCrossing) != 2 {
		t.FailNow()
	}

	if canDriveCrossing(2, track, &crossing) {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 1 || len(crossing.CarsOnCrossing) != 2 {
		t.FailNow()
	}


	//Car one not anymore on crossing
	//Car two on crossing
	//Car three can still not pass
	//Car four wants to pass crossing from other direction -> wait
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 4,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 250, CarNo: 0}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 1}
	track[2] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 6,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 2}
	track[3] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 6,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 3}

	if !canDriveCrossing(0, track, &crossing) {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 1 || len(crossing.CarsOnCrossing) != 1 {
		t.FailNow()
	}

	if !canDriveCrossing(1, track, &crossing) {
		t.FailNow()
	}

	if canDriveCrossing(2, track, &crossing) {
		t.FailNow()
	}

	if canDriveCrossing(3, track, &crossing) {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 2 || len(crossing.CarsOnCrossing) != 1 {
		t.FailNow()
	}


	//Car one not anymore on crossing
	//Car two on crossing
	//Car three can still not pass
	//Car four can still not pass
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 4,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 250, CarNo: 0}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 1}
	track[2] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 6,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 2}
	track[3] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 6,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 3}

	if !canDriveCrossing(0, track, &crossing) {
		t.FailNow()
	}

	if !canDriveCrossing(1, track, &crossing) {
		t.FailNow()
	}

	if canDriveCrossing(2, track, &crossing) {
		t.FailNow()
	}

	if canDriveCrossing(3, track, &crossing) {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 2 || len(crossing.CarsOnCrossing) != 1 {
		t.FailNow()
	}

	//Car one not anymore on crossing
	//Car two not anymore on crossing
	//Car three can pass
	//Car four can pass
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 4,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 250, CarNo: 0}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 4,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 1}
	track[2] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 6,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 2}
	track[3] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 6,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 3}

	if !canDriveCrossing(0, track, &crossing) {
		t.FailNow()
	}

	if !canDriveCrossing(1, track, &crossing) {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 2 || len(crossing.CarsOnCrossing) != 0 {
		t.FailNow()
	}

	if !canDriveCrossing(2, track, &crossing) {
		t.FailNow()
	}

	if _, carInQueue := tryRemoveCarFromQueue(2, &crossing); !carInQueue {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 1 || len(crossing.CarsOnCrossing) != 1 {
		t.FailNow()
	}

	if !canDriveCrossing(3, track, &crossing) {
		t.FailNow()
	}

	if _, carInQueue := tryRemoveCarFromQueue(3, &crossing); !carInQueue {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 0 || len(crossing.CarsOnCrossing) != 2 {
		t.FailNow()
	}

	//Move car three, four from crossing, so that crossing is free
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 5,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 250, CarNo: 0}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 4,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 1}
	track[2] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 7,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 2}
	track[3] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 7,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 3}

	if !canDriveCrossing(0, track, &crossing) {
		t.FailNow()
	}

	if !canDriveCrossing(1, track, &crossing) {
		t.FailNow()
	}

	if !canDriveCrossing(2, track, &crossing) {
		t.FailNow()
	}

	if !canDriveCrossing(3, track, &crossing) {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 0 || len(crossing.CarsOnCrossing) != 0 {
		t.FailNow()
	}

	//Car 1 on crossing, car 2 wants to go on crossing from other side
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 250, CarNo: 0}
	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 6,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 0, CarNo: 1}

	if !canDriveCrossing(0, track, &crossing) {
		t.FailNow()
	}

	if canDriveCrossing(1, track, &crossing) {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 1 || len(crossing.CarsOnCrossing) != 1 {
		t.FailNow()
	}
}

func TestCloseToCrossing(t *testing.T) {
	t.Log("TEST: Starting testCrossing")
	crossing := NewCrossing(3, 6)

	//Crossing has tileNo 3
	var track= getEmptyStatusArray()

	//Car 1 on crossing,
	//Car 2 is one tile before crossing and wants to go on crossing from other side
	//Car 2 is far enough (can drive on)
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 250, CarNo: 0}

	if !canDriveCrossing(0, track, &crossing) {
		t.FailNow()
	}

	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-100 * time.Millisecond), PosTileNo: 5,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 800, CarNo: 1}

	if !canDriveCrossing(1, track, &crossing) {
		t.FailNow()
	}

	//Car 1 on crossing,
	//Car 2 is one tile before crossing and wants to go on crossing from other side
	//Car 2 is a bit away but very fast (has to break earlier)
	track[0] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now(), PosTileNo: 3,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 250, CarNo: 0}

	if !canDriveCrossing(0, track, &crossing) {
		t.FailNow()
	}

	track[1] = anki.Status{ MsgID: 70, MsgTimestamp: time.Now(), TransitionTimestamp: time.Now().Add(-200 * time.Millisecond), PosTileNo: 5,
		PosLocation: 1, LaneNo: 2, MaxTileNo: 8, LaneLength: 800, CarSpeed: 800, CarNo: 1}

	if canDriveCrossing(1, track, &crossing) {
		t.FailNow()
	}

	if crossing.CrossingWaitingCarQueue.Len() != 1 || len(crossing.CarsOnCrossing) != 1 {
		t.FailNow()
	}
}

func getEmptyStatusArray() []anki.Status {
	return []anki.Status { getEmptyStatus(), getEmptyStatus(), getEmptyStatus(), getEmptyStatus(), getEmptyStatus()}
}

func getEmptyStatus() anki.Status {
	return anki.Status{}
}
