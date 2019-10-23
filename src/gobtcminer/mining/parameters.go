/*
Author: Mathieu Mailhos
Filename: parameters.go
Description: Define mining parameters specific to Bitcoin BP023
*/

package mining

import (
	"strconv"
	"strings"
)

//Gettarget TODO Calculate the right target depending on the difficulty. This one is totally made up for testing purpose.
func Gettarget(difficulty float64, bits uint32) string {
	const padding int = 17
	target := strings.Repeat("0", padding) + strconv.Itoa(int(uint32(difficulty)*bits))
	suffix := strings.Repeat("f", 64-len(target))
	return target + suffix
}
