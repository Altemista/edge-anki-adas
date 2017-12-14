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
	ticker := time.NewTicker(500 * 1e6) // 1e6 = ms, 1e9 = s
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

	//TODO: Iterative approach to find most left lane etc.
	//TODO: Zero timestmap

	if !getStateForCarNo(carNo, track).MsgTimestamp.IsZero() {
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
		if !otherCarState.MsgTimestamp.IsZero() && index != carNo && otherCarState.LaneNo == currentCarState.LaneNo &&
			otherCarState.PosTileNo == currentCarState.PosTileNo {
			mlog.Printf("WARNING: Cars are on same lane and tile")
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

	//Check all other car states
	for index, otherCarState := range track {
		if !otherCarState.MsgTimestamp.IsZero() && index != carNo && otherCarState.LaneNo == currentCarState.LaneNo+1 &&
			otherCarState.PosTileNo == currentCarState.PosTileNo {
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
			if !otherCarState.MsgTimestamp.IsZero() && index != carNo && otherCarState.LaneNo == currentCarState.LaneNo-1 &&
				otherCarState.PosTileNo == currentCarState.PosTileNo {
				mlog.Printf("WARNING: Other car on right lane, no change possible")
				return false
			}
		}
	}
	mlog.Printf("DEBUG: Can change right")
	return true
}

/**
CRITERIA:
Find car that is blocking us and adjust speed to the speed of the blocking car
 */
func adjustSpeed(carNo int, track []anki.Status, cmdCh chan anki.Command) bool {
	var currentCarState = getStateForCarNo(carNo, track)
	var blockingCarState anki.Status

	//Check all other car states
	for index, otherCarState := range track {
		if index > carNo {
			if !otherCarState.MsgTimestamp.IsZero() && index != carNo && otherCarState.LaneNo == currentCarState.LaneNo && otherCarState.PosTileNo == currentCarState.PosTileNo {
				mlog.Printf("WARNING: Other car in front")
				blockingCarState = otherCarState
			}
		}

	}

	//Change speed according to car before
	cmd := anki.Command{ CarNo: carNo, Command: "s", Param1: string(blockingCarState.CarSpeed)}
	cmdCh <- cmd
	return true
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