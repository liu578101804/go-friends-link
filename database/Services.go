package database

// GetAllFriends 获取所有的友链
func GetAllFriends() ([]*Friends, error) {
	friends := make([]*Friends, 0)
	err := D.Find(&friends).Error
	return friends, err
}
