package database

import (
	sqlxadapter "github.com/Blank-Xu/sqlx-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/jmoiron/sqlx"
)

func NewEnforcer(db *sqlx.DB) (*casbin.Enforcer, error) {
	a, err := sqlxadapter.NewAdapter(db, "casbin")
	if err != nil {
		return nil, err
	}

	// Initialize the model from Go code.
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act")

	// Create the enforcer.
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	if err = e.LoadPolicy(); err != nil {
		return nil, err
	}

	return e, nil
}
