package main

type KubectlInterface interface {
	getRolloutHistory(resourceType, resourceName, namespace string) (string, error)
	getRolloutHistoryWithRevision(resourceType, resourceName string, revision int, namespace string) (string, error)
}

type OutputWriter interface {
	Write(string)
}
