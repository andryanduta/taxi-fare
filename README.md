# Taxi Fare

A program to evaluate and calculate taxi fare based on input Time and Distance. Program will read user's input via terminal.

## Prerequisites

To follow the example in this article, you will need:
- A Go workspace set up 

note: I use go version v1.14.4 when developing the program

## Requirements

The base fare is 400 yen for up to 1 km.
Up to 10 km, 40 yen is added every 400 meters.
Over 10km, 40 yen is added every 350 meters.
This taxi is equipped with the following two meters. Only one of the most recent real values is recorded on these meters.
- Distance Meter
- Fare Meter

The specifications of the distance meter are as follows.
- Space-separated first half is elapsed time (Max 99:99:99.999), second half is mileage.(the unit is meters, Max: 99999999.9)
- It keeps only the latest values.
- It calculates and creates output of the mileage per minute, but an error of less than 1 second occurs.

Error Definition
Error occurs under the following conditions.
- Not under the format of, hh:mm:ss.fff<SPACE>xxxxxxxx.f<LF>, but under an improper format.
- Blank Line
- When the past time has been sent.
- The interval between records is more than 5 minutes apart.
- When there are less than two lines of data.
- When the total mileage is 0.0m.

## How to run

Check out the guideline to run the program
```
Available commands for taxi-fare:

Usage:

    go build                                Compile the project
    ./taxi-fare                             Run binary
    go test ./...                           Run tests
```

### Input

```
elapsed_time<SPACE>mileage
```
Input contains two main object. Check below format as an example. Press enter to terminate the input and processing the output. 
```
00:00:00.000 0.0
00:03:00.123 377.5
00:05:00.128 1111.1
<PRESS ENTER>
```

### Output

Taxi fare output 
```
411.11
```