package main

import (
	"errors"
	"fmt"
)

const EARTH_HOURTOSEC = 3600;
const EARTH_MINUTETOSEC = 60;
const ROKETIN_HOURTOSEC = 10000;
const ROKETIN_MINUTETOSEC = 100;

func main() {
	var hours int32;
	var minutes int32;
	var seconds int32;
	for {
		fmt.Print("Enter hours, minutes, seconds: ");
		_, errorUserInput := fmt.Scanln(&hours, &minutes, &seconds);
		
		if errorUserInput != nil {
        	fmt.Println("\nInvalid input! Please enter an integer.")
        	break;
    	}

		var newHours, newMinutes, newSeconds, err = convertTime(hours, minutes, seconds);
		if err != nil {
			fmt.Println("User input is not valid. Please use the correct format.");
			break;
		} else {
			fmt.Printf("Time in Roketin are %02d:%02d:%02d", newHours, newMinutes, newSeconds);
		}
		
	}
}

func convertTime(hours int32, minutes int32, seconds int32) (uint8, uint8, uint8, error) {
	if seconds >= 60 || seconds < 0 {
		return 0, 0, 0, errors.New("seconds value must be in the range 0 to 60");
	}
	if minutes >= 60 || minutes < 0 {
		return 0, 0, 0, errors.New("minutes value must be in the range 0 to 60");;
	}
	if hours >= 24 || hours < 0 {
		return 0, 0, 0, errors.New("hours value must be in the range 0 to 24");;
	}

	var totalTime uint32 = uint32(hours * EARTH_HOURTOSEC + minutes * EARTH_MINUTETOSEC + seconds);

	var newHours uint8 = uint8(totalTime / ROKETIN_HOURTOSEC);
	totalTime %= ROKETIN_HOURTOSEC;
	var newMinutes uint8 = uint8(totalTime / ROKETIN_MINUTETOSEC);
	totalTime %= ROKETIN_MINUTETOSEC
	var newSeconds uint8 = uint8((totalTime));
	
	return newHours, newMinutes, newSeconds, nil;
}