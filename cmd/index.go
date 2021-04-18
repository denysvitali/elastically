package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"sort"
)

type IndexCmd struct {
	List *IndexListCmd `arg:"subcommand:list"`
	Delete *IndexDeleteCmd `arg:"subcommand:delete"`
}

func (i IndexCmd) Run(c *elastic.Client, logger *logrus.Logger) error {
	if i.List != nil {
		return i.List.Run(c, logger)
	}
	if i.Delete != nil {
		return i.Delete.Run(c, logger)
	}

	return fmt.Errorf("please specify a subcommand")
}

type IndexListCmd struct {}
type IndexDeleteCmd struct {
	IndexName string `arg:"positional,required"`
}

var _ RunnableCmd = IndexDeleteCmd{}
func (i IndexDeleteCmd) Run(c *elastic.Client, l *logrus.Logger) error {
	idxDeleteResp, err := c.DeleteIndex(i.IndexName).Do(context.Background())
	if err != nil {
		return fmt.Errorf("unable to delete index %s: %v", i.IndexName, err)
	}

	if !idxDeleteResp.Acknowledged {
		return fmt.Errorf("deletion not acknowledged")
	}

	fmt.Printf("%s deleted successfully!\n", i.IndexName)
	return nil
}

type sortByName []string

func (s sortByName) Len() int {
	return len(s)
}

func (s sortByName) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var _ sort.Interface = sortByName{}

func (i IndexListCmd) Run(c *elastic.Client, logger *logrus.Logger) error {
	indexStats, err := c.IndexStats().Do(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get cluster stats: %v", err)
	}

	var indices []string
	for k, _ := range indexStats.Indices {
		indices = append(indices, k)
	}
	sort.Sort(sortByName(indices))

	for _, idx := range indices {
		fmt.Printf("%s\n", idx)
	}

	return nil
}

var _ RunnableCmd = IndexCmd{}
var _ RunnableCmd = IndexListCmd{}