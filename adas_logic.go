// Copyright 2018 NTT Group

// Permission is hereby granted, free of charge, to any person obtaining a copy of this
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify,
// merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following
// conditions:

// The above copyright notice and this permission notice shall be included in all copies
// or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
// FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package main

import (
	"time"
	"strconv"
	"math"
	anki "github.com/okoeth/edge-anki-base"
	"sync"
	"container/list"
)

var lock sync.Mutex
var crossing = NewCrossing(3, 7)
var driveAheadWaitingQueue *list.List

func driveCars(track []anki.Status, cmdCh chan anki.Command, statusCh chan anki.Status) {
	driveAheadWaitingQueue = list.New()

	//ticker := time.NewTicker(200 * 1e6) // 1e6 = ms, 1e9 = s
	ticker := time.NewTicker(200 * 1e6) // 1e6 = ms, 1e9 = s
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// TODO: Migrate carNo to int and solve 0 vs 1 issue

			logicBegin := time.Now()
			mlog.Println("INFO: ======== Logic begin =========")

			// Update tile ids for obstacle
			for i, car := range track {
				//-1 and -2 are obstacle
				if car.CarNo == -1 {
					track[i].PosTileNo = crossing.Tile1No
					track[i].TransitionTimestamp = car.MsgTimestamp
				} else if car.CarNo == -2 {
					track[i].PosTileNo = crossing.Tile2No
					track[i].TransitionTimestamp = car.MsgTimestamp
				}
			}

			// Drive cars
			for i, car := range track {
				//-1 and -2 are obstacle
				if car.CarNo != -1 && car.CarNo != -2 {
					driveCar(i, track, cmdCh)
				}
			}

			mlog.Printf("INFO: ======= Logic ended after %f ms ==========", time.Since(logicBegin).Seconds()*1000)
		}
	}
}

func driveCar(carNo int, track []anki.Status, cmdCh chan anki.Command) {
	defer anki.Track_execution_time(anki.Start_execution_time("driveCar"))

	var currentCarState = getStateForCarNo(carNo, track)
	if messageWithIntegrity(currentCarState) {
		mlog.Printf("CarNo %d, %+v", carNo, currentCarState)
		if canDriveCrossing(carNo, track, &crossing, cmdCh) {
			if canDriveOn(carNo, track, cmdCh) {
				if currentCarState.LightsOn {
					//Reset light
					track[carNo].LightsOn = false
					lightsOff(carNo, cmdCh)
				}

				driveAhead(carNo, track, cmdCh)
			} else {
				if !currentCarState.LightsOn {
					track[carNo].LightsOn = true
					lightsOn(carNo, cmdCh)
				}

				if availableLane := getAvailableLane(carNo, track, cmdCh); availableLane != -1 {
					driveAhead(carNo, track, cmdCh)
					changeLane(carNo, track, cmdCh, availableLane)
				} else {
					adjustSpeed(carNo, track, cmdCh)
				}
			}
		} else {
			stopCar(carNo, cmdCh)
		}
	}
}

func messageWithIntegrity(currentCarState anki.Status) bool {
	if !currentCarState.MsgTimestamp.IsZero() &&
		!currentCarState.TransitionTimestamp.IsZero() &&
			currentCarState.MaxTileNo != 0 &&
				currentCarState.LaneLength != 0{
			return true
	}
	return false
}

func stopCar(carNo int, cmdCh chan anki.Command) {
	mlog.Printf("DEBUG: Stopping car")
	cmd := anki.Command{ CarNo: carNo, Command: "s", Param1: strconv.Itoa(0)}
	cmdCh <- cmd
}

/**
CRITERIA:
Check the lane of the current car
If there is any car on the current lane and the same tile
	return false
Else
	return true
 */
func canDriveOn(carNo int, track []anki.Status, cmdCh chan anki.Command) bool {
	defer anki.Track_execution_time(anki.Start_execution_time("canDriveOn"))

	var currentCarState = getStateForCarNo(carNo, track)

	//Check all other car states
	for index, otherCarState := range track {
		if timeStampsValid(otherCarState) && index != carNo &&
			otherCarState.CarSpeed < currentCarState.CarSpeed &&
			hasCarInFront(otherCarState, currentCarState, currentCarState.LaneNo) {
			mlog.Printf("WARNING: Other car on same lane")
			return false
		}
	}
	mlog.Printf("DEBUG: Can drive on")
	return true
}


