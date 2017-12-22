package main

import(
	queue "github.com/golang-collections/collections/queue"
	anki "github.com/okoeth/edge-anki-base"
)

type (
	// Command represents a valid command with can be sent to the Anki Overdrive controller
	Crossing struct {
		Tile1No int
		Tile2No int
		CrossingTileCarQueue queue.Queue
		CrossingWaitingCarQueue queue.Queue
	}
)

func canDriveCrossing(carNo int, track []anki.Status, crossing *Crossing) bool {
	var currentCarState= getStateForCarNo(carNo, track)
	var lastTileNo= (currentCarState.PosTileNo - 1) % currentCarState.MaxTileNo
	//var nextTileNo = (currentCarState.PosTileNo+1)%currentCarState.MaxTileNo

	//Check if last tile was crossing
	if lastTileNo == crossing.Tile1No || lastTileNo == crossing.Tile2No {
		crossing.CrossingTileCarQueue.Dequeue()
	}

	//Check if current tile is crossing
	if (currentCarState.PosTileNo == crossing.Tile1No || currentCarState.PosTileNo == crossing.Tile2No) &&
		crossing.CrossingTileCarQueue.Peek() != currentCarState.CarNo {
		if crossing.CrossingTileCarQueue.Len() > 0 {
			if crossing.CrossingWaitingCarQueue.Peek() != currentCarState.CarNo {
				crossing.CrossingWaitingCarQueue.Enqueue(currentCarState.CarNo)
			}
			mlog.Println("WARNING: Can not pass crossing")
			mlog.Printf("CrossingWaitingCarQueue: %+v\n", crossing.CrossingWaitingCarQueue)
			mlog.Printf("CrossingTileCarQueue: %+v\n", crossing.CrossingTileCarQueue)
			return false
		} else {
			crossing.CrossingTileCarQueue.Enqueue(currentCarState.CarNo)
		}
	}

	mlog.Println("DEBUG: Can pass crossing")
	return true
}