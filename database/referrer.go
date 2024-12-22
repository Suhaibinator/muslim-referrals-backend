package database

import "fmt"

func (db *DbDriver) CreateReferrer(record *Referrer) (*Referrer, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if err := db.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (db *DbDriver) UpdateReferrer(userId uint64, record *Referrer) (*Referrer, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Ensure the referrer belongs to the user
	if record.UserId != userId {
		return nil, fmt.Errorf("unauthorized to update this referrer")
	}

	// Save updates to the referrer if it belongs to the user
	result := db.db.Model(record).Where("user_id = ?", userId).Save(record)
	if result.Error != nil {
		return nil, result.Error // Handle the error, could be due to a database issue
	}

	// Fetch the updated record
	var updatedRecord Referrer
	if err := db.db.Where("referrer_id = ? AND user_id = ?", record.ReferrerId, userId).First(&updatedRecord).Error; err != nil {
		return nil, err
	}

	return &updatedRecord, nil
}

func (db *DbDriver) DeleteReferrer(userId uint64, record *Referrer) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Ensure the referrer belongs to the user
	if record.UserId != userId {
		return fmt.Errorf("unauthorized to delete this referrer")
	}

	// Delete the referrer if it belongs to the user
	result := db.db.Where("user_id = ?", userId).Delete(record)
	if result.Error != nil {
		return result.Error // Handle the error, could be due to a database issue
	}

	return nil
}

func (db *DbDriver) GetReferrerById(id uint64) *Referrer {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var referrer Referrer
	db.db.Preload("User").First(&referrer, id)
	return &referrer
}

func (db *DbDriver) GetReferrerByUserId(userId uint64) *Referrer {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var referrer Referrer
	db.db.Preload("User").Where("user_id = ?", userId).First(&referrer)
	return &referrer
}
