package types

type promptType string

const (
	promptType_Generic        = ""
	promptType_TrueOrFalse    = "true-or-false"
	promptType_MultipleChoice = "multiple-choice"
)

type PreInitial struct {
}

func DefinePreInitial(name string, labelsParam []string) PreInitial {
	definePreInitialInternal(name, labelsParam)

	return PreInitial{}
}

func FDefinePreInitial(name string, labelsParam []string) PreInitial {

	definePreInitialInternal(name, labelsParam)
	globalState.Focus = true

	return PreInitial{}
}

func definePreInitialInternal(name string, labelsParam []string) {
	globalState.clearEphemeral()

	globalState.name = name
	globalState.labels = labelsParam
}

func (pi PreInitial) ResourceURLs(urls ...string) PreInitial {
	globalState.resourceURLs = urls
	return pi
}

func (pi PreInitial) Start(f func()) any {
	f()
	return nil
}

type Initial struct {
	name         string
	labels       []string
	prompt       string
	focus        bool
	resourceURLS []string
	promptType   promptType
}

func (m Initial) Name(nameParam string) Initial {
	m.name = nameParam
	return m
}

func (m Initial) MultipleChoice(promptParam string) Initial {
	m.prompt = promptParam
	m.promptType = promptType_MultipleChoice
	m.labels = mergeStringSlices(m.labels, []string{"multiple-choice"})
	return m
}

func (m Initial) TrueOrFalse(promptParam string) Initial {
	m.prompt = promptParam
	m.promptType = promptType_TrueOrFalse
	m.labels = mergeStringSlices(m.labels, []string{"true-or-false"})
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
	globalState.globalEvaluations = append(globalState.globalEvaluations, e)
}

type globalStateType struct {
	// These fields should be reset by define():
	name         string
	labels       []string
	Focus        bool
	resourceURLs []string
	// Note: when a new field is added, ensure that the field is cleared in clearEphemeral()

	// These fields should not be reset by define:
	globalEvaluations []Evaluation
}

func (gs *globalStateType) clearEphemeral() {
	gs.name = ""
	gs.labels = []string{}
	gs.Focus = false
	gs.resourceURLs = nil
}

var globalState globalStateType = globalStateType{}

func FDefine(name string, labelsParam []string) Initial {

	res := Define(name, labelsParam)

	res.focus = true

	return res

}

func Define(name string, labelsParam []string) Initial {

	res := Initial{
		name:         globalState.name + " -> " + name,
		labels:       []string{},
		focus:        globalState.Focus,
		promptType:   promptType_Generic,
		resourceURLS: globalState.resourceURLs,
	}

	// evaluations can also add labels, in addition to those defined at file scope
	res.labels = append(res.labels, globalState.labels...)
	res.labels = mergeStringSlices(res.labels, labelsParam)

	return res

}

func Labels(labels ...string) []string {
	return labels
}

func Evaluations() []Evaluation {
	return globalState.globalEvaluations
}

// mergeStringSlices takes two string slices and merges them,
// returning a new slice with all unique strings.
func mergeStringSlices(slice1, slice2 []string) []string {
	allKeys := make(map[string]struct{})

	// Add all elements from the first slice to the map.
	for _, item := range slice1 {
		allKeys[item] = struct{}{}
	}

	// Add all elements from the second slice to the map.
	// Duplicates will be ignored automatically.
	for _, item := range slice2 {
		allKeys[item] = struct{}{}
	}

	// Create a new slice to hold the result.
	// We can pre-allocate the capacity for efficiency.
	result := make([]string, 0, len(allKeys))

	// Convert the map keys back into a slice.
	for key := range allKeys {
		result = append(result, key)
	}

	return result
}
