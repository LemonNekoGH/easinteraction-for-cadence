package templates

import (
	_ "embed"
	"text/template"
)

var (
	//go:embed contract_interaction.gohtml
	templateInteractionStr string
	//go:embed query_script.gohtml
	templateQueryScriptStr string
	//go:embed tx_script.gohtml
	templateTxScriptStr string

	TemplateInteraction = template.Must(template.New("Interaction").Parse(templateInteractionStr))
	TemplateQueryScript = template.Must(template.New("QueryScript").Parse(templateQueryScriptStr))
	TemplateTxScript    = template.Must(template.New("TxScript").Parse(templateTxScriptStr))
)
