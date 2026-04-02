package checks

import "argocd-ai-benchmark/types"

var _ = types.FDefinePreInitial("tests from 'https://argo-cd.readthedocs.io/en/stable/user-guide/sync-options/'",
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

				What is the name of the Argo CD annotation that I can add to '.metadata.annotations' for a resource, to require manual confirmation before that resource is pruned?

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

	// JGW-TODO: stopped here

})
