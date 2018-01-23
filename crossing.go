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
	defer anki.Track_execution_time(anki.Start_execution_time("canDriveCrossing"))

	lockTime := time.Now()
	lock.Lock()
	defer lock.Unlock()
	mlog.Printf("INFO: ======= Waited at canDriveCrossing lock for %f ms =======", time.Since(lockTime).Seconds()*1000)

	var currentCarState = getStateForCarNo(carNo, track)

	var lastTileNo = (currentCarState.PosTileNo - 1) % currentCarState.MaxTileNo
	if lastTileNo == -1 {
		lastTileNo = 7
	}
	var nextTileNo = (currentCarState.PosTileNo + 1) % currentCarState.MaxTileNo

	var currentDistanceTravelled = CalculateDistanceTravelled(float32(currentCarState.CarSpeed), getTimeDelta(currentCarState.TransitionTimestamp))
	var distanceToCrossing = getRemainingTileDistance(currentCarState.LaneLength, currentDistanceTravelled)
	var distanceInTimeStep = CalculateDistanceTravelled(float32(currentCarState.CarSpeed), 600)

	mlog.Printf("DEBUG: current tile %d", currentCarState.PosTileNo)
	mlog.Printf("DEBUG: last tile %d", lastTileNo)
	mlog.Printf("DEBUG: next tile %d", nextTileNo)
	mlog.Printf("DEBUG: Distance to crossing %f", distanceToCrossing)
	mlog.Printf("DEBUG: Speed %d", currentCarState.CarSpeed)
	mlog.Printf("DEBUG: Distance in timestep %f", distanceInTimeStep)

	//Check if last tile was crossing
	if lastTileNo == crossing.Tile1No || lastTileNo == crossing.Tile2No {
		delete(crossing.CarsOnCrossing, currentCarState.CarNo)
	}

	//Search backwards for the last crossing tile id and clean-up
	for i := 1; i < 4; i++ {
		nextSearchTileNo := (currentCarState.PosTileNo - i) % currentCarState.MaxTileNo
		if nextSearchTileNo == -1 {
			nextSearchTileNo = 7
		}

		if nextSearchTileNo == crossing.Tile1No || nextSearchTileNo == crossing.Tile2No {
			delete(crossing.CarsOnCrossing, currentCarState.CarNo)
			break
		} 
	}

	//Check outdated queue entries
	for key, value := range crossing.CarsOnCrossing {
		mlog.Printf("DEBUG: time since on crossing %f", time.Since(value.Timestamp).Seconds())
		if time.Since(value.Timestamp).Seconds() > 2 {
			delete(crossing.CarsOnCrossing, key)
		}
	}

	//Check if the car is currently on the crossing tile or close or already waiting (braking)
	if (isCarCloseToCrossingTile(nextTileNo, crossing, distanceToCrossing, distanceInTimeStep) ||
		isCarOnCrossingTile(currentCarState.PosTileNo, crossing) ||
		isCarInWaitingQueue(currentCarState.CarNo, crossing)) &&
		!isCarActiveOnCrossing(currentCarState.CarNo, crossing) {


			//Check if other cars are on the crossing
			mlog.Printf("DEBUG: length of cars on crossing %d", len(crossing.CarsOnCrossing))
		if len(crossing.CarsOnCrossing) > 0 &&
			!isCarGoingInSameDirectionAsActiveCar(currentCarState.PosTileNo, crossing)  &&
				!isCarGoingInSameDirectionAsActiveCar(nextTileNo, crossing){

			//Check if car is already waiting
			if _, inQueue := getCarFromWaitingQueue(carNo, crossing.CrossingWaitingCarQueue); !inQueue{
				crossing.CrossingWaitingCarQueue.PushBack(
					CarActionState{
						Timestamp: time.Now(),
						CarNo: carNo,
						Lane: currentCarState.LaneNo,
						Speed: currentCarState.CarSpeed })
			}

			mlog.Println("WARNING: Can not pass crossing")
			mlog.Printf("DEBUG: Speed: %d\n", currentCarState.CarSpeed)
			mlog.Printf("DEBUG: CrossingWaitingCarQueue: %+v\n", crossing.CrossingWaitingCarQueue)
			mlog.Printf("DEBUG: CrossingTileCarQueue: %+v\n", crossing.CarsOnCrossing)
			return false
		} else {
			var crossingNo = -1

			if isCarOnCrossingTile(currentCarState.PosTileNo, crossing) {
				crossingNo = currentCarState.PosTileNo
			} else {
				//Search forward for the next crossing tile id
				for i := 1; i < 4; i++ {
					nextSearchTileNo := (currentCarState.PosTileNo + i) % currentCarState.MaxTileNo
					if nextSearchTileNo == crossing.Tile1No {
						crossingNo = crossing.Tile1No
						break
					} else if nextSearchTileNo == crossing.Tile2No {
						crossingNo = crossing.Tile2No
						break
					}
				}
			}

			if crossingNo != -1 {
				mlog.Println("DEBUG: Adding car to cars on crossing")
				crossing.CarsOnCrossing[carNo] =
					CarActionState{CarNo: carNo,
						Timestamp: time.Now(),
						Lane: currentCarState.LaneNo,
						Speed: currentCarState.CarSpeed,
						PosTileNo: crossingNo}
			}
		}
	}

	mlog.Println("DEBUG: Can pass crossing")
	return true
}

