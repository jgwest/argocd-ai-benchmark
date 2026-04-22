package evaluations

import "argocd-ai-benchmark/types"

var _ = types.DefinePreInitial("tests from 'https://argo-cd.readthedocs.io/en/stable/operator-manual/user-management/'",
	types.Labels("simple")).
	ResourceURLs("https://raw.githubusercontent.com/argoproj/argo-cd/refs/heads/master/docs/operator-manual/user-management/index.md").Start(func() {

	// Local users --------------------------------------------------------------
	types.Define("name of admin user", types.Labels()).
		Prompt(`
				In Argo CD, in the context of user management, what is the name of the built-in user that has full access to the system?
				Provide ONLY the answer. The answer is 1 word.`).Execute().
		ExactAnswers("admin").
		Evaluate()

	types.Define("max length of local user username", types.Labels()).
		Prompt(`
				In Argo CD, in the context of user management, what is the maximum length of a local account's username?
				Provide ONLY the answer. The answer is a single number.`).Execute().
		ExactAnswers("32").
		Evaluate()

	types.Define("detect valid local user configuration", types.Labels()).
		TrueOrFalse(`
				The following is a simplified Argo CD ConfigMap 'argocd-cm', which can be used to configure local user accounts in Argo CD.

				True or false: This is a valid example of how to configure Argo CD 'argocd-cm' ConfigMap for enabling a local user account named 'dave':
				` + "```yaml" + `
				apiVersion: v1
				kind: ConfigMap
				metadata:
				  name: argocd-cm
				  namespace: argocd
				data:
				  accounts.dave: apiKey, login

				  # (...other config map data...)
				` + "```" + `
				`).Execute().
		ExactAnswers("t").Evaluate()

	types.Define("detect typo in local user configuration", types.Labels()).
		TrueOrFalse(`
				The following is a simplified Argo CD ConfigMap 'argocd-cm', which can be used to configure local user accounts in Argo CD.

				True or false: This is a valid example of how to configure Argo CD 'argocd-cm' ConfigMap for enabling a local user account named 'dave':
				` + "```yaml" + `
				apiVersion: v1
				kind: ConfigMap
				metadata:
				  name: argocd-cm
				  namespace: argocd
				data:
				  account.dave: apiKey, login

				  # (...other config map data...)
				` + "```" + `
				`).Execute().
		ExactAnswers("f").Evaluate() // False because the correct value is 'accounts.dave', not 'account.dave'

	// Dex ----------------------------------------------------------------------

	types.Define("name of ConfigMap that should be used to enable dex sso", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the ConfigMap that should be modified in order to enable SSO via dex?
      			Provide ONLY the answer.`).Execute().
		ExactAnswers("argocd-cm").
		Evaluate()

	types.Define("name of the field in argocd-cm Configmap that should be modified for dex sso", types.Labels()).
		Prompt(`
				In Argo CD, the 'argocd-cm' ConfigMap can be modified in order to enable SSO support via dex. What is the name of the field that should be modified in this ConfigMap to enable dex?

				Provide ONLY the answer. The answer is the name of the field.`).Execute().
		ExactAnswers("dex.config", "data.dex.config", ".data.dex.config").
		Evaluate()

	types.Define("detect valid dex configuration for GitHub", types.Labels()).
		TrueOrFalse(`
				In Argo CD, the 'argocd-cm' ConfigMap can be modified in order to enable SSO support via dex. In this case, we are using dex to authenticate via GitHub.

				True or False: this example configuration is a correct and valid example of how to configure Argo CD to connect to GitHub using Dex:
				` + "```yaml" + `
				apiVersion: v1
				kind: ConfigMap
				metadata:
				  name: argocd-cm
				  namespace: argocd
				data:
				  url: https://argocd.example.com
				  dex.config: |
				    connectors:
				      # GitHub example
				      - type: github
				        id: github
				        name: GitHub
				        config:
				          clientID: aabbccddeeff00112233
				          clientSecret: $dex.github.clientSecret
				          orgs:
				          - name: your-github-org
				  # (...)
				` + "```" + `

				If the field names are not correct, or the data is not in the correct format, the answer should be false.
				If the field names are correct, and the data is in the correct format, the answer should be true.
          `).Execute().
		ExactAnswers("t").
		Evaluate()

	types.Define("detect invalid dex configuration due to incorrect field names", types.Labels()).
		TrueOrFalse(`
				In Argo CD, the 'argocd-cm' ConfigMap can be modified in order to enable SSO support via dex. In this case, we are using dex to authenticate via GitHub.

				True or False: this example configuration is a correct and valid example of how to configure Argo CD to connect to GitHub using Dex:
				` + "```yaml" + `
				apiVersion: v1
				kind: ConfigMap
				metadata:
				  name: argocd-cm
				  namespace: argocd
				data:
				  url: https://argocd.example.com
				  dex.config: |
				    connectors:
				      # GitHub example
				      - type: github
				        id: github
				        name: GitHub
				        config:
				          clientAccountID: aabbccddeeff00112233
				          clientAccountSecret: $dex.github.clientSecret
				          orgs:
				          - name: your-github-org
				  # (...)
				` + "```" + `

				If the field names are not correct, or the data is not in the correct format, the answer should be false.
				If the field names are correct, and the data is in the correct format, the answer should be true.
          `).Execute().
		ExactAnswers("f"). // clientAccountID should be clientID, clientAccountSecer should be clientSecret
		Evaluate()

	// OIDC Configuration with DEX ----------------------------------------------------------------------

	types.Define("detect valid OIDC Configuration with dex", types.Labels()).
		TrueOrFalse(`
				In Argo CD, the 'argocd-cm' ConfigMap can be modified in order to enable SSO support via dex. In this case, we are using dex to authenticate via OIDC.

				True or False: this example configuration is a correct and valid example of how to configure Argo CD to use OIDC via dex:
				` + "```yaml" + `
				apiVersion: v1
				kind: ConfigMap
				metadata:
				  name: argocd-cm
				  namespace: argocd
				data:
				  url: https://argocd.example.com
				  dex.config: |
				    connectors:
				      - type: oidc
				        id: oidc
				        name: OIDC
				        config:
				          issuer: https://example-OIDC-provider.example.com
				          clientID: aaaabbbbccccddddeee
				          clientSecret: $dex.oidc.clientSecret
				  # (...)
				` + "```" + `

				If the field names are not correct, or the data is not in the correct format, the answer should be false.
				If the field names are correct, and the data is in the correct format, the answer should be true.
          `).Execute().
		ExactAnswers("t").
		Evaluate()

	types.Define("detect invalid OIDC Configuration with dex", types.Labels()).
		TrueOrFalse(`
				In Argo CD, the 'argocd-cm' ConfigMap can be modified in order to enable SSO support via dex. In this case, we are using dex to authenticate via OIDC.

				True or False: this example configuration is a correct and valid example of how to configure Argo CD to use OIDC via dex:
				` + "```yaml" + `
				apiVersion: v1
				kind: ConfigMap
				metadata:
				  name: argocd-cm
				  namespace: argocd
				data:
				  url: https://argocd.example.com
				  dex.config: |
				    connectors:
				      - type: oidc-client
				        id: oidc-client
				        name: oidc-client
				        config:
				          issuer: https://example-provider.example.com
				          clientID: aaaabbbbccccddddeee
				          clientSecret: $dex.oidc.clientSecret
				  # (...)
				` + "```" + `

				If the field names are not correct, or the data is not in the correct format, the answer should be false.
				If the field names are correct, and the data is in the correct format, the answer should be true.
          `).Execute().
		ExactAnswers("f"). // false as 'type' field should be 'oidc', not 'oidc-client'
		Evaluate()

	// Skipping certificate verification on OIDC provider connections -----------
	types.Define("output name of 'oidc.tls.insecure.skip.verify field'", types.Labels()).
		Prompt(`
				In Argo CD, what is the name of the field that can be set in the 'argocd-cm' ConfigMap that, when set, will tell Argo CD to skip validation of TLS certificates when connecting to an OIDC provider.

				Provide ONLY the answer. The answer is the name of the field to set in the ConfigMap.
				`).
		Execute().
		ExactAnswers("oidc.tls.insecure.skip.verify", "data.oidc.tls.insecure.skip.verify").
		Evaluate()
})
