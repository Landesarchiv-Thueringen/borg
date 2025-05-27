package internal

func GetLocalization() map[string]string {
	return serverConfig.Localization.Dict
}
