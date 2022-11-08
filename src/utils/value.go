package utils

type UserID string
type CardID string

func (ui *UserID) ToString() string {
	return string(*ui)
}
func (ci *CardID) ToString() string {
	return string(*ci)
}
