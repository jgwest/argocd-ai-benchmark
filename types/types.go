package types

type Initial struct {
	name   string
	labels []string
	prompt string
	focus  bool
}

func (m Initial) Name(nameParam string) Initial {
	m.name = nameParam
	return m
}

func (m Initial) Prompt(promptParam string) Initial {
	m.prompt = promptParam
	return m
}

func (m Initial) Execute() Evaluation {
	return Evaluation{
		initial: m,
	}
}

type Evaluation struct {
	initial      Initial
	exactAnswers []string
}

func (e Evaluation) ExactAnswers(answers ...string) Evaluation {
	e.exactAnswers = answers
	return e
}

// Terminal
func (e Evaluation) Evaluate() {
	globalState.GlobalEvaluations = append(globalState.GlobalEvaluations, e)
}

type globalStatus struct {
	// These fields should be reset by define():
	Name   string
	Labels []string
	Focus  bool

	// These fields should not be reset by define:
	GlobalEvaluations []Evaluation
}

func (gs *globalStatus) clearEphemeral() {
	gs.Name = ""
	gs.Labels = []string{}
	gs.Focus = false
}

var globalState globalStatus = globalStatus{}

func FDefineChecks(nameParam string, labelsParam []string, f func()) any {
	globalState.clearEphemeral()
	globalState.Focus = true
	return defineChecksInner(nameParam, labelsParam, f)
}

func DefineChecks(nameParam string, labelsParam []string, f func()) any {
	globalState.clearEphemeral()
	return defineChecksInner(nameParam, labelsParam, f)
}

func defineChecksInner(nameParam string, labelsParam []string, f func()) any {
	globalState.Name = nameParam
	globalState.Labels = labelsParam

	f()
	return nil
}

func Define(labelsParam ...string) Initial {
	return Initial{
		name:   globalState.Name,
		labels: append(globalState.Labels, labelsParam...),
		focus:  globalState.Focus,
	}
}

func Evaluations() []Evaluation {
	return globalState.GlobalEvaluations
}
