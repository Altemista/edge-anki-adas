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

	anki "github.com/okoeth/edge-anki-base"
	"strconv"
)

func driveCars(track []anki.Status, cmdCh chan anki.Command, statusCh chan anki.Status) {
	//ticker := time.NewTicker(200 * 1e6) // 1e6 = ms, 1e9 = s
	ticker := time.NewTicker(200 * 1e6) // 1e6 = ms, 1e9 = s
	defer ticker.Stop()
	for {
		select {
		case s := <-statusCh:
			mlog.Printf("INFO: Received status update from channel")
			anki.UpdateTrack(track, s)
		case <-ticker.C:
			// TODO: Migrate carNo to int and solve 0 vs 1 issue
			for i, object := range track {
				if object.CarID != "-1" {
					driveCar(i, track, cmdCh)
				}
			}
		}
	}
}

func driveCar(carNo int, track []anki.Status, cmdCh chan anki.Command) {
	//mlog.Printf("Status of carNo %d, %+v\n", carNo, track)

	//Get position info for car
	//for i := range track {
	//	if track[i].CarNo == carNo {
	//		mlog.Printf("Track of carNo %d, %+v\n", carNo, track[i])
	//		break
	//	}
	//}

	//TODO:
	//Can't drive on
	//While lane > 0 and lane <= 4
	//	Try left (lane)
	//	Try right (lane)
	//If no change possible, adjust speed

	//TODO: Iterative (recursive?) approach to find most left lane etc.
	var currentCarState = getStateForCarNo(carNo, track)
	if !currentCarState.MsgTimestamp.IsZero() &&
		!currentCarState.TransitionTimestamp.IsZero(){
		mlog.Printf("CarNo %d", carNo)
		if canDriveOn(carNo, track, cmdCh) {
			driveAhead(carNo, track, cmdCh)
		} else {
			if availableLane := getAvailableLane(carNo, track, cmdCh); availableLane != -1 {
				driveAhead(carNo, track, cmdCh)
				changeLane(carNo, track, cmdCh, availableLane)
			} else {
				adjustSpeed(carNo, track, cmdCh)
			}
		}
	}
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
	var currentCarState = getStateForCarNo(carNo, track)

	//Check all other car states
	for index, otherCarState := range track {
		if timeStampsValid(otherCarState) && index != carNo &&
			hasCarInFront(otherCarState, currentCarState, currentCarState.LaneNo) &&
				otherCarState.CarSpeed < currentCarState.CarSpeed {
			mlog.Printf("WARNING: Other car on same lane")
			return false
		}
	}
	mlog.Printf("DEBUG: Can drive on")
	return true
}


func getAvailableLane(carNo int, track []anki.Status, cmdCh chan anki.Command) int {
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
	return !carState.MsgTimestamp.IsZero() && !carState.TransitionTimestamp.IsZero()
}

func hasCarInFront(otherCarState anki.Status, currentCarState anki.Status, laneNo int) bool {
	var currentTimeDelta = time.Since(currentCarState.TransitionTimestamp).Seconds() * 1000
	var currentDistanceTravelled = CalculateDistanceTravelled(currentCarState.CarSpeed, currentTimeDelta)
	var otherTimeDelta = time.Since(otherCarState.TransitionTimestamp).Seconds() * 1000
	var otherDistanceTravelled = CalculateDistanceTravelled(otherCarState.CarSpeed, otherTimeDelta)
	var distanceInTimeStep = CalculateDistanceTravelled(currentCarState.CarSpeed*1.5, 200)

	var nextTileNo = (currentCarState.PosTileNo+1) % currentCarState.MaxTileNo

	// Other car must be on same lane and on next or current tile
	if otherCarState.LaneNo == laneNo && (
		otherCarState.PosTileNo == currentCarState.PosTileNo ||
			otherCarState.PosTileNo == nextTileNo) {
		var distanceDelta float64 = 0
		//1. Check if cars are on same tiles
		if otherCarState.PosTileNo == currentCarState.PosTileNo &&
			otherDistanceTravelled > currentDistanceTravelled{
			distanceDelta = otherDistanceTravelled - currentDistanceTravelled

		} else if otherCarState.PosTileNo == nextTileNo {
			//2. Check if other car is on next tile
			distanceDelta = otherDistanceTravelled +
				(float64(currentCarState.LaneLength) - currentDistanceTravelled)
		}

		mlog.Printf("DEBUG: Car pos: tile %d, pos: %f", currentCarState.PosTileNo, currentDistanceTravelled)
		mlog.Printf("DEBUG: Car pos other: tile %d, pos: %f", otherCarState.PosTileNo, otherDistanceTravelled)
		mlog.Printf("DEBUG: Lane length %d", currentCarState.LaneLength)
		mlog.Printf("DEBUG: lane distance delta %f", float64(currentCarState.LaneLength) - currentDistanceTravelled)
		mlog.Println("DEBUG: Distance is ", distanceDelta)

		// Check if distance is enough in respect to speed
		// It is not enough if the distanceDelta is lower than
		// What the car would travel in 1.5 intervals
		if distanceDelta < distanceInTimeStep {
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
	speed := calculateSpeed(carNo, track)

	//Change speed according to car before
	cmd := anki.Command{ CarNo: carNo, Command: "s", Param1: string(speed)}
	cmdCh <- cmd
	return true
}

func calculateSpeed(carNo int, track []anki.Status) int {
	var currentCarState = getStateForCarNo(carNo, track)
	var blockingCarState anki.Status

	//Check all other car states
	for index, otherCarState := range track {
		if index > carNo {
			if !otherCarState.MsgTimestamp.IsZero() && index != carNo && otherCarState.LaneNo == currentCarState.LaneNo &&
				otherCarState.PosTileNo == currentCarState.PosTileNo {
				mlog.Printf("WARNING: Other car in front")
				blockingCarState = otherCarState
				return blockingCarState.CarSpeed
			}
		}

	}
	return 0
}

/**
Simply drive on
 */
func driveAhead(carNo int, track []anki.Status, cmdCh chan anki.Command) {
	mlog.Printf("INFO: Drive ahead")
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