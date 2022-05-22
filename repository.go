package good

import (
	"github.com/overflowingd/good/repository"
	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

func (r *Repository) Create(entities any) (repository.RowsAffected, error) {
	result := r.Db.Create(entities)
	return repository.RowsAffected(result.RowsAffected), result.Error
}

func (r *Repository) Save(entities any) (repository.RowsAffected, error) {
	result := r.Db.Save(entities)
	return repository.RowsAffected(result.RowsAffected), result.Error
}

func (r *Repository) First(dest any, conditions ...any) (repository.RowsAffected, error) {
	result := r.Db.First(dest, conditions...)
	return repository.RowsAffected(result.RowsAffected), result.Error
}

func (r *Repository) FirstByID(dest any, id any) (repository.RowsAffected, error) {
	return r.First(dest, "id = ?", id)
}

func (r *Repository) Take(dest any, conditions ...any) (repository.RowsAffected, error) {
	result := r.Db.Take(dest, conditions...)
	return repository.RowsAffected(result.RowsAffected), result.Error
}

func (r *Repository) TakeByID(dest any, id any) (repository.RowsAffected, error) {
	return r.Take(dest, "id = ?", id)
}

func (r *Repository) Last(dest any, conditions ...any) (repository.RowsAffected, error) {
	result := r.Db.Last(dest, conditions...)
	return repository.RowsAffected(result.RowsAffected), result.Error
}

func (r *Repository) LastByID(dest any, id any) (repository.RowsAffected, error) {
	return r.Last(dest, "id = ?", id)
}