func isCarOnCrossingTile(posTileNo int, crossing *Crossing) bool {
	if posTileNo == crossing.Tile1No ||
		posTileNo == crossing.Tile2No {
			mlog.Println("DEBUG: Car is on crossing tile")
			return true
	}
	return false
}

func isCarInWaitingQueue(carNo int, crossing *Crossing) bool {
	if _, inQueue := getCarFromWaitingQueue(carNo, crossing.CrossingWaitingCarQueue); inQueue {
		mlog.Println("DEBUG: Car is already in waiting queue")
		return true
	}
	return false
}

func isCarCloseToCrossingTile(nextTileNo int, crossing *Crossing,
	distanceToCrossing float64, distanceInTimeStep float64) bool {
	if nextTileNo == crossing.Tile1No ||
		nextTileNo == crossing.Tile2No {
			if distanceToCrossing < distanceInTimeStep {
				mlog.Println("DEBUG: Car is close to crossing tile")
				return true
			}
	}
	mlog.Println("DEBUG: Car is not close to crossing tile")
	return false
}

func isCarActiveOnCrossing(carNo int, crossing *Crossing) bool {
	_, carActiveOnCrossing := crossing.CarsOnCrossing[carNo]
	if carActiveOnCrossing {
		mlog.Println("DEBUG: Car is already active on crossing")
		return true
	}
	mlog.Println("DEBUG: Car is not active on crossing")
	return false
}

func isCarGoingInSameDirectionAsActiveCar(posTileNo int, crossing *Crossing) bool {
	currentCrossingDirection := getTileOfFirstCarOnCrossing(crossing)
	if posTileNo == currentCrossingDirection {
		mlog.Println("DEBUG: Car is going in same direction as active car")
		return true
	}
	return false
}

func getTileOfFirstCarOnCrossing(crossing *Crossing) int {
	for _, value := range crossing.CarsOnCrossing {
		return value.PosTileNo
	}
	return -1
}

func getCarFromWaitingQueue(carNo int, list *list.List) (*list.Element, bool) {
	// Iterate through list and print its contents.
	for listElement := list.Front();
		listElement != nil; listElement = listElement.Next() {
		if listElement.Value.(CarActionState).CarNo == carNo {
			return listElement, true
		}
	}
	return nil, false
}

func tryRemoveCarFromQueue(carNo int, list *list.List) (CarActionState, bool) {
	// here a reactivate has to happen if car is stopped
	if listElement, inQueue := getCarFromWaitingQueue(carNo, list); inQueue {
		crossing.CrossingWaitingCarQueue.Remove(listElement)
		stoppedCarState := listElement.Value.(CarActionState)
		return stoppedCarState, true
	}
	return CarActionState{}, false
}