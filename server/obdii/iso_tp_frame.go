package obdii

import (
	"log"
	"strconv"

	"github.com/tntmeijs/pitona/server/utility"
)

type frameType int

const (
	invalid frameType = iota - 1
	single
	first
	consecutive
	flowControl
)

// Represents a single frame of the ISO 15765-2 standard
//
// Reference: https://en.wikipedia.org/wiki/ISO_15765-2
type isoTpFrame struct {
	frameType frameType
}

func (isoTpFrame *isoTpFrame) parse(bytes ...byte) error {
	if len(bytes) == 0 {
		return utility.GenericErrorMessage{Message: "No bytes passed to frame - unable to parse frame"}
	}

	log.Println(bytes[0])

	// Bytes 7 through 4 identify the frame's type
	switch frameType((bytes[0] & 0b1111_0000) >> 4) {
	case single:
		isoTpFrame.frameType = single
	case first:
		isoTpFrame.frameType = first
	case consecutive:
		isoTpFrame.frameType = consecutive
	case flowControl:
		isoTpFrame.frameType = flowControl
	default:
		isoTpFrame.frameType = invalid
	}

	log.Println("Frame type " + strconv.Itoa(int(isoTpFrame.frameType)))

	return nil
}
