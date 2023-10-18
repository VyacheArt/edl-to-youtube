/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package edl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	// TransitionCut is the cut transition.
	TransitionCut Transition = "C"

	// TransitionDissolve is the dissolve transition.
	TransitionDissolve Transition = "D"
)

type (
	Clip struct {
		ClipNumber  int
		TrackNumber string
		Transition  Transition
		SourceIn    Timecode
		SourceOut   Timecode
		RecordIn    Timecode
		RecordOut   Timecode
		Color       string
		Marker      string
	}

	Transition string
)

func (c *Clip) Parse(s string) error {
	//we can receive a single clip header or meta information at the second line
	lines := strings.Split(s, "\n")

	if err := c.parseHeader(lines[0]); err != nil {
		return err
	}

	return nil
}

func (c *Clip) parseHeader(headerLine string) error {
	re := regexp.MustCompile(`\s{2,}`)
	normalizedHeaderLine := strings.TrimSpace(re.ReplaceAllString(headerLine, " "))

	//Clip is in the format of
	//001 V C        00:00:00:00 00:00:00:00 00:00:00:00 00:00:00:00
	parts := strings.Split(normalizedHeaderLine, " ")
	if len(parts) != 7 && len(parts) != 8 {
		return fmt.Errorf("invalid clip format: %s", normalizedHeaderLine)
	}

	clipNumber, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid clip format: %s, position: 0", normalizedHeaderLine)
	}

	//sometimes there is an audio track number after the clip number
	//for example after export in DaVinci Resolve
	offset := 0
	if len(parts) == 8 {
		offset = 1
	}

	c.ClipNumber = clipNumber
	c.TrackNumber = parts[1+offset]
	c.Transition = Transition(parts[2+offset])

	if err := c.SourceIn.Parse(parts[3+offset]); err != nil {
		return err
	}

	if err := c.SourceOut.Parse(parts[4+offset]); err != nil {
		return err
	}

	if err := c.RecordIn.Parse(parts[5+offset]); err != nil {
		return err
	}

	if err := c.RecordOut.Parse(parts[6+offset]); err != nil {
		return err
	}

	return nil
}

func (c *Clip) parseMeta(s string) error {
	parts := strings.Split(s, "|")
	pairs := map[string]func(string){
		"C:": func(s string) { c.Color = s },
		"M:": func(s string) { c.Marker = s },
	}

partLoop:
	for _, part := range parts {
		for key, set := range pairs {
			if strings.HasPrefix(part, key) {
				set(strings.TrimSpace(strings.TrimPrefix(part, key)))
				continue partLoop
			}
		}
	}

	return nil
}
