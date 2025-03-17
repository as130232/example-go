package deviceType

const (
	Web     = "Web"
	Client  = "Client"
	Ios     = "Ios"
	Android = "Android"
)

func IsMobile(deviceType string) bool {
	if Ios == deviceType || Android == deviceType {
		return true
	}

	return false
}

func IsDeviceType(deviceType string) bool {
	if Web == deviceType || Client == deviceType || Ios == deviceType || Android == deviceType {
		return true
	}

	return false
}
