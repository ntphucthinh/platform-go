package dbutil

import (
	api "platform-go/generated/api"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

// Paginator ページネーションの結果を保持する構造体
type Paginator[T any] struct {
	Data          []T               `json:"data"`          // 取得したレコードデータ
	PaginatorInfo api.PaginatorInfo `json:"paginatorInfo"` // ページネーション情報
}

const (
	minPage     = 1  // 最小ページ番号
	maxPageSize = 50 // 1ページあたりの最大レコード数
)

var (
	ErrInvalidPageSize = errors.Newf("invalid page size: must be between 1 and %d", maxPageSize)
	ErrInvalidPage     = errors.New("invalid page number: must be greater than 0")
)

// Paginate ページネーション処理を実行する
// paginator, err := dbutil.ModelPaginate[User](db.Model(&User{}), 10, 1)
// パラメータ:
// - query: 実行するクエリ（db.Model(model)の形式）
// - limit: 1ページあたりの取得レコード数
// - page: 取得するページ番号
func ModelPaginate[T any](query *gorm.DB, limit int, page int) (*Paginator[T], error) {
	// バリデーションチェック
	if limit < 1 || limit > maxPageSize {
		return nil, ErrInvalidPageSize
	}
	if page < minPage {
		page = 1
	}

	// オフセット計算
	offset := (page - 1) * limit

	// 総レコード数を取得（キャッシュ可能な場合は利用）
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, err
	}

	// メインクエリ実行
	var data []T
	if err := query.Offset(offset).Limit(limit).Find(&data).Error; err != nil {
		return nil, err
	}

	// 総ページ数を計算（0件の場合は0）
	totalPages := 0
	if totalRecords > 0 {
		totalPages = int((totalRecords + int64(limit) - 1) / int64(limit))
	}

	// 結果を設定
	return &Paginator[T]{
		Data: data,
		PaginatorInfo: api.PaginatorInfo{
			Count:        len(data),
			CurrentPage:  page,
			HasMorePages: page < totalPages,
			LastPage:     totalPages,
			PerPage:      limit,
			Total:        totalRecords,
		},
	}, nil
}
