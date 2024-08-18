package database

func (db *DbDriver) CreateCompany(record *Company) (*Company, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if err := db.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (db *DbDriver) UpdateCompany(record *Company) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db.Save(record)
}

func (db *DbDriver) DeleteCompany(record *Company) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.db.Delete(record)
}

func (db *DbDriver) GetCompanyById(id uint64) *Company {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var company Company
	db.db.Preload("Domains").First(&company, id)
	return &company
}

func (db *DbDriver) GetAllCompanies() []Company {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var companies []Company
	db.db.Preload("Domains").Find(&companies)
	return companies
}
