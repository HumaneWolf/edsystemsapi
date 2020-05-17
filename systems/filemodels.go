package systems

// EDSMSystem is a star system in the EDSM dump.
type EDSMSystem struct {
	ID     int32           `json:"id"`
	ID64   int64           `json:"id64"`
	Name   string          `json:"name"`
	Coords EDSMCoordinates `json:"coords"`
	Date   string          `json:"date"`
}

// EDSMCoordinates are the coordinates of a system in the EDSM dump.
type EDSMCoordinates struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}
