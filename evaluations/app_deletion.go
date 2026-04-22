package evaluations

import "argocd-ai-benchmark/types"

var _ = types.DefinePreInitial("tests from 'https://argo-cd.readthedocs.io/en/stable/user-guide/app_deletion/'",
	types.Labels("simple")).
	ResourceURLs("https://raw.githubusercontent.com/argoproj/argo-cd/refs/heads/master/docs/user-guide/app_deletion.md").Start(func() {

	types.Define("kubectl delete with no finalizer", types.Labels()).
		MultipleChoice(`
				The following is an Argo CD Application:
				` + "```" + `
				apiVersion: argoproj.io/v1alpha1
				kind: Application
				metadata:
				  name: jgw-app
				  namespace: jgw
				spec:
				  destination:
				    namespace: jgw
				    server: https://kubernetes.default.svc
				  project: default
				  source:
				    path: kustomize-guestbook
				    repoURL: https://github.com/argoproj/argocd-example-apps
				    targetRevision: master
				  syncPolicy:
				    automated:
				      prune: true
				      selfHeal: true
				` + "```" + `

				The Application is Synced and Healthy, and it is being used to deploy a number of resources to the cluster.

				A user issues the following command:
				` + "```" + `
				kubectl delete -n jgw application/jgw-app
				` + "```" + `

				What happens to the resources that were deployed via the Argo CD Application:
				A) The Application and all child resources of the Application are deleted.
				B) The Application is deleted, but the child resources of the Application are not deleted.
				C) The Application is deleted, but only the immediate child resources of the Application are deleted.
				D) The Application is not deleted, because the child resources are not deleted.

          `).Execute().
		ExactAnswers("B").
		Evaluate()

	types.Define("kubectl delete with finalizer", types.Labels()).
		MultipleChoice(`
				The following is an Argo CD Application:
				` + "```" + `
				apiVersion: argoproj.io/v1alpha1
				kind: Application
				metadata:
				  name: jgw-app
				  namespace: jgw
				  finalizers:
				    - resources-finalizer.argocd.argoproj.io
				spec:
				  destination:
				    namespace: jgw
				    server: https://kubernetes.default.svc
				  project: default
				  source:
				    path: kustomize-guestbook
				    repoURL: https://github.com/argoproj/argocd-example-apps
				    targetRevision: master
				  syncPolicy:
				    automated:
				      prune: true
				      selfHeal: true
				` + "```" + `

				The Application is Synced and Healthy, and it is being used to deploy a number of resources to the cluster.

				A user issues the following command:
				` + "```" + `
				kubectl delete -n jgw application/jgw-app
				` + "```" + `

				What happens to the resources that were deployed via the Argo CD Application:
				A) The Application and all child resources of the Application are deleted.
				B) The Application is deleted, but the child resources of the Application are not deleted.
				C) The Application is deleted, but only the immediate child resources of the Application are deleted.
				D) The Application is not deleted, because the child resources are not deleted.

          `).Execute().
		ExactAnswers("A").
		Evaluate()

	types.Define("argocd app delete cascade=true with finalizer", types.Labels()).
		MultipleChoice(`
				The following is an Argo CD Application:
				` + "```" + `
				apiVersion: argoproj.io/v1alpha1
				kind: Application
				metadata:
				  name: jgw-app
				  namespace: jgw
				spec:
				  destination:
				    namespace: jgw
				    server: https://kubernetes.default.svc
				  project: default
				  source:
				    path: kustomize-guestbook
				    repoURL: https://github.com/argoproj/argocd-example-apps
				    targetRevision: master
				  syncPolicy:
				    automated:
				      prune: true
				      selfHeal: true
				` + "```" + `

				The Application is Synced and Healthy, and it is being used to deploy a number of resources to the cluster.

				A user issues the following command:
				` + "```" + `
				argocd app delete jgw-app --cascade
				` + "```" + `

				What happens to the resources that were deployed via the Argo CD Application:
				A) The Application and all child resources of the Application are deleted.
				B) The Application is deleted, but the child resources of the Application are not deleted.
				C) The Application is deleted, but only the immediate child resources of the Application are deleted.
				D) The Application is not deleted, because the child resources are not deleted.
          `).Execute().
		ExactAnswers("A").
		Evaluate()

	types.Define("argocd app delete with cascade=false", types.Labels()).
		MultipleChoice(`
				The following is an Argo CD Application:
				` + "```" + `
				apiVersion: argoproj.io/v1alpha1
				kind: Application
				metadata:
				  name: jgw-app
				  namespace: jgw
				spec:
				  destination:
				    namespace: jgw
				    server: https://kubernetes.default.svc
				  project: default
				  source:
				    path: kustomize-guestbook
				    repoURL: https://github.com/argoproj/argocd-example-apps
				    targetRevision: master
				  syncPolicy:
				    automated:
				      prune: true
				      selfHeal: true
				` + "```" + `

				The Application is Synced and Healthy, and it is being used to deploy a number of resources to the cluster.

				A user issues the following command:
				` + "```" + `
				argocd app delete jgw-app --cascade=false
				` + "```" + `

				What happens to the resources that were deployed via the Argo CD Application:
				A) The Application and all child resources of the Application are deleted.
				B) The Application is deleted, but the child resources of the Application are not deleted.
				C) The Application is deleted, but only the immediate child resources of the Application are deleted.
				D) The Application is not deleted, because the child resources are not deleted.
          `).Execute().
		ExactAnswers("B").
		Evaluate()

	types.Define("remove finalizer via kubectl, then kubectl delete command", types.Labels()).
		MultipleChoice(`
				The following is an Argo CD Application:
				` + "```" + `
				apiVersion: argoproj.io/v1alpha1
				kind: Application
				metadata:
				  name: jgw-app
				  namespace: jgw
				  finalizers:
				    - resources-finalizer.argocd.argoproj.io
				spec:
				  destination:
				    namespace: jgw
				    server: https://kubernetes.default.svc
				  project: default
				  source:
				    path: kustomize-guestbook
				    repoURL: https://github.com/argoproj/argocd-example-apps
				    targetRevision: master
				  syncPolicy:
				    automated:
				      prune: true
				      selfHeal: true
				` + "```" + `

				The Application is Synced and Healthy, and it is being used to deploy a number of resources to the cluster.

				A user issues the following commands:
				` + "```" + `
				kubectl patch app jgw-app -n jgw -p '{"metadata": {"finalizers": null}}' --type merge
				kubectl delete app jgw-app
				` + "```" + `

				What happens to the resources that were deployed via the Argo CD Application:
				A) The Application and all child resources of the Application are deleted.
				B) The Application is deleted, but the child resources of the Application are not deleted.
				C) The Application is deleted, but only the immediate child resources of the Application are deleted.
				D) The Application is not deleted, because the child resources are not deleted.
          `).Execute().
		ExactAnswers("B").
		Evaluate()

	types.Define("add finalizer via kubectl, then kubectl delete command", types.Labels()).
		MultipleChoice(`
				The following is an Argo CD Application:
				` + "```" + `
				apiVersion: argoproj.io/v1alpha1
				kind: Application
				metadata:
				  name: jgw-app
				  namespace: jgw
				spec:
				  destination:
				    namespace: jgw
				    server: https://kubernetes.default.svc
				  project: default
				  source:
				    path: kustomize-guestbook
				    repoURL: https://github.com/argoproj/argocd-example-apps
				    targetRevision: master
				  syncPolicy:
				    automated:
				      prune: true
				      selfHeal: true
				` + "```" + `

				The Application is Synced and Healthy, and it is being used to deploy a number of resources to the cluster.

				A user issues the following commands:
				` + "```" + `
				kubectl patch app jgw-app -n jgw -p '{"metadata": {"finalizers": ["resources-finalizer.argocd.argoproj.io"]}}' --type merge
				kubectl delete app jgw-app -n jgw
				` + "```" + `

				What happens to the resources that were deployed via the Argo CD Application:
				A) The Application and all child resources of the Application are deleted.
				B) The Application is deleted, but the child resources of the Application are not deleted.
				C) The Application is deleted, but only the immediate child resources of the Application are deleted.
				D) The Application is not deleted, because the child resources are not deleted.
          `).Execute().
		ExactAnswers("A").
		Evaluate()

})
