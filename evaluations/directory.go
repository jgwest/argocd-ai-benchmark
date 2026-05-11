package evaluations

import "argocd-ai-benchmark/types"

var _ = types.DefinePreInitial("tests from 'https://argo-cd.readthedocs.io/en/stable/user-guide/directory/'",
	types.Labels("simple")).
	ResourceURLs("https://raw.githubusercontent.com/argoproj/argo-cd/refs/heads/master/docs/user-guide/directory.md").Start(func() {

	types.Define("directory-type Argo CD app will read yaml and JSON only", types.Labels()).MultipleChoice(`
			A directory-type Argo CD Application will read plain manifest files contained within the git repository. Which file types will be interpreted as Kubernetes manifests?
			A) .yml, .yaml
			B) .yml, .yaml, .json
			C) .yml, .yaml, .json, kustomize
			D) .yml, .yaml, .json, kustomize, .jsonnet
			E) None of the above
			`).Execute().ExactAnswers("b").Evaluate()

	types.Define("T/F don't need to add directory", types.Labels()).TrueOrFalse(`
			When defining an Argo CD Application, and deploying plain manifests from a directory defined within a Git repository (a directory-type Argo CD Application):

			True or false: It's unnecessary to explicitly add the 'spec.source.directory' field to Argo CD Application resource, except to add additional configuration options. Argo CD will automatically detect that the source repository/path contains plain manifest files.
			`).Execute().ExactAnswers("t").Evaluate()

	types.Define("T/F don't need to add directory (reverse)", types.Labels()).TrueOrFalse(`
			When defining an Argo CD Application, and deploying plain manifests from a directory defined within a Git repository (a directory-type Argo CD Application):

			True or false: It's necessary to explicitly add the 'spec.source.directory' field to Argo CD Application resource, whether or not additional configuration options are specified. If this field is not specified, Argo CD will not automatically detect that the source repository/path contains plain manifest files.
			`).Execute().ExactAnswers("f").Evaluate()

	types.Define("enable recursive directory via CLI", types.Labels()).Prompt(`
			We are defining an Argo CD Application, and deploying plain manifests from a directory defined within a Git repository (a directory-type Argo CD Application).

			However, at present, in our Argo CD Application, Argo CD is NOT recursively reading resources across all directories in the Git repository.

			What 'argocd cli' command line parameter can be added to 'argocd app set (appname)', in order to enable this recursive reading behaviour?

			Specify only the name of the parameter. Do not specify any other text.
		`).Execute().ExactAnswers("--directory-recurse", "directory-recurse").Evaluate()

	types.Define("enable recursive directory via spec field", types.Labels()).Prompt(`

			We are defining an Argo CD Application, and deploying plain manifests from a directory defined within a Git repository (a directory-type Argo CD Application).

			However, at present, in our Argo CD Application, Argo CD is NOT recursively reading resources across all directories in the Git repository.

			What field can we enable in the Argo CD Application resource, in order to enable recursive reading of resources?

			Specify only the name/value of the field. For example: "spec.source.directory.(key): (value)" (without quotes). Include BOTH the key and the value.

		`).Execute().ExactAnswers("spec.source.directory.recurse: true").Evaluate()

	types.Define("only including certain YAML files via CLI", types.Labels()).Prompt(`

			We are defining an Argo CD Application, and deploying plain manifests from a directory defined within a Git repository (a directory-type Argo CD Application).

			However, at present, in our Argo CD Application, Argo CD is deploying files with types ".yaml", ".yml", and ".json". I only want Argo CD to deploy "*.yaml" files defined within that directory.

			What 'argocd cli' command line parameters can be added to 'argocd app set (appname)', in order to tell Argo CD to only deploy YAML files?

			Specify only the name of the parameter, and the value to that parameter. For example: ("--name value", without quotes).

		`).Execute().ExactAnswers("--directory-include \"*.yaml\"", "--directory-include *.yaml").Evaluate()

	types.Define("only including certain YAML files via spec", types.Labels()).Prompt(`

			We are defining an Argo CD Application, and deploying plain manifests from a directory defined within a Git repository (a directory-type Argo CD Application).

			However, at present, in our Argo CD Application, Argo CD is deploying files with types ".yaml", ".yml", and ".json". I only want Argo CD to deploy "*.yaml" files defined within that directory.

			What field can we set in the Argo CD Application resource, in order to tell Argo CD to only deploy "*.yaml" files?

			Only specify the name of the field and the value of the field. For example: "spec.source.directory.(key): (value)". Include BOTH the key and the value.

		`).Execute().ExactAnswers(
		"spec.source.directory.include: '*.yaml'",
		"spec.source.directory.include: \"*.yaml\"",
		"spec.source.directory.include: *.yaml").Evaluate()

	types.Define("exclude certain files via CLI", types.Labels()).Prompt(`
		We are defining an Argo CD Application, and deploying plain manifests from a directory defined within a Git repository (a directory-type Argo CD Application).

		However, at present, in our Argo CD Application, Argo CD is deploying files with types ".yaml", ".yml", and ".json".

		I want to exclude Argo CD from deploying "*.json" files defined within that directory, while still continuing to deploy "*.yaml" and "*.yml".

		What 'argocd cli' command line parameters can be added to 'argocd app set (appname)', in order to tell Argo CD to avoid deploying "*.json" files?

		Provide ONLY the answer: specify ONLY the name/value of the parameter. For example: "--name value", without quotes

		`).Execute().ExactAnswers(
		"--directory-exclude \"*.json\"",
		"--directory-exclude '*.json'",
		"--directory-exclude *.json",
		"--directory-exclude \"{*.json}\"",
	).Evaluate()

	types.Define("exclude certain files via spec field", types.Labels()).Prompt(`
		We are defining an Argo CD Application, and deploying plain manifests from a directory defined within a Git repository (a directory-type Argo CD Application).

		However, at present, in our Argo CD Application, Argo CD is deploying files with types ".yaml", ".yml", and ".json".

		I want to exclude Argo CD from deploying "*.json" files defined within that directory, while still continuing to deploy "*.yaml" and "*.yml".

		What field can we set in the Argo CD Application resource, in order to tell Argo CD to avoid deploying "*.json" files?

		Specify only the name/value of the field. For example: "spec.source.directory.(key): (value)" (without quotes). Include BOTH the key and the value.

		`).Execute().ExactAnswers(
		"spec.source.directory.exclude: \"*.json\"",
		"spec.source.directory.exclude: '*.json'",
		"spec.source.directory.exclude: *.json",
		// NOT VALID GLOB: "spec.source.directory.exclude: '{*.json}'",
		// NOT VALID GLOB: "spec.source.directory.exclude: \"{*.json}\"",
	).Evaluate()

	types.Define("files can be marked with '+argocd:skip-file-rendering'", types.Labels()).Prompt(`
		We are defining an Argo CD Application, and deploying plain manifests from a directory defined within a Git repository (a directory-type Argo CD Application).

		Within the Git repository (that Argo CD is deploying) are files that resemble Kubernetes manifests, but are not intended to be deployed as Kubernetes resources. (For example, a Helm 'values.yaml' file.)

		What comment can I add to these files to tell Argo CD to avoid attempting to deploy them?

		Answer with only the name of the comment that can be applied to the file. Don't quote or use markdown.

		`).Execute().ExactAnswers("+argocd:skip-file-rendering", "# +argocd:skip-file-rendering").Evaluate()

})
