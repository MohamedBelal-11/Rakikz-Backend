package slices

func Concat[T any](slices ...[]T) []T {
    var result []T
    for _, s := range slices {
        result = append(result, s...)
    }
    return result
}

func Map[T any, R any](items []T, fn func(T) R) []R {
    result := make([]R, len(items))

    for i, item := range items {
        result[i] = fn(item)
    }

    return result
}