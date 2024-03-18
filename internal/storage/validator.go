package storage

func (f *Film) Validate() error {
	if len(f.Title) < 1 || len(f.Title) > 150 || len(f.Description) > 1000 || f.Rating < 0 || f.Rating > 10 {
		return ErrInvalidData
	}
	return nil
}
