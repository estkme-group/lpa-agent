package lpac

type ProfileState uint8

const (
	ProfileStateDisabled ProfileState = 0
	ProfileStateEnabled  ProfileState = 1
)

type ProfileClass uint8

const (
	ProfileClassTest         ProfileClass = 0
	ProfileClassProvisioning ProfileClass = 1
	ProfileClassOperational  ProfileClass = 2
)
