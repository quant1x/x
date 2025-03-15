package concurrent

type coreOnce interface {
	Do(f func())
	doSlow(f func())
}
