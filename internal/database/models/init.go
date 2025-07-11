package models

type ModelList []interface{}

func All() ModelList {
	return ModelList{
		&Users{},
		&UserStatus{},
		&UserLocation{},
		&UserCredentials{},
	}
}