func getAvailableLane(carNo int, track []anki.Status, cmdCh chan anki.Command) int {
	defer anki.Track_execution_time(anki.Start_execution_time("getAvailableLane"))

	var currentCarState= getStateForCarNo(carNo, track)

	//Check for all possible lanes
	var suggestedLaneIndex= -1
	var laneAvailable= false

	for _, laneOffset := range []int{1, -1, 2, -2, 3, -3} {
		laneAvailable = true
		suggestedLaneIndex = currentCarState.LaneNo + laneOffset

		if suggestedLaneIndex < 1 || suggestedLaneIndex > 4 {
			continue
		}

		//Check all other car states of same and next tile
		for index, otherCarState := range track {

			if timeStampsValid(otherCarState) && index != carNo &&
				hasCarInFront(otherCarState, currentCarState, suggestedLaneIndex) {
				mlog.Printf("WARNING: Other car on lane %d, no change possible", suggestedLaneIndex)
				laneAvailable = false
				break
			}
		}

		if laneAvailable {
			mlog.Printf("DEBUG: Can change to lane %d", suggestedLaneIndex)
			return suggestedLaneIndex
		}
	}

	return -1
}

func timeStampsValid(carState anki.Status) bool {
	msgTimestampZero := carState.MsgTimestamp.IsZero()
	transitionTimestampZero := carState.TransitionTimestamp.IsZero()
	//mlog.Printf("Timestamps for car %d, msgZero: %b, transitionZero %b", carState.CarNo, msgTimestampZero, transitionTimestampZero)
	return !msgTimestampZero && !transitionTimestampZero
}

func getTimeDelta(transitionTimestamp time.Time) float64 {
	return time.Since(transitionTimestamp).Seconds() * 1000
}

func getRemainingTileDistance(laneLength int, currentDistanceTravelled float64) float64 {
	return float64(laneLength) - currentDistanceTravelled
}

func hasCarInFront(otherCarState anki.Status, currentCarState anki.Status, laneNo int) bool {
	var currentTimeDelta = getTimeDelta(currentCarState.TransitionTimestamp)
	var currentDistanceTravelled = CalculateDistanceTravelled(float32(currentCarState.CarSpeed), currentTimeDelta)
	var otherTimeDelta = getTimeDelta(otherCarState.TransitionTimestamp)
	var otherDistanceTravelled = CalculateDistanceTravelled(float32(otherCarState.CarSpeed), otherTimeDelta)

	//Lane change takes about 900ms, so we have to change far before
	var distanceInTimeStep = math.Max(200.0, CalculateDistanceTravelled(float32(currentCarState.CarSpeed), 900))

	var nextTileNo = (currentCarState.PosTileNo+1) % currentCarState.MaxTileNo

	// Other car must be on same lane or neighbor lane and on next or current tile
	if (otherCarState.LaneNo == laneNo ||
		otherCarState.LaneNo == laneNo-1 ||
		otherCarState.LaneNo == laneNo+1) && (
		otherCarState.PosTileNo == currentCarState.PosTileNo ||
			otherCarState.PosTileNo == nextTileNo) {

		var distanceDelta float64 = -1

		//1. Check if cars are on same tiles
		if otherCarState.PosTileNo == currentCarState.PosTileNo &&
			(math.Floor(currentDistanceTravelled) <= math.Floor(otherDistanceTravelled) ||
				currentTimeDelta < otherTimeDelta) {
			distanceDelta = otherDistanceTravelled - currentDistanceTravelled
		} else if otherCarState.PosTileNo == nextTileNo {
			//2. Check if other car is on next tile
			distanceDelta = otherDistanceTravelled +
				getRemainingTileDistance(currentCarState.LaneLength, currentDistanceTravelled)
		}

		mlog.Printf("DEBUG: Next probable position %f", distanceInTimeStep)
		mlog.Printf("DEBUG: Car pos: tile %d, pos: %f", currentCarState.PosTileNo, currentDistanceTravelled)
		mlog.Printf("DEBUG: Car pos other: tile %d, pos: %f", otherCarState.PosTileNo, otherDistanceTravelled)
		mlog.Printf("DEBUG: Car speed: %d", currentCarState.CarSpeed)
		mlog.Printf("DEBUG: Car speed other: %d", otherCarState.LaneNo)
		mlog.Printf("DEBUG: Car lane: %d", currentCarState.LaneNo)
		mlog.Printf("DEBUG: Car lane other: %d", otherCarState.CarSpeed)
		mlog.Printf("DEBUG: Car time delta: %f", currentTimeDelta)
		mlog.Printf("DEBUG: Car time delta other: %f", otherTimeDelta)
		mlog.Printf("DEBUG: Lane length %d", currentCarState.LaneLength)
		mlog.Printf("DEBUG: lane distance delta %f", getRemainingTileDistance(currentCarState.LaneLength, currentDistanceTravelled))
		mlog.Println("DEBUG: Distance is ", distanceDelta)

		// Check if distance is enough in respect to speed
		// It is not enough if the distanceDelta is lower than
		// What the car would travel in 1.5 intervals
		if distanceDelta > -1 && distanceDelta < distanceInTimeStep {
			return true
		}
	}

	return false
}

