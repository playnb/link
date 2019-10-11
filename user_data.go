package link

type UserData struct {
	userData interface{}
}

func (ud *UserData) GetUserData() interface{} {
	return ud.userData
}
func (ud *UserData) SetUserData(data interface{}) {
	ud.userData = data
}
