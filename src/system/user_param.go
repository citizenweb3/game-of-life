package system

import "math/rand"

type UsersParam struct {
	Volts      uint64
	Amperes    uint64
	Cyberlinks uint64
	Kw         uint64
	// Karma uint64
	// AverageRankWeight uint64
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
		Volts:      rand.Uint64() % 1000,
		Amperes:    rand.Uint64() % 1000,
		Cyberlinks: rand.Uint64() % 1000,
		Kw:         rand.Uint64() % 1000,
	}
}
