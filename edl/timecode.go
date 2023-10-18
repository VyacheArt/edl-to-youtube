/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package edl

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Timecode struct {
	Time  time.Time
	Frame int
}

func (t *Timecode) Parse(s string) error {
	//Timecode is in the format of hh:mm:ss:ff
	//where ff is the frame number.
	parts := strings.Split(s, ":")
	if len(parts) != 4 {
		return fmt.Errorf("invalid timecode format: %s", s)
	}

	numbers := make([]int, 4)
	for i, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil {
			return fmt.Errorf("invalid timecode format: %s, position: %d", s, i)
		}

		numbers[i] = n
	}

	hours, minutes, seconds, frames := numbers[0], numbers[1], numbers[2], numbers[3]

	t.Time = time.Date(0, 0, 0, hours, minutes, seconds, 0, time.UTC)
	t.Frame = frames

	return nil
}

func (t *Timecode) String() string {
	return t.Time.Format("15:04:05") + ":" + fmt.Sprintf("%02d", t.Frame)
}
