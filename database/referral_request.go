package database

func (db *DbDriver) CreateReferralRequest(record *ReferralRequest) (*ReferralRequest, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if err := db.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (db *DbDriver) UpdateReferralRequest(record *ReferralRequest) (*ReferralRequest, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Save the updated record
	result := db.db.Save(record)
	if result.Error != nil {
		return nil, result.Error // Handle the error, could be due to a database issue
	}

	// Fetch the updated record
	var updatedRecord ReferralRequest
	if err := db.db.Where("referral_request_id = ?", record.ReferralRequestId).First(&updatedRecord).Error; err != nil {
		return nil, err
	}

	return &updatedRecord, nil
}

func (db *DbDriver) DeleteReferralRequest(record *ReferralRequest) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.db.Delete(record).Error
}

func (db *DbDriver) GetReferralRequestById(id uint64) *ReferralRequest {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var referralRequest ReferralRequest
	db.db.Preload("Candidate").
		Preload("Candidate.User").
		Preload("Company").
		Preload("Referrer").
		Preload("Referrer.User").
		Preload("JobLinks").
		Preload("Locations").
		First(&referralRequest, id)
	if referralRequest.ReferralRequestId == 0 { // Checking if the referral request was found
		return nil
	}
	return &referralRequest
}

func (db *DbDriver) GetReferralRequestsByReferrerId(referrerId uint64) []ReferralRequest {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var referralRequests []ReferralRequest
	db.db.Preload("Candidate").
		Preload("Candidate.User").
		Preload("Company").
		Preload("Referrer").
		Preload("Referrer.User").
		Preload("JobLinks").
		Preload("Locations").
		Where("referrer_id = ?", referrerId).
		Find(&referralRequests)
	return referralRequests
}

func (db *DbDriver) GetReferralRequestsByCandidateId(candidateId uint64) []ReferralRequest {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var referralRequests []ReferralRequest
	db.db.Preload("Candidate").
		Preload("Candidate.User").
		Preload("Company").
		Preload("Referrer").
		Preload("Referrer.User").
		Preload("JobLinks").
		Preload("Locations").
		Where("candidate_id = ?", candidateId).
		Find(&referralRequests)
	return referralRequests
}

func (db *DbDriver) GetReferralRequestsByCompanyId(companyId uint64) []ReferralRequest {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var referralRequests []ReferralRequest
	db.db.Preload("Candidate").
		Preload("Candidate.User").
		Preload("Company").
		Preload("Referrer").
		Preload("Referrer.User").
		Preload("JobLinks").
		Preload("Locations").
		Where("company_id = ?", companyId).
		Find(&referralRequests)
	return referralRequests
}

func (db *DbDriver) GetReferralRequestByIdAndCandidateId(referralRequestId, candidateId uint64) *ReferralRequest {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var referralRequest ReferralRequest
	db.db.Preload("Candidate").
		Preload("Candidate.User").
		Preload("Company").
		Preload("Referrer").
		Preload("Referrer.User").
		Preload("JobLinks").
		Preload("Locations").
		Where("id = ? AND candidate_id = ?", referralRequestId, candidateId).
		First(&referralRequest)
	if referralRequest.ReferralRequestId == 0 { // Checking if the referral request was found
		return nil
	}
	return &referralRequest
}
