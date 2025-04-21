package blacklist

type BlacklistApi interface {	
	Add(jti, ip string) error
	Delete(jti string) error
	ContainsAndGetIp(jti string) (string, bool)
}

type Blacklist struct {
	BlacklistApi
}

func New() *Blacklist {
	return &Blacklist{
		BlacklistApi: newCache(),
	}
}

