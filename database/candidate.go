package database

import "fmt"

func (db *DbDriver) CreateCandidate(record *Candidate) (*Candidate, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if err := db.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (db *DbDriver) UpdateCandidate(userId uint64, record *Candidate) (*Candidate, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if record.UserId != userId {
		return nil, fmt.Errorf("unauthorized to update this candidate")
	}

	result := db.db.Model(record).Where("user_id = ?", userId).Save(record)
	if result.Error != nil {
		return nil, result.Error
	}

	var updatedRecord Candidate
	if err := db.db.Where("candidate_id = ? AND user_id = ?", record.CandidateId, userId).First(&updatedRecord).Error; err != nil {
		return nil, err
	}

	return &updatedRecord, nil
}

func (db *DbDriver) DeleteCandidate(userId uint64, record *Candidate) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if record.UserId != userId {
		return fmt.Errorf("unauthorized to delete this candidate")
	}

	result := db.db.Where("user_id = ?", userId).Delete(record)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DbDriver) GetCandidateById(userId, id uint64) *Candidate {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var candidate Candidate
	db.db.Preload("User").Where("candidate_id = ? AND user_id = ?", id, userId).First(&candidate)
	if candidate.CandidateId == 0 {
		return nil
	}
	return &candidate
}

func (db *DbDriver) GetBulkCandidatesByIds(userId uint64, ids []uint64) *[]Candidate {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var candidates []Candidate
	db.db.Preload("User").Where("user_id = ? AND candidate_id IN ?", userId, ids).Find(&candidates)
	return &candidates
}

func (db *DbDriver) GetCandidateByUserId(userId uint64) *Candidate {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var candidate Candidate
	db.db.Preload("User").Where("user_id = ?", userId).First(&candidate)
	if candidate.CandidateId == 0 {
		return nil
	}
	return &candidate
}
