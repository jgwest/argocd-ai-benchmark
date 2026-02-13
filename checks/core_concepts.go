package checks

import "argocd-ai-benchmark/types"

var _ = types.DefineChecks("tests from https://argo-cd.readthedocs.io/en/stable/core_concepts/",
	[]string{"simple"}, func() {

		types.Define().
			Prompt(
				`What Argo CD concept does this describe: When this is triggered, the latest resource state in Git is compared with the live state. It determines what is different (but don't necessarily act on it just yet)
      			Provide ONLY the answer.
        		The answer is 1 word.`).Execute().
			ExactAnswers("Refresh").
			Evaluate()

		types.Define().
			Prompt(
				`What Argo CD concept does this describe: The process of making an application move to its target state. E.g. by applying changes to a Kubernetes cluster.
      			Provide ONLY the answer.
         		The answer is 1 word.`).Execute().
			ExactAnswers("Sync").
			Evaluate()

		types.Define().
			Prompt(
				`"Unknown", "Synced", and "OutOfSync" are an example of what Argo CD concept?
				Provide ONLY the answer.
				The answer is 2 words.`).Execute().
			ExactAnswers("Sync status").
			Evaluate()

		types.Define().
			Prompt(
				`"Unknown", "Progressing", "Healthy", "Suspended", "Degraded", and "Missing", are an example of what Argo CD concept?
        		Provide ONLY the answer.
          		The answer is 2 words.`).Execute().
			ExactAnswers("Application health", "Health status").
			Evaluate()

	})
