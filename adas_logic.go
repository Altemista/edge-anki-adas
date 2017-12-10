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
	ticker := time.NewTicker(5000 * 1e6) // 1e6 = ms, 1e9 = s
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

func canDriveOn(carNo int, track []anki.Status, cmdCh chan anki.Command) bool {
	return false
}

func canChangeLeft(carNo int, track []anki.Status, cmdCh chan anki.Command) bool {
	return false
}

func canChangeRight(carNo int, track []anki.Status, cmdCh chan anki.Command) bool {
	return false
}

func adjustSpeed(carNo int, track []anki.Status, cmdCh chan anki.Command) bool {
	return false
}

func driveAhead(carNo int, track []anki.Status, cmdCh chan anki.Command) {
}

func changeToLeftLane(carNo int, track []anki.Status, cmdCh chan anki.Command) {
}

func changeToRightLane(carNo int, track []anki.Status, cmdCh chan anki.Command) {
}
