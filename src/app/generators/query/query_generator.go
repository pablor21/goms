package query

import (
	"github.com/pablor21/goms/app/models"
	"github.com/pablor21/goms/pkg/database"
	"gorm.io/gen"
)

type Filter interface {
	// SELECT * FROM @@table WHERE @@column=@value
	FilterWithColumn(column string, value string) ([]*gen.T, error)
}

func generateRepositories() {

	g := gen.NewGenerator(gen.Config{
		OutPath: "./app/repositories",                                               // output path
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(database.GetConnection("default").(*database.GormConnection).Conn()) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(models.User{}, models.Asset{}, models.Tag{}, models.AssetLibrary{}, models.AssetFolder{})

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	// g.ApplyInterface(func(Querier) {}, models.User{})

	g.ApplyInterface(func(Filter) {}, models.Tag{})

	// Generate the code
	g.Execute()
}
