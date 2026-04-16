package database

import "context"

// txKey 用于在 context 中存储 Ent 事务
type txKey struct{}

// ExtractTx 从 context 中提取事务客户端，若不在事务中则返回原始 client
// 注意：此函数需要配合具体 Ent client 使用，实际使用时需类型断言或泛型封装
func ExtractTx(ctx context.Context, client any) any {
	if tx := ctx.Value(txKey{}); tx != nil {
		return tx
	}
	return client
}

func WithTx(ctx context.Context, tx any) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func GetTx(ctx context.Context) (any, bool) {
	v := ctx.Value(txKey{})
	return v, v != nil
}
