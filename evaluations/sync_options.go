package evaluations

import "argocd-ai-benchmark/types"

var _ = types.DefinePreInitial("tests from 'https://argo-cd.readthedocs.io/en/stable/user-guide/sync-options/'",
	types.Labels("simple")).
	ResourceURLs("https://raw.githubusercontent.com/argoproj/argo-cd/refs/heads/master/docs/user-guide/sync-options.md").Start(func() {

	// No Prune Resources -------------------------------------------------------
	types.Define("'argocd.argoproj.io/sync-options: Prune=false' annotation can prevent pruning of resource", types.Labels()).
		Prompt(`
				I am using Argo CD to deploy resources to my kubernetes cluster. There is a kubernetes resource that I am deploying, which I don't want to be pruned under any circumstances.

				What is the name and value of the annotation that I can add to '.metadata.annotations' for a resource, to prevent it from being pruned?

				Provide ONLY the answer. The answer is the name and value of the annotation to add. Dont quote the answer or use markdown.`).Execute().
		ExactAnswers("argocd.argoproj.io/sync-options: Prune=false").
		Evaluate()

	// Resource Pruning With Confirmation ---------------------------------------
	types.Define("'argocd.argoproj.io/sync-options: Prune=confirm' annotation can be used to require confirmation of pruning", types.Labels()).
		Prompt(`
				I am using Argo CD to deploy resources to my kubernetes cluster. There is a kubernetes resource that I am deploying, which I don't want to be pruned except when I manually confirm it should be pruned.

				What is the name of the Argo CD annotation that I can add to '.metadata.annotations' for a resource, to require manual confirmation before that resource is pruned? Provide the annotation specific to prune behaviour.

				Provide ONLY the answer. The answer is the name and value of the annotation to add. Dont quote the answer or use markdown.`).Execute().
		ExactAnswers("argocd.argoproj.io/sync-options: Prune=confirm").
		Evaluate()

	// Disable Kubectl Validation -----------------------------------------------
	types.Define("'argocd.argoproj.io/sync-options: Validate=false' annotation can be used to require kubectl apply with --validate=false", types.Labels()).
		Prompt(`
				I am using Argo CD to deploy resources to my kubernetes cluster. There is a kubernetes resource that I am deploying, which requires that is be applied using 'kubectl apply' with the '--validate=false' flag.

				What is the name of the Argo CD annotation that I can add to '.metadata.annotations' for a resource, to require that that resource is applied with the '--validate=false' flag of 'kubectl apply'?

				Provide ONLY the answer. The answer is the name and value of the annotation to add. Dont quote the answer or use markdown.`).Execute().
		ExactAnswers("argocd.argoproj.io/sync-options: Validate=false").
		Evaluate()

	// Skip Dry Run for new custom resources types ------------------------------

	types.Define("'argocd.argoproj.io/sync-options: SkipDryRunOnMissingResource=true' annotation can be used to skip missing resources on dry run", types.Labels()).Prompt(`
			I am using Argo CD to deploy resources to my kubernetes cluster. One of the resources I am deploying is a custom kubernetes resource which is defined via CustomResourceDefinition.

			However, Argo CD has not yet deployed the CustomResourceDefinition for the resource. This means Argo CD returns an error when I attempt to deploy my custom resource.

			What is the name of the Argo CD annotation that I can add to '.metadata.annotations' for a resource, to skip dry run for that resource so that Argo CD will still attempt to deploy that resource?

			Provide ONLY the answer. The answer is the name and value of the annotation to add. Dont quote the answer or use markdown.`).Execute().
		ExactAnswers("argocd.argoproj.io/sync-options: SkipDryRunOnMissingResource=true").Evaluate()

	types.Define("Application spec.syncPolicy.syncOptions[SkipDryRunOnMissingResource=true] can be used to skip all missing resources for an Application", types.Labels()).Prompt(`
			I am using Argo CD to deploy resources to my kubernetes cluster. One of the resources I am deploying is a custom kubernetes resource which is defined via CustomResourceDefinition (CRD).

			However, Argo CD has not yet deployed the CustomResourceDefinition for the resource. This means Argo CD returns an error when I attempt to deploy my custom resource.

			What is the name of the Argo CD sync option that I can add to 'spec.syncPolicy.syncOptions' of an Argo CD Application, to skip dry run for all resources in that Application, so that Argo CD will still attempt to deploy any resources which are missing the CRD?

			Provide ONLY the answer. The answer is the name and value of the sync option to add. Dont quote the answer or use markdown.`).Execute().
		ExactAnswers("SkipDryRunOnMissingResource=true").Evaluate()

	// No Resource Deletion -----------------------------------------------------

	types.Define("'argocd.argoproj.io/sync-options: Delete=false' annotation can be added to prevent deletion", types.Labels()).Prompt(`
			I am using Argo CD to deploy resources to my kubernetes cluster.

			For one of the resources I am deploying, I don't want that resource to be deleted when the parent Argo CD Application (that contains the resource) is deleted.

			What is the name of the Argo CD annotation that I can add to '.metadata.annotations' for a resource, to prevent that resource from being cleaned up when the Argo CD Application is deleted?

			Provide ONLY the answer. The answer is the name and value of the annotation to add. Dont quote the answer or use markdown.
			`).Execute().
		ExactAnswers("argocd.argoproj.io/sync-options: Delete=false").Evaluate()

	// Resource Deletion With Confirmation --------------------------------------

	types.Define("'argocd.argoproj.io/sync-options: Delete=confirm' annotation can be added to require deletion confirmation", types.Labels()).Prompt(`

			I am using Argo CD to deploy resources to my kubernetes cluster.

			For one of the resources I am deploying, I don't want that resource to be deleted without manual confirmation.

			What is the name of the Argo CD annotation that I can add to '.metadata.annotations' for a resource, to prevent that resource from being deleted without manual confirmation?

			Provide ONLY the answer. The answer is the name and value of the annotation to add. Dont quote the answer or use markdown.
			`).Execute().
		ExactAnswers("argocd.argoproj.io/sync-options: Delete=confirm").Evaluate()

	// Selective Sync -----------------------------------------------------------

	types.Define("spec.syncPolicy.syncOptions[ApplyOutOfSyncOnly=true] can be used to enable selective sync", types.Labels()).Prompt(`
			I am using Argo CD to deploy resources to my kubernetes cluster.

			I am deploying thousands of objects, which means it can take a long time for Argo CD to sync all those objects (and also puts pressure on API server).

			What is the name of the Argo CD sync option that I can add to Application's 'spec.syncPolicy.syncOptions' field, to ensure that Argo CD only attempts to sync resources that are out of sync?

			Provide ONLY the answer. The answer is the name and value of the sync policy to add. Dont quote the answer or use markdown.`).Execute().
		ExactAnswers("ApplyOutOfSyncOnly=true").Evaluate()

	// Prune Last ---------------------------------------------------------------

	types.Define("'argocd.argoproj.io/sync-options: PruneLast=true' annotation can be added to trigger prune as final, implicit wave of sync", types.Labels()).Prompt(`
			I am using Argo CD to deploy resources to my kubernetes cluster.

			What is the name of the Argo CD annotation that I can add to a resource's '.metadata.annotations' field to trigger the pruning of that resource as a final, implicit wave of a sync operation (after all other resources are synced/healthy)?

			Provide ONLY the answer. The answer is the name and value of the annotation to add. Dont quote the answer or use markdown.`).Execute().
		ExactAnswers("argocd.argoproj.io/sync-options: PruneLast=true").Evaluate()

	types.Define("spec.syncPolicy.syncOptions[PruneLast=true] on Application can be used to trigger prune as final, implicit wave of sync", types.Labels()).Prompt(`
			I am using Argo CD to deploy resources to my kubernetes cluster.

			What is the name of the Argo CD sync option that I can add to Application's 'spec.syncPolicy.syncOptions' field to trigger the pruning of resources as a final, implicit wave of a sync operation (after all other resources are synced/healthy)?

			Provide ONLY the answer. The answer is the name and value of the sync policy to add. Dont quote the answer or use markdown.`).Execute().
		ExactAnswers("PruneLast=true").Evaluate()

	// Replace Resource Instead Of Applying Changes -----------------------------

	types.Define("'argocd.argoproj.io/sync-options: Replace=true' annotation can be added to enable replace, rather than apply, on a resource", types.Labels()).Prompt(`
			I am using Argo CD to deploy resources to my kubernetes cluster.

			What is the name of the Argo CD annotation that I can add to '.metadata.annotations' of a resource, to use 'kubectl replace' rather than 'kubectl apply', when modifying that resource (when it is out of sync).

			Provide ONLY the answer. The answer is the name and value of the annotation to add. Dont quote the answer or use markdown.`).Execute().
		ExactAnswers("argocd.argoproj.io/sync-options: Replace=true").Evaluate()

	types.Define("spec.syncPolicy.syncOptions[Replace=true] can be used to enable replace, rather than apply, on an Application", types.Labels()).Prompt(`

			I am using Argo CD to deploy resources to my kubernetes cluster.

			What is the name of the Argo CD sync option that I can add to Application's 'spec.syncPolicy.syncOptions' field, to use 'kubectl replace' rather than 'kubectl apply' when modifying resources that are out of sync.

			Provide ONLY the answer. The answer is the name and value of the sync policy to add. Dont quote the answer or use markdown.`).Execute().ExactAnswers("Replace=true").Evaluate()

	// Server-Side Apply --------------------------------------------------------

	types.Define("'argocd.argoproj.io/sync-options: ServerSideApply=true' annotation can be used to enable server side apply", types.Labels()).Prompt(`
		I am using Argo CD to deploy resources to my kubernetes cluster.

		What is the name of the Argo CD annotation that I can add to an individual resource's '.metadata.annotations' field, to trigger server side apply for that resource?

		Provide ONLY the answer. The answer is the name and value of the annotation to add. Dont quote the answer or use markdown.
		`).Execute().
		ExactAnswers("argocd.argoproj.io/sync-options: ServerSideApply=true").
		Evaluate()

	types.Define("'argocd.argoproj.io/sync-options: ServerSideApply=false' annotation can be used to disable server side apply", types.Labels()).Prompt(`
		I am using Argo CD to deploy resources to my kubernetes cluster.

		What is the name of the Argo CD annotation that I can add to an individual resource's '.metadata.annotations' field, to disable server side apply for that resource (when the Application that resource is part of has it enabled?)

		Provide ONLY the answer. The answer is the name and value of the annotation to add. Dont quote the answer or use markdown.
		`).Execute().
		ExactAnswers("argocd.argoproj.io/sync-options: ServerSideApply=false").
		Evaluate()

	types.Define("spec.syncPolicy.syncOptions[ServerSideApply=true] can be used to enabled server side apply", types.Labels()).Prompt(`

			I am using Argo CD to deploy resources to my kubernetes cluster.

			What is the name of the Argo CD sync option that I can add to Application's 'spec.syncPolicy.syncOptions' field, to use server side apply, when modifying resources that are out of sync.

			Provide ONLY the answer. The answer is the name and value of the sync policy to add. Dont quote the answer or use markdown.`).Execute().ExactAnswers("ServerSideApply=true").Evaluate()

	types.Define("Replace=true takes precedence over ServerSideApply=true", types.Labels()).TrueOrFalse(`

		I am using Argo CD to deploy resources to my kubernetes cluster.

		True or false: When deploying resources via Argo CD, the ServerSideApply=true sync option takes precedence over Replace=true sync option.`).Execute().ExactAnswers("F"). //  Replace=true takes precedence over ServerSideApply=true.
		Evaluate()

	types.Define("Argo CD can use server side apply to patch existing resources on the cluster that are not fully managed by Argo CD.", types.Labels()).TrueOrFalse(`
		True or false: Argo CD does not currently support the patching of existing resources on the cluster that are not already fully managed by Argo CD.
		`).Execute().
		ExactAnswers("F"). // server side apply can be used to "[patch] existing resources on the cluster that are not fully managed by Argo CD."
		Evaluate()

	// Fail the sync if a shared resource is found -------------------------------
	//
	types.Define("FailOnSharedResource=true can be used to fail the sync if there are shared resources", types.Labels()).Prompt(`

		I am using Argo CD to deploy resources to my kubernetes cluster.

		What is the name of the Argo CD sync option that I can add to Application's 'spec.syncPolicy.syncOptions' field, to tell Argo CD to fail the sync operation if Argo CD detects that there are resources shared between multiple Argo CD Applications deploying to a cluster?

		Provide ONLY the answer. The answer is the name and value of the sync policy to add. Dont quote the answer or use markdown.`).
		Execute().
		ExactAnswers("FailOnSharedResource=true").
		Evaluate()

	// Respect ignore differences configs ----------------------------------------

	types.Define("RespectIgnoreDifferences can be enabled to respect ignore differences during sync", types.Labels()).Prompt(`

		I am using Argo CD to deploy resources to my kubernetes cluster.

		What is the name of the Argo CD sync option that I can add to Application's 'spec.syncPolicy.syncOptions' field, to tell Argo CD to consider the configurations made in the Application's '.spec.ignoreDifferences' field during the sync stage (that is, not JUST when computing the diff between the live and desired state).

		Provide ONLY the answer. The answer is the name and value of the sync policy to add. Dont quote the answer or use markdown.`).
		Execute().
		ExactAnswers("RespectIgnoreDifferences=true").
		Evaluate()

	// Create Namespace ----------------------------------------------------------

	types.Define("CreateNamespace=true will create a target namespace if needed", types.Labels()).Prompt(`

		I am using Argo CD to deploy resources to my kubernetes cluster.

		What is the name of the Argo CD sync option that I can add to Application's 'spec.syncPolicy.syncOptions' field, to tell Argo CD to create the namespace specified in the spec.destination.namespace if it doesn't exist?

		Provide ONLY the answer. The answer is the name and value of the sync policy to add. Dont quote the answer or use markdown.`).
		Execute().
		ExactAnswers("CreateNamespace=true").
		Evaluate()

	types.Define("namespace metadata can be specified via .spec.syncPolicy.managedNamespaceMetadata", types.Labels()).Prompt(`

		I am using Argo CD to deploy resources to my kubernetes cluster.

		What is the name of the Argo CD Application field that can be used to add labels/annotations to Namespaces that are created via the 'Namespace=True' sync option?

		Provide ONLY the answer. The answer is the name of the field in the Argo CD Application resource. Dont quote the answer or use markdown.`).
		Execute().
		ExactAnswers(
			"managedNamespaceMetadata",
			"spec.syncPolicy.managedNamespaceMetadata",
			".spec.syncPolicy.managedNamespaceMetadata").
		Evaluate()
})
