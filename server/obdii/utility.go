package obdii

import (
	"log"
	"strconv"
	"strings"
)

// Convert two bytes into a DTC
//
// Reference: https://en.wikipedia.org/wiki/OBD-II_PIDs#Service_03_(no_PID_required)
func BytesToDtc(byteA, byteB byte) string {
	firstCharacter := "?"

	switch (0b1100_0000 & byteA) >> 6 {
	case 0b00:
		firstCharacter = "P" // Powertrain
	case 0b01:
		firstCharacter = "C" // Chassis
	case 0b10:
		firstCharacter = "B" // Body
	case 0b11:
		firstCharacter = "U" // Network
	}

	secondCharacter := strconv.FormatUint(uint64((0b0011_0000&byteA)>>4), 16)
	thirdCharacter := strconv.FormatUint(uint64(0b0000_1111&byteA), 16)
	fourthCharacter := strconv.FormatUint(uint64((0b1111_0000&byteB)>>4), 16)
	fifthCharacter := strconv.FormatUint(uint64(0b0000_1111&byteB), 16)

	troubleCode := strings.ToUpper(firstCharacter + secondCharacter + thirdCharacter + fourthCharacter + fifthCharacter)
	log.Println("Decoded DTC: " + troubleCode)
	return troubleCode
}
