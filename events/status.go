package events

import (
	"github.com/CmdrVasquess/watched"
)

type StatusFlag uint32

type statusT string

const StatusEvent = statusT(watched.StatStatusName)

func (t statusT) New() Event     { return new(Status) }
func (t statusT) String() string { return string(t) }

type Status struct {
	Common
	Flags     StatusFlag
	Pips      [3]int
	FireGroup int
	Fuel      struct {
		Main      float64 `json:"FuelMain"`
		Reservoir float64 `json:"FuelReservoi"`
	}
}

func (_ *Status) EventType() Type { return StatusEvent }

// Read: https://forums.frontier.co.uk/forums/elite-api-and-tools/
const (
	FStatusDocked StatusFlag = (1 << iota)
	FStatusLanded
	FStatusGearDown
	FStatusShieldsUp
	FStatusSupercruise

	FStatusFAOff
	FStatusHPDeployed
	FStatusInWing
	FStatusLightsOn
	FStatusCSDeployed

	FStatusSilentRun
	FStatusFuelScooping
	FStatusSrvHandbrake
	FStatusSrvTurret
	FStatusSrvUnderShip

	FStatusSrvDriveAssist
	FStatusFsdMassLock
	FStatusFsdCharging
	FStatusCooldown
	FStatusLowFuel

	FStatusOverHeat
	FStatusHasLatLon
	FStatusIsInDanger
	FStatusInterdicted
	FStatusInMainShip

	FStatusInFighter
	FStatusInSrv
	FStatusHudAnalysis
	FStatusNightVis
	FStatusAltAvgR
)

func (s *Status) AnyFlag(fs StatusFlag) bool {
	return s.Flags&fs > 0
}

func (s *Status) AllFlags(fs StatusFlag) bool {
	return s.Flags&fs == fs
}

func init() {
	RegisterType(string(StatusEvent), StatusEvent)
}
