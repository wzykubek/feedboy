package parser

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Scheme struct {
	Url         string `yaml:"url"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Selectors   struct {
		Item            string `yaml:"item"`
		ItemTitle       string `yaml:"item_title"`
		ItemDescription string `yaml:"item_description"`
		ItemContent     string `yaml:"item_content"`
		ItemUrl         string `yaml:"item_url"`
		ItemDate        string `yaml:"item_date"`
	} `yaml:"selectors"`
	DateFormat string `yaml:"date_format"`
}

func NewScheme(filename string) (*Scheme, error) {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var scheme *Scheme
	err = yaml.Unmarshal([]byte(yamlFile), &scheme)
	if err != nil {
		return nil, err
	}
	return scheme, nil
}
