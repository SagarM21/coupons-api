package db

import "e-commerce/pkg/models"


func (d *DB) SaveAgeEstimate(c models.AgeEstimate) error {
	return d.Session.Query(
		`INSERT INTO age_estimates (name, age, count) VALUES (?, ?, ?)`,
		c.Name, c.Age, c.Count,
	).Exec()
}

func (d *DB) GetAgeEstimate(name string) (*models.AgeEstimate, error){
	query := `SELECT name, age, count FROM age_estimates WHERE name = ? LIMIT 1`

	var res models.AgeEstimate
	if err := d.Session.Query(query, name).Scan(&res.Name, &res.Age, &res.Count); err != nil {
		return nil, err
	}
	return &res, nil
}