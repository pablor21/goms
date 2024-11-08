// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package repositories

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q            = new(Query)
	Asset        *asset
	AssetFolder  *assetFolder
	AssetLibrary *assetLibrary
	Tag          *tag
	User         *user
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Asset = &Q.Asset
	AssetFolder = &Q.AssetFolder
	AssetLibrary = &Q.AssetLibrary
	Tag = &Q.Tag
	User = &Q.User
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:           db,
		Asset:        newAsset(db, opts...),
		AssetFolder:  newAssetFolder(db, opts...),
		AssetLibrary: newAssetLibrary(db, opts...),
		Tag:          newTag(db, opts...),
		User:         newUser(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Asset        asset
	AssetFolder  assetFolder
	AssetLibrary assetLibrary
	Tag          tag
	User         user
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:           db,
		Asset:        q.Asset.clone(db),
		AssetFolder:  q.AssetFolder.clone(db),
		AssetLibrary: q.AssetLibrary.clone(db),
		Tag:          q.Tag.clone(db),
		User:         q.User.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:           db,
		Asset:        q.Asset.replaceDB(db),
		AssetFolder:  q.AssetFolder.replaceDB(db),
		AssetLibrary: q.AssetLibrary.replaceDB(db),
		Tag:          q.Tag.replaceDB(db),
		User:         q.User.replaceDB(db),
	}
}

type queryCtx struct {
	Asset        IAssetDo
	AssetFolder  IAssetFolderDo
	AssetLibrary IAssetLibraryDo
	Tag          ITagDo
	User         IUserDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Asset:        q.Asset.WithContext(ctx),
		AssetFolder:  q.AssetFolder.WithContext(ctx),
		AssetLibrary: q.AssetLibrary.WithContext(ctx),
		Tag:          q.Tag.WithContext(ctx),
		User:         q.User.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
