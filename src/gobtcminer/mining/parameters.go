/*
Author: Mathieu Mailhos
Filename: parameters.go
Description: Define mining parameters specific to Bitcoin BP023
*/

package mining

import "strconv"

//Gettarget TODO Calculate the right target depending on the difficulty. This one is totally made up for testing purpose.
func Gettarget(difficulty float64, bits uint32) string {
	const padding int = 17
	var target = ""
	for i := 0; i < padding; i++ {
		target = target + "0"
	}
	target = target + strconv.Itoa(int(uint32(difficulty)*bits))
	for i := len(target); i < 64; i++ {
		target = target + "f"
	}
	return target
}
