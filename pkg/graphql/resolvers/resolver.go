package resolvers

type Resolver struct{}

func (r *Resolver) Query() *QueryResolver {
	return &QueryResolver{}
}

func (r *Resolver) Mutation() *MutationResolver {
	return &MutationResolver{}
}
