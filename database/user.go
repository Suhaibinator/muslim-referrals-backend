package database

import "gorm.io/gorm"

func (db *DbDriver) CreateUser(record *User) (*User, error) {
	record.Id = 0
	db.mu.Lock()
	defer db.mu.Unlock()
	if err := db.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (db *DbDriver) UpdateUser(record *User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	result := db.db.Save(record)
	if result.Error != nil {
		return result.Error
	}

	return nil
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

func (db *DbDriver) GetUserByEmail(email string) *User {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var user User

	result := db.db.Where("email = ?", email).First(&user)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return nil
	}
	return &user
}
