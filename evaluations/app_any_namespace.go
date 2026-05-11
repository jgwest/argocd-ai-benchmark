package evaluations

import (
	"argocd-ai-benchmark/types"
)

var _ = types.DefinePreInitial("tests from 'https://argo-cd.readthedocs.io/en/stable/operator-manual/app-any-namespace/'",
	types.Labels("simple")).
	ResourceURLs("https://raw.githubusercontent.com/argoproj/argo-cd/refs/heads/master/docs/operator-manual/app-any-namespace.md").Start(func() {

	types.Define("T/F cluster-scoped only", types.Labels()).
		TrueOrFalse(`
					True or false: The 'applications in any namespace' feature of Argo CD can be used in either a cluster-scoped Argo CD instance, OR a namespace-scoped Argo CD instance.`).Execute().
		ExactAnswers("f").
		Evaluate()

	types.Define("T/F cluster-scoped only (reverse)", types.Labels()).
		TrueOrFalse(`
					True or false: The 'applications in any namespace' feature of Argo CD can only be used in a cluster-scoped Argo CD instance.`).Execute().
		ExactAnswers("t").
		Evaluate()

	types.Define("T/F feature is not enabled by default", types.Labels()).
		TrueOrFalse(`
					True or false: The 'applications in any namespace' feature of Argo CD must be explicitly enabled and configured appropriately (that is, it is not enabled by default).`).Execute().
		ExactAnswers("t").
		Evaluate()

	types.Define("T/F feature is not enabled by default", types.Labels()).
		TrueOrFalse(`
					True or false: The 'applications in any namespace' feature of Argo CD does not need to be explicitly enabled and configured, because it is enabled by default.`).Execute().
		ExactAnswers("f").
		Evaluate()

	types.Define("T/F argocd server needs to be configured for apps in any namespace", types.Labels()).
		TrueOrFalse(`
				True or false: In order to enable the 'applications in any namespace' feature of Argo CD, you must enable it either via a container parameter, or via an environment variable, in the Argo CD API server component.`).Execute().
		ExactAnswers("t").
		Evaluate()

	types.Define("T/F repo server does NOT need to be configured for apps in any namespace", types.Labels()).
		TrueOrFalse(`
				True or false: In order to enable the 'applications in any namespace' feature of Argo CD, you must enable it either via a container parameter, or via an environment variable, in the Argo CD repository server component.`).Execute().
		ExactAnswers("f").
		Evaluate()

	types.Define("T/F application-controller correctly configured via arg", types.Labels()).
		TrueOrFalse(`
				The following is a YAML manifest for the Argo CD Application Controller Component:
				` + "```" + `
				apiVersion: apps/v1
				kind: StatefulSet
				metadata:
				  name: argocd-application-controller
				# (...)
				spec:
				  # (...)
				  template:
				    metadata:
				      labels:
				        app.kubernetes.io/name: argocd-application-controller
				    spec:
				      containers:
				        - args:
				          - /usr/local/bin/argocd-application-controller
				          - --application-namespaces=one,two,three
				        env:
				          # (...)
				      image: quay.io/argoproj/argocd:latest
				      imagePullPolicy: Always
				      name: argocd-application-controller
				      # (...)
				` + "```" + `
				True or false: This StatefulSet has correctly enabled the 'applications in any namespace feature' for namespaces one, two, and three.`).Execute().
		ExactAnswers("t").
		Evaluate()

	types.Define("T/F application-controller configured with incorrect parameter should not be reported as correct configuration", types.Labels()).
		TrueOrFalse(`
				The following is a YAML manifest for the Argo CD Application Controller Component:
				` + "```" + `
				apiVersion: apps/v1
				kind: StatefulSet
				metadata:
				  name: argocd-application-controller
				# (...)
				spec:
				  # (...)
				  template:
				    metadata:
				      labels:
				        app.kubernetes.io/name: argocd-application-controller
				    spec:
				      containers:
				        - args:
				          - /usr/local/bin/argocd-application-controller
				          - --app-namespaces=one,two,three
				        env:
				          # (...)
				      image: quay.io/argoproj/argocd:latest
				      imagePullPolicy: Always
				      name: argocd-application-controller
				      # (...)
				` + "```" + `
				True or false: This StatefulSet has correctly enabled the 'applications in any namespace feature' for namespaces one, two, and three.`).Execute().
		ExactAnswers("f").
		Evaluate()

	types.Define("should be able to name the container parameter that is added to enable apps in any namespace", types.Labels()).
		Prompt(`
				When enabling the Argo CD feature 'applications in any namespace', what is the name of the parameter that should be added to Argo CD containers to enable this feature?

				Provide ONLY the name of the parameter.
			`).Execute().
		ExactAnswers("--application-namespaces", "application-namespaces").
		Evaluate()

	types.Define("provide the exact parameter to add to container, to enable feature for two namespaces", types.Labels()).
		Prompt(`
				I would like to enable the Argo CD feature 'applications in any namespace'. I would like to allow Argo CD to reconcile Argo CD Applications in namespaces 'ns-one' and 'ns-two'. What are the exact parameters that I should add to Argo CD Application Controller container, to enable it to use this feature with these namespaces.

				Provide ONLY the exact parameter value.
			`).Execute().
		ExactAnswers("--application-namespaces=ns-one,ns-two"). // Gemini 2.5 Pro wrapped its answer in markdown; if this happens again, consider updating the prompt.
		Evaluate()

	types.Define("provide the name of AppProject field that needs to be configured to enabled the feature", types.Labels()).
		Prompt(`

			I would like to enable the Argo CD feature 'applications in any namespace'. For this, I must modify Argo CD's AppProject resource.
			What is the name of the field that I must modify in AppProject CR to enable this feature?

			Provide ONLY the name of the field.
			`).Execute().
		ExactAnswers("sourceNamespaces", ".spec.sourceNamespaces", "spec.sourceNamespaces").
		Evaluate()

	types.Define("wildcards ARE supported via parameter", types.Labels()).
		TrueOrFalse(`
				In Argo CD, the 'applications in any namespace' feature can be configured using the '--application-namespaces' parameter.

				True or false: The '--application-namespaces' parameter does not support wildcards. Only the specific namespaces listed will have their Applications reconciled.`).Execute().
		ExactAnswers("f").Evaluate()

	types.Define("provide name of argocd-cm-params-cm field to enable the feature", types.Labels()).
		Prompt(`

			I would like to enable the Argo CD feature 'applications in any namespace'. For this, I can modify Argo CD's ConfigMap argocd-cmd-params-cm. What is the name of the field I need to add to the ConfigMap to enable this feature?

			Provide ONLY the name of the field.
			`).Execute().
		ExactAnswers("application.namespaces", ".data.application.namespaces", "data.application.namespaces").
		Evaluate()

	types.Define("never grant access to the argocd namespace within the AppProject.", types.Labels()).
		TrueOrFalse(`
				True or false: It is perfectly fine to grant access to the argocd namespace within an Argo CD AppProject, as long as Argo CD is configured to allow this.`).Execute().
		ExactAnswers("f").
		Evaluate()

	types.Define("A) applications in any namespace will be reconciled if both are true: must be listed as application namespaces via CLI param, AND in the .spec.sourceNamespaces in AppProject", types.Labels()).TrueOrFalse(`
		True or false: In order for an application to be managed and reconciled outside the Argo CD's control plane namespace, the only thing that is required is to add the namespace '--application-namespaces' parameter to application controller and server workloads.

		Once this parameter is added, Applications will be managed and reconciled within the named namespaces.`).Execute().
		ExactAnswers("f").
		Evaluate()

	types.Define("B) applications in any namespace will be reconciled if both are true: applications must be listed as application namespaces via CLI param, AND in the .spec.sourceNamespaces in AppProject", types.Labels()).TrueOrFalse(`
		True or false: In order for an application to be managed and reconciled outside the Argo CD's control plane namespace, the only thing that is required is for the AppProject named in the Application's '.spec.project' field to include that namespace in 'spec.sourceNamespaces'.

		As long as the AppProject includes the namespace within 'spec.sourceNamespaces', the Application will be reconciled`).Execute().ExactAnswers("f").Evaluate()

})
