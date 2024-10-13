package mutex

type DistributedMutex interface {
	Release() error
}
