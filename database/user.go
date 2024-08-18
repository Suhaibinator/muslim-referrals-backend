package database

func (db *DbDriver) CreateUser(record *User) (*User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if err := db.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (db *DbDriver) UpdateUser(record *User) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db.Save(record)
}

func (db *DbDriver) DeleteUser(record *User) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db.Delete(record)
}

func (db *DbDriver) GetUserById(id uint64) *User {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var user User
	db.db.First(&user, id)
	return &user
}
