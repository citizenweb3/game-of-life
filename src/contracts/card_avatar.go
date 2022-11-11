package contracts

import "net/url"

type CardAvatarI interface {
	GetAvatarLink() url.URL
}

type CardAvatarsI interface {
	GenerateAvatar() CardAvatarI
}

type CardAvatar struct {
	Link url.URL
}
