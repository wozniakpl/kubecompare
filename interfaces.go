package main

type KubectlInterface interface {
	getRolloutHistory(resourceType, resourceName string) (string, error)
	getRolloutHistoryWithRevision(resourceType, resourceName string, revision int) (string, error)
}

type OutputWriter interface {
	Write(string)
}
