package main

import tests "MapsBenchmark/internal"

func main() {

	tests.InitMaps()
	tests.MapsAddTest()
	tests.MapsGetTest()
	//tests.MapsGetAndAddTest()
}
