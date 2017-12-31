package main

import(
	"container/list"
	anki "github.com/okoeth/edge-anki-base"
	"time"
)

type (
	// Command represents a valid command with can be sent to the Anki Overdrive controller
	Crossing struct {
		Tile1No int
		Tile2No int
		CarsOnCrossing map[int]CarActionState
		CrossingWaitingCarQueue *list.List
	}
)

func NewCrossing(tile1No int, tile2No int) Crossing {
	return Crossing { Tile1No: tile1No,
	Tile2No: tile2No,
	CarsOnCrossing: make(map[int]CarActionState),
	CrossingWaitingCarQueue: list.New() }
}

func canDriveCrossing(carNo int, track []anki.Status, crossing *Crossing) bool {
	var currentCarState = getStateForCarNo(carNo, track)

	var lastTileNo = (currentCarState.PosTileNo - 1) % currentCarState.MaxTileNo
	if lastTileNo == -1 {
		lastTileNo = 7
	}

	var nextTileNo = (currentCarState.PosTileNo + 1) % currentCarState.MaxTileNo

	var currentDistanceTravelled = CalculateDistanceTravelled(float32(currentCarState.CarSpeed), getTimeDelta(currentCarState.TransitionTimestamp))
	var distanceToCrossing = getRemainingTileDistance(currentCarState.LaneLength, currentDistanceTravelled)
	var distanceInTimeStep = CalculateDistanceTravelled(float32(currentCarState.CarSpeed), 500)

	mlog.Printf("DEBUG: last tile %d", lastTileNo)
	mlog.Printf("DEBUG: next tile %d", nextTileNo)
	mlog.Printf("DEBUG: Distance to crossing %f", distanceToCrossing)

	//Check if last tile was crossing
	if lastTileNo == crossing.Tile1No || lastTileNo == crossing.Tile2No {
		delete(crossing.CarsOnCrossing, currentCarState.CarNo)
	}

	//Check outdated queue entries
	for key, value := range crossing.CarsOnCrossing {
		mlog.Printf("DEBUG: time since on crossing %f", time.Since(value.Timestamp).Seconds())
		if time.Since(value.Timestamp).Seconds() > 2 {
			delete(crossing.CarsOnCrossing, key)
		}
	}

	//Check if the car is currently on the crossing tile
	//And if it is already on the crossing (drive on)
	if (isCarCloseToCrossingTile(nextTileNo, crossing, distanceToCrossing, distanceInTimeStep) ||
		isCarOnCrossingTile(currentCarState.PosTileNo, crossing)) &&
		!isCarActiveOnCrossing(currentCarState.CarNo, crossing) {

		if len(crossing.CarsOnCrossing) > 0 &&
			!isCarGoingInSameDirectionAsActiveCar(currentCarState.PosTileNo, crossing)  &&
				!isCarGoingInSameDirectionAsActiveCar(nextTileNo, crossing){

			//Check if car is already waiting
			if _, inQueue := getCarFromWaitingQueue(carNo, crossing); !inQueue{

				crossing.CrossingWaitingCarQueue.PushBack(
					CarActionState{
						Timestamp: time.Now(),
						CarNo: carNo,
						Lane: currentCarState.LaneNo,
						Speed: currentCarState.CarSpeed })
			}

			mlog.Println("WARNING: Can not pass crossing")
			mlog.Printf("DEBUG: CrossingWaitingCarQueue: %+v\n", crossing.CrossingWaitingCarQueue)
			mlog.Printf("DEBUG: CrossingTileCarQueue: %+v\n", crossing.CarsOnCrossing)
			return false
		} else if isCarOnCrossingTile(currentCarState.PosTileNo, crossing) {
			crossing.CarsOnCrossing[carNo] =
				CarActionState{CarNo: carNo,
					Timestamp: time.Now(),
					Lane: currentCarState.LaneNo,
					Speed: currentCarState.CarSpeed,
					PosTileNo: currentCarState.PosTileNo}
		}
	}

	mlog.Println("DEBUG: Can pass crossing")
	return true
}

func isCarOnCrossingTile(posTileNo int, crossing *Crossing) bool {
	return posTileNo == crossing.Tile1No ||
		posTileNo== crossing.Tile2No
}

func isCarCloseToCrossingTile(nextTileNo int, crossing *Crossing,
	distanceToCrossing float64, distanceInTimeStep float64) bool {
	if nextTileNo == crossing.Tile1No ||
		nextTileNo == crossing.Tile2No {
			if distanceToCrossing < distanceInTimeStep {
				return true
			}
	}
	return false
}

func isCarActiveOnCrossing(carNo int, crossing *Crossing) bool {
	_, carActiveOnCrossing := crossing.CarsOnCrossing[carNo]
	return carActiveOnCrossing
}

func isCarGoingInSameDirectionAsActiveCar(posTileNo int, crossing *Crossing) bool {
	currentCrossingDirection := getTileOfFirstCarOnCrossing(crossing)
	return posTileNo == currentCrossingDirection
}

func getTileOfFirstCarOnCrossing(crossing *Crossing) int {
	for _, value := range crossing.CarsOnCrossing {
		return value.PosTileNo
	}
	return -1
}

func getCarFromWaitingQueue(carNo int, crossing *Crossing) (*list.Element, bool) {
	// Iterate through list and print its contents.
	for listElement := crossing.CrossingWaitingCarQueue.Front();
		listElement != nil; listElement = listElement.Next() {
		if listElement.Value.(CarActionState).CarNo == carNo {
			return listElement, true
		}
	}
	return nil, false
}

func tryRemoveCarFromQueue(carNo int, crossing *Crossing) (CarActionState, bool) {
	// here a reactivate has to happen if car is stopped
	if listElement, inQueue := getCarFromWaitingQueue(carNo, crossing); inQueue {
		crossing.CrossingWaitingCarQueue.Remove(listElement)
		stoppedCarState := listElement.Value.(CarActionState)
		return stoppedCarState, true
	}
	return CarActionState{}, false
}