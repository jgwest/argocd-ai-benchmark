package checks

import "argocd-ai-benchmark/types"

var _ = types.DefinePreInitial("tests from 'https://argo-cd.readthedocs.io/en/stable/core_concepts/'",
	types.Labels("simple")).
	ResourceURLs("https://raw.githubusercontent.com/argoproj/argo-cd/refs/heads/master/docs/core_concepts.md").Start(func() {

	types.Define("refresh", types.Labels()).
		Prompt(`
			What Argo CD concept does this describe: When this is triggered, the latest resource state in Git is compared with the live state. determines what is different (but don't necessarily act on it just yet)
			Provide ONLY the answer.
			The answer is 1 word.`).Execute().
		ExactAnswers("Refresh").
		Evaluate()

	types.Define("sync", types.Labels()).
		Prompt(`
			What Argo CD concept does this describe: The process of making an application move to its target state. E.g. by applying changes to a Kubernetes cluster.
			Provide ONLY the answer.
			The answer is 1 word.`).Execute().
		ExactAnswers("Sync").
		Evaluate()

	types.Define("sync status", types.Labels()).
		Prompt(`
			"Unknown", "Synced", and "OutOfSync" are an example of what Argo CD concept?
			Provide ONLY the answer.
			The answer is 2 words.`).Execute().
		ExactAnswers("Sync status").
		Evaluate()

	types.Define("health status", types.Labels()).
		Prompt(`
			"Unknown", "Progressing", "Healthy", "Suspended", "Degraded", and "Missing", are an example of what Argo CD concept?
			Provide ONLY the answer.
			The answer is 2 words.`).Execute().
		ExactAnswers("Application health", "Health status").
		Evaluate()

})
