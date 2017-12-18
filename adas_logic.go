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

	if !getStateForCarNo(carNo, track).MsgTimestamp.IsZero() {
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
	var nextTileNo = (currentCarState.PosTileNo+1) % currentCarState.MaxTileNo

	//Check all other car states
	for index, otherCarState := range track {
		if !otherCarState.MsgTimestamp.IsZero() && index != carNo &&
			hasCarInFront(otherCarState, currentCarState.LaneNo, nextTileNo) &&
				otherCarState.CarSpeed < currentCarState.CarSpeed {
			mlog.Printf("WARNING: Other car on same lane next  tile")
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
func getAvailableLane(carNo int, track []anki.Status, cmdCh chan anki.Command) int {
	var currentCarState = getStateForCarNo(carNo, track)
	var nextTileNo = (currentCarState.PosTileNo+1) % currentCarState.MaxTileNo

	//Check for all possible lanes
	var suggestedLaneIndex = -1
	var laneAvailable = false

	for _, laneOffset := range []int {1, -1, 2, -2, 3, -3} {
		laneAvailable = true
		suggestedLaneIndex = currentCarState.LaneNo + laneOffset

		if suggestedLaneIndex < 1 || suggestedLaneIndex > 4 {
			continue
		}

		//Check all other car states of same and next tile
		for index, otherCarState := range track {
			if !otherCarState.MsgTimestamp.IsZero() && index != carNo &&
				(hasCarInFront(otherCarState, suggestedLaneIndex, currentCarState.PosTileNo) ||
						hasCarInFront(otherCarState, suggestedLaneIndex, nextTileNo)){
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

func hasCarInFront(otherCarState anki.Status, laneNo int, tileNo int) bool {
	if otherCarState.LaneNo == laneNo && otherCarState.PosTileNo == tileNo {
		return true
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
			if !otherCarState.MsgTimestamp.IsZero() && index != carNo && otherCarState.LaneNo == currentCarState.LaneNo && otherCarState.PosTileNo == currentCarState.PosTileNo {
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