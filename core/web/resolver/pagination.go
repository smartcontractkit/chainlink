package resolver

type PaginationMetadataResolver struct {
	total int32
}

func NewPaginationMetadata(total int32) *PaginationMetadataResolver {
	return &PaginationMetadataResolver{total: total}
}

func (r *PaginationMetadataResolver) Total() int32 {
	return r.total
}
