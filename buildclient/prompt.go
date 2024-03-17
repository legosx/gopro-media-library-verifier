package buildclient

import "github.com/manifoldco/promptui"

func PromptInput(label string, mask rune, hideEntered bool, validate func(value string) error) (value string, err error) {
	prompt := promptui.Prompt{
		Label:       label,
		Mask:        mask,
		HideEntered: hideEntered,
		Validate:    validate,
	}

	return prompt.Run()
}

func PromptSelect(label string, items interface{}) (index int, err error) {
	prompt := promptui.Select{Label: label, Items: items}

	index, _, err = prompt.Run()

	return index, err
}
