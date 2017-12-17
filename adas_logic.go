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
	//TODO: Zero timestmap
	var currentCarState = getStateForCarNo(carNo, track)
	if !currentCarState.MsgTimestamp.IsZero() &&
		!currentCarState.TransitionTimestamp.IsZero(){
		mlog.Printf("CarNo %d", carNo)
		if canDriveOn(carNo, track, cmdCh) {
			driveAhead(carNo, track, cmdCh)
		} else {
			if canChangeLeft(carNo, track, cmdCh) {
				driveAhead(carNo, track, cmdCh)
				changeToLeftLane(carNo, track, cmdCh)
			} else {
				if canChangeRight(carNo, track, cmdCh) {
					driveAhead(carNo, track, cmdCh)
					changeToRightLane(carNo, track, cmdCh)
				} else {
					adjustSpeed(carNo, track, cmdCh)
				}
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

/**
CRITERIA:
Check if the car is on the most left lane (4) -> False
Else
	If there is any car on the current lane + 1 and the same tile
		return false
	Else
		return true

 */
func canChangeLeft(carNo int, track []anki.Status, cmdCh chan anki.Command) bool {
	var currentCarState = getStateForCarNo(carNo, track)

	if currentCarState.LaneNo >= 4 {
		return false
	}

	//Check all other car states of same and next tile
	for index, otherCarState := range track {
		if timeStampsValid(otherCarState) && index != carNo &&
			hasCarInFront(otherCarState, currentCarState, currentCarState.LaneNo+1) {
			mlog.Printf("WARNING: Other car on left lane, no change possible")
			return false
		}
	}

	mlog.Printf("DEBUG: Can change left")
	return true
}

/**
CRITERIA:
Check if the car is on the most right lane (1)
Else
	If there is any car on the current lane - 1 and the same tile
		return false
	Else
		return true
 */
func canChangeRight(carNo int, track []anki.Status, cmdCh chan anki.Command) bool {
	var currentCarState = getStateForCarNo(carNo, track)

	if currentCarState.LaneNo <= 1 {
		return false
	}

	//Check all other car states
	for index, otherCarState := range track {
		if index > carNo {
			if timeStampsValid(otherCarState) && index != carNo &&
				hasCarInFront(otherCarState, currentCarState, currentCarState.LaneNo-1) {
				mlog.Printf("WARNING: Other car on right lane, no change possible")
				return false
			}
		}
	}
	mlog.Printf("DEBUG: Can change right")
	return true
}

func timeStampsValid(carState anki.Status) bool {
	return !carState.MsgTimestamp.IsZero() && !carState.TransitionTimestamp.IsZero()
}

func hasCarInFront(otherCarState anki.Status, currentCarState anki.Status, laneNo int) bool {
	var currentTimeDelta = time.Since(currentCarState.TransitionTimestamp).Seconds() * 1000
	var currentDistanceTravelled = CalculateDistanceTravelled(currentCarState.CarSpeed, currentTimeDelta)
	var otherTimeDelta = time.Since(otherCarState.TransitionTimestamp).Seconds() * 1000
	var otherDistanceTravelled = CalculateDistanceTravelled(otherCarState.CarSpeed, otherTimeDelta)
	var distanceInTimeStep = CalculateDistanceTravelled(currentCarState.CarSpeed, 200)

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

		mlog.Println("DEBUG: Distance is ", distanceDelta)

		// Check if distance is enough in respect to speed
		// It is not enough if the distanceDelta is lower than
		// What the car would travel in 1.5 intervals
		if distanceDelta < distanceInTimeStep*1.5 {
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
func changeToLeftLane(carNo int, track []anki.Status, cmdCh chan anki.Command) {
	mlog.Printf("INFO: Changing to left lane")
	cmd := anki.Command{ CarNo: carNo, Command: "c", Param1: "3"}
	cmdCh <- cmd

	mlog.Printf("Command sent %+v\n", cmd)
}

/**
Initiate right change
 */
func changeToRightLane(carNo int, track []anki.Status, cmdCh chan anki.Command) {
	mlog.Printf("INFO: Changing to right lane")
	cmd := anki.Command{ CarNo: carNo, Command: "c", Param1: "", Param2: "right"}
	cmdCh <- cmd

	mlog.Printf("Command sent %+v\n", cmd)
}

func getStateForCarNo(carNo int, track []anki.Status) anki.Status {
	return track[carNo]
}