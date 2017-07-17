package user_init

import (
	"strings"
	"chess/api/models"
)

func DeviceInit(uId int, from, uniqueId, idfv, idfa string) error {

	from = strings.ToLower(from)
	T := 1
	if strings.Contains(from, "ios") || strings.Contains(from, "android") {
		if strings.Contains(from, "ios") {
			T = 2
		}
		device := &models.DeviceModel{
			UserId:   uId,
			UniqueId: uniqueId,
			Type:     T,
			Idfa:     idfa,
			Idfv:     idfv,
		}
		err := device.Upsert()

		return err
	}
	return nil

}