/**
CRITERIA:
Find car that is blocking us and adjust speed to the speed of the blocking car
 */
func adjustSpeed(carNo int, track []anki.Status, cmdCh chan anki.Command) bool {
	var currentCarState = getStateForCarNo(carNo, track)
	speed := calculateSpeed(carNo, track)

	if speed < currentCarState.CarSpeed {
		driveAheadWaitingQueue.PushBack(
			CarActionState{
				Timestamp: time.Now(),
				CarNo: carNo,
				Lane: currentCarState.LaneNo,
				Speed: currentCarState.CarSpeed })
	}

	//Change speed according to car before
	cmd := anki.Command{ CarNo: carNo, Command: "s", Param1: string(speed)}
	cmdCh <- cmd
	return true
}

func lightsOn(carNo int, cmdCh chan anki.Command) {
	cmd := anki.Command{ CarNo: carNo, Command: "l", }
	cmdCh <- cmd
}

func lightsOff(carNo int, cmdCh chan anki.Command) {
	cmd := anki.Command{ CarNo: carNo, Command: "lp", }
	cmdCh <- cmd
}

func calculateSpeed(carNo int, track []anki.Status) int {
	var currentCarState = getStateForCarNo(carNo, track)
	var blockingCarState anki.Status

	//Check all other car states
	for index, otherCarState := range track {
		if !otherCarState.MsgTimestamp.IsZero() && index != carNo && otherCarState.LaneNo == currentCarState.LaneNo &&
			otherCarState.PosTileNo == currentCarState.PosTileNo {
			mlog.Printf("WARNING: Other car in front")
			blockingCarState = otherCarState
			return blockingCarState.CarSpeed
		}
	}
	return 0
}

/**
Simply drive on
 */
func driveAhead(carNo int, track []anki.Status, cmdCh chan anki.Command) {
	defer anki.Track_execution_time(anki.Start_execution_time("driveAhead"))

	mlog.Printf("INFO: Drive ahead")

	mlog.Printf("CrossingWaitingCarQueue: %+v\n", crossing.CrossingWaitingCarQueue)
	mlog.Printf("CrossingTileCarQueue: %+v\n", crossing.CarsOnCrossing)

	// here a reactivate has to happen if car is stopped, and waiting in queue
	if carActionState, inQueue := tryRemoveCarFromQueue(carNo, crossing.CrossingWaitingCarQueue); inQueue {
		mlog.Printf("DEBUG: Reactivating car from crossing waiting with speed %d\n", carActionState.Speed)
		cmd := anki.Command{CarNo: carNo, Command: "s", Param1: strconv.Itoa(carActionState.Speed)}
		cmdCh <- cmd
	}

	// if drive ahead stopped a car entirely, assemble old speed here
	if carActionState, inQueue := tryRemoveCarFromQueue(carNo, driveAheadWaitingQueue); inQueue {
		mlog.Printf("DEBUG: Reactivating car from drive ahead queue with speed %d\n", carActionState.Speed)
		cmd := anki.Command{CarNo: carNo, Command: "s", Param1: strconv.Itoa(carActionState.Speed)}
		cmdCh <- cmd
	}

	// if car is faster than 700, limit to 700
	carState := getStateForCarNo(carNo, track)
	if carState.CarSpeed > 650 {
		cmd := anki.Command{CarNo: carNo, Command: "s", Param1: strconv.Itoa(600)}
		cmdCh <- cmd
	}
}

/**
Initiate left change
 */
func changeLane(carNo int, track []anki.Status, cmdCh chan anki.Command, laneNo int) {
	mlog.Printf("INFO: Changing to lane %d", laneNo)
	cmd := anki.Command{ CarNo: carNo, Command: "c", Param1: strconv.Itoa(laneNo)}
	cmdCh <- cmd

	mlog.Printf("Command sent %+v\n", cmd)
}

func getStateForCarNo(carNo int, track []anki.Status) anki.Status {
	return track[carNo]
}