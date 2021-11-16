//Package builder はSQLのAND, OR, WHERE句を適切に組み合わせるためのビルドツール
package builder

type Builder interface {
	Build() string
}
