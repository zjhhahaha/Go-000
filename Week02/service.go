package demo

type service struct {
	dao Dao
}

func (s *service) GetUser() (*User, error) {
	user, err := s.dao.GetUser()
	if err != nil {
		if IsNotFound(err) {
			//mock
			return &User{}, nil
		}
		return nil, err
	}
	return user, nil
}
