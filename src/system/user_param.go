package system

import (
	"gameoflife/utils"
)

type UsersParam struct {
	Volts          uint64
	Amperes        uint64
	Cyberlinks     uint64
	Kw             uint64
	Hydrogen       uint64
	LockedHydrogen uint64

	// Karma uint64
	// AverageRankWeight uint64
}

func (up *UsersParam) LockHydogen(count uint64) error {
	if count > up.Hydrogen {
		return utils.ErrSystemNotEnoughH
	}
	up.LockedHydrogen += count
	up.Hydrogen -= count
	return nil
}

func (up *UsersParam) GetLockedHydrogen() uint64 {
	return up.LockedHydrogen
}

func (up *UsersParam) GetHydrogen() uint64 {
	return up.Hydrogen
}

func (up *UsersParam) GetVolts() uint64 {
	return up.Volts
}

func (up *UsersParam) GetAmperes() uint64 {
	return up.Amperes
}

func (up *UsersParam) GetCountCyberlinks() uint64 {
	return up.Cyberlinks
}

func (up *UsersParam) GetKw() uint64 {
	return up.Kw
}

func GenerateRandomUserParam() UsersParam {
	return UsersParam{
		Volts:          utils.GetRandomNumberUint64(1000),
		Amperes:        utils.GetRandomNumberUint64(1000),
		Cyberlinks:     utils.GetRandomNumberUint64(1000),
		Kw:             utils.GetRandomNumberUint64(1000),
		Hydrogen:       utils.GetRandomNumberUint64(1000),
		LockedHydrogen: utils.GetRandomNumberUint64(1000),
	}
}
