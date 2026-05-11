package evaluations

import "argocd-ai-benchmark/types"

var _ = types.FDefinePreInitial("tests from 'https://argo-cd.readthedocs.io/en/stable/operator-manual/high_availability/'",
	types.Labels("simple")).
	ResourceURLs("https://raw.githubusercontent.com/argoproj/argo-cd/refs/heads/master/docs/operator-manual/high_availability.md").Start(func() {

	// argocd-repo-server -----------------------------------------------------

	types.Define("parallelismlimit flag", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of repository server process parameter that can be used to configure how many manifests generations can run concurrently?

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("--parallelismlimit", "parallelismlimit").
		Evaluate()

	types.Define("ARGOCD_GIT_ATTEMPTS_COUNT env var", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of repository server environment variable that can be set on Argo CD repository server process, in order to retry 'git ls-remote' calls when they fail.

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("ARGOCD_GIT_ATTEMPTS_COUNT").
		Evaluate()

	types.Define("repo-cache-expiration duration flag", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the repository server process parameter that controls how long the respository server maintains its cache of generated manifests?

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("repo-cache-expiration", "--repo-cache-expiration").
		Evaluate()

	types.Define("ARGOCD_EXEC_TIMEOUT env var", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the repository server environment variable that controls how long Argo CD waits while it executes config management tools such as helm or kustomize?

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("ARGOCD_EXEC_TIMEOUT").
		Evaluate()

	// argocd-application-controller ------------------------------------------

	types.Define("status-processors", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the application controller process parameter that controls the number of status processors for each queue?

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("--status-processors", "status-processors").
		Evaluate()

	types.Define("default # of status-processors", types.Labels()).
		Prompt(`
				In Argo CD, in Argo CD application controller, what is default number of status processors for each queue?

				Provide ONLY the answer. The answer is 1 number.`).Execute().
		ExactAnswers("20").
		Evaluate()

	types.Define("operator-processors", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the application controller process parameter that controls the number of operator processors for each queue?

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("--operation-processors", "operation-processors").
		Evaluate()

	types.Define("default # of operator-processors", types.Labels()).
		Prompt(`
				In Argo CD, in Argo CD application controller, what is default number of operator processors for each queue?

				Provide ONLY the answer. The answer is 1 number.`).Execute().
		ExactAnswers("10").
		Evaluate()

	types.Define("--repo-server-timeout-seconds", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of application controller process parameter that controls the amount of time that application controller will wait for repo server to perform manifest generation, before a context deadline exceeded is returned.

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("--repo-server-timeout-seconds", "repo-server-timeout-seconds").
		Evaluate()

	types.Define("default git poll time", types.Labels()).
		Prompt(`
				In Argo CD, what is the default interval that Argo CD will wait between polls of Git?

				Provide ONLY the answer. The answer is 1 word, expressed as a duration string (e.g. 60s, 1m, 1h, etc).`).Execute().
		ExactAnswers("3m").
		Evaluate()

	types.Define("timeout.reconciliation", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the setting that can be specified in argocd-cm ConfigMap which controls how often Argo CD will poll git?

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("timeout.reconciliation").
		Evaluate()

	types.Define("ARGOCD_CONTROLLER_REPLICAS env var", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the environment variable which can be specified on application controller which enables sharding and tells the application controller how many shards it should use?

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("ARGOCD_CONTROLLER_REPLICAS").
		Evaluate()

	types.Define("ARGO_CD_UPDATE_CLUSTER_INFO_TIMEOUT env var", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the environment variable which can be set on application controller process to control the interval between cluster information updates?

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("ARGO_CD_UPDATE_CLUSTER_INFO_TIMEOUT").
		Evaluate()

	types.Define("--sharding-method parameter", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the application controller process parameter that can be set to control the sharding method?

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("--sharding-method", "sharding-method").
		Evaluate()

	types.Define("ARGOCD_CONTROLLER_SHARDING_ALGORITHM env var", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the environment variable that can be set on the application controller process parameter to control the sharding method?

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("ARGOCD_CONTROLLER_SHARDING_ALGORITHM").
		Evaluate()

	shardingMethods := map[string]bool{
		"legacy":             true,
		"round-robin":        true,
		"consistent-hashing": true,
		// we make up some fake but plausible sharding methods to penalize guessing
		"fixed-hashing":       false,
		"fixed":               false,
		"constant-hashing":    false,
		"constant":            false,
		"scalable-hashing":    false,
		"scalable":            false,
		"lightweight-hashing": false,
		"lightweight":         false,
	}
	for k, v := range shardingMethods {

		expectedValue := "T"
		if !v {
			expectedValue = "F"
		}

		types.Define(k+" sharding method", types.Labels()).TrueOrFalse(`

		I am using Argo CD to deploy resources to my kubernetes cluster.

		True or false: '` + k + `' is the name of a supported sharding method in Argo CD, when configuring Argo CD Application Controller via parameter or environment variables`).Execute().ExactAnswers(expectedValue).Evaluate()

	}

	// Rate Limiting Application Reconciliations ------------------------------

	types.Define("WORKQUEUE_BUCKET_SIZE env var", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the environment variable which enable a bucket-based rate limiter that prevents a large number of apps from being queued at the same time, and controls the size of that bucket?

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("WORKQUEUE_BUCKET_SIZE").
		Evaluate()

	types.Define("WORKQUEUE_BUCKET_SIZE env var", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the environment variable which enable a bucket-based rate limiter that prevents a large number of apps from being queued at the same time, and controls the number of items that can be queried per second?

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("WORKQUEUE_BUCKET_QPS").
		Evaluate()

	types.Define("ARGOCD_K8SCLIENT_RETRY_MAX env var", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the environment variable which can be used to set the maximum number of retries for k8s requests? 

				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("ARGOCD_K8SCLIENT_RETRY_MAX").
		Evaluate()

})
