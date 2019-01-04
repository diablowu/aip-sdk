package unit

type UnitSkillService struct {
}

func (unit UnitSkillService) Say() string {
	return AuthService{}.Token(API_URL)
}
