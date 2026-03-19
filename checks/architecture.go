package checks

import (
	"argocd-ai-benchmark/types"
	"log"
	"math/rand"
)

var _ = types.DefinePreInitial("tests from 'https://argo-cd.readthedocs.io/en/stable/operator-manual/architecture/'",
	types.Labels("simple")).
	ResourceURLs("https://raw.githubusercontent.com/argoproj/argo-cd/refs/heads/master/docs/operator-manual/architecture.md").Start(func() {

	const (
		apiServer     = "API Server"
		appController = "Application Controller"
		repoServer    = "Repository Server"
	)

	type question struct {
		responsibility string
		owner          string
	}

	questions := []question{
		{
			responsibility: "This component is a gRPC/REST server which exposes the API consumed by the Web UI, CLI, and CI/CD systems.",
			owner:          apiServer,
		},
		// I've removed this one as it's too generic: it could also apply to application controller.
		// {
		// 	responsibility: "application management and status reporting",
		// 	owner:          apiServer,
		// },
		{
			responsibility: "invoking of application operations (e.g. sync, rollback, user-defined actions)",
			owner:          apiServer,
		},
		{
			responsibility: "repository and cluster credential management (stored as K8s secrets)",
			owner:          apiServer,
		},
		{
			responsibility: "authentication and auth delegation to external identity providers",
			owner:          apiServer,
		},
		{
			responsibility: "RBAC enforcement",
			owner:          apiServer,
		},
		{
			responsibility: "listener/forwarder for Git webhook events",
			owner:          apiServer,
		},
		{
			responsibility: "This component is an internal service which maintains a local cache of the Git repository holding the application manifests.",
			owner:          repoServer,
		},
		{
			responsibility: "This component is responsible for generating and returning the Kubernetes manifests when provided the following inputs: repository URL, revision, application path, etc.",
			owner:          repoServer,
		},
		{
			responsibility: "This component is a Kubernetes controller which continuously monitors running applications and compares the current, live state against the desired target state (as specified in the repo).",
			owner:          appController,
		},
		{
			responsibility: "This component detects OutOfSync application state and optionally takes corrective action.",
			owner:          appController,
		},
		{
			responsibility: "This component is responsible for invoking any user-defined hooks for lifecycle events (PreSync, Sync, PostSync).",
			owner:          appController,
		},
	}

	nameToMultipleChoice := map[string]string{
		apiServer:     "A",
		appController: "B",
		repoServer:    "C",
	}

	// Shuffle questions with a specific hardcoded seed for deterministic results
	rng := rand.New(rand.NewSource(12))
	shuffledQuestions := make([]question, len(questions))
	copy(shuffledQuestions, questions)
	rng.Shuffle(len(shuffledQuestions), func(i, j int) {
		shuffledQuestions[i], shuffledQuestions[j] = shuffledQuestions[j], shuffledQuestions[i]
	})

	for _, question := range shuffledQuestions {

		expectedAnswer, ok := nameToMultipleChoice[question.owner]
		if !ok {
			log.Fatalln("unexpected owner:" + question.owner)
		}

		types.Define(question.responsibility, types.Labels()).
			MultipleChoice(`
					Which Argo CD component does this description describe:
					"` + question.responsibility + `"

					A) API Server
					B) Application Controller
					C) Repository Server
					D) None of the above
					`).Execute().
			ExactAnswers(expectedAnswer).
			Evaluate()
	}

})
