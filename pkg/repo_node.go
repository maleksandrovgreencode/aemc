package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wttech/aemc/pkg/common/fmtx"
	"github.com/wttech/aemc/pkg/common/langx"
	"github.com/wttech/aemc/pkg/common/stringsx"
	"golang.org/x/exp/maps"
)

// RepoNode represents single node in JCR repository
type RepoNode struct {
	repo Repo

	path string
}

func NewNode(repo Repo, path string) RepoNode {
	return RepoNode{
		repo: repo,
		path: path,
	}
}

func (n RepoNode) Path() string {
	return n.path
}

func (n RepoNode) Root() bool {
	return n.path == "/"
}

func (n RepoNode) Name() string {
	return stringsx.AfterLast(n.path, "/")
}

func (n RepoNode) Extension() string {
	return stringsx.AfterLast(n.Name(), ".")
}

func (n RepoNode) Content() RepoNode {
	return n.Child("jcr:content")
}

func (n RepoNode) Parent() RepoNode {
	return NewNode(n.repo, stringsx.BeforeLast(n.path, "/"))
}

func (n RepoNode) Parents() <-chan RepoNode {
	result := make(chan RepoNode)
	go func() {
		current := n
		for !current.Root() {
			current = current.Parent()
			result <- current
		}
		close(result)
	}()
	return result
}

func (n RepoNode) Child(name string) RepoNode {
	return NewNode(n.repo, fmt.Sprintf("%s/%s", n.path, name))
}

func (n RepoNode) Children() langx.Iterator[RepoNode] {
	return RepoNodeIterator{nextNodes: func(node RepoNode) []RepoNode {
		return []RepoNode{} // TODO impl this
	}}
}

func (n RepoNode) Siblings() langx.Iterator[RepoNode] {
	return RepoNodeIterator{nextNodes: func(node RepoNode) []RepoNode {
		return []RepoNode{} // TODO impl this
	}}
}

func (n RepoNode) Sibling(name string) RepoNode {
	return n.Parent().Child(name)
}

type RepoNodeState struct {
	Path       string         `yaml:"path" json:"path"`
	Exists     bool           `yaml:"exists" json:"exists"`
	Properties map[string]any `yaml:"properties" json:"properties"`
}

func (n RepoNode) State() (*RepoNodeState, error) {
	exists, err := n.ReadExists()
	if err != nil {
		return nil, err
	}
	if !exists {
		return &RepoNodeState{
			Path:   n.path,
			Exists: false,
		}, nil
	}
	props, err := n.ReadProps()
	if err != nil {
		return nil, err
	}
	return &RepoNodeState{
		Path:       n.path,
		Exists:     true,
		Properties: props,
	}, nil
}

func (n RepoNode) ReadExists() (bool, error) {
	return n.repo.Exists(n.path)
}

func (n RepoNode) ReadProps() (map[string]any, error) {
	return n.repo.Read(n.path)
}

func (n RepoNode) Save(props map[string]any) error {
	return n.repo.Save(n.path, props)
}

func (n RepoNode) SaveWithChanged(props map[string]any) (bool, error) {
	state, err := n.State()
	if err != nil {
		return false, err
	}
	if !state.Exists {
		err = n.repo.Save(n.path, props)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	propsBefore := maps.Clone(state.Properties)
	if n.repo.PropsEqual(propsBefore, props) {
		return false, nil
	}
	err = n.repo.Save(n.path, props)
	if err != nil {
		return false, err
	}
	state, err = n.State()
	if err != nil {
		return false, err
	}
	return !n.repo.PropsEqual(propsBefore, state.Properties), nil
}

func (n RepoNode) DeleteWithChanged() (bool, error) {
	state, err := n.State()
	if err != nil {
		return false, err
	}
	if !state.Exists {
		return false, nil
	}
	err = n.repo.Delete(n.path)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (n RepoNode) Delete() error {
	state, err := n.State()
	if err != nil {
		return err
	}
	if !state.Exists {
		return fmt.Errorf("node '%s' cannot be deleted as it does not exist on instance '%s'", n.path, n.repo.instance.ID())
	}
	return n.repo.Delete(n.path)
}

func (n RepoNode) DeleteProp(name string) error {
	return n.SaveProp(name, nil)
}

func (n RepoNode) SaveProp(name string, value any) error {
	return n.repo.Save(n.path, map[string]any{name: value})
}

type RepoNodeIterator struct {
	nodes     langx.Stack[RepoNode]
	nextNodes func(node RepoNode) []RepoNode
}

func (i RepoNodeIterator) Next() (RepoNode, bool) {
	node, ok := i.nodes.Pop()
	if ok {
		for _, v := range i.nextNodes(node) {
			i.nodes.Push(v)
		}
		return node, true
	}
	var zero RepoNode
	return zero, false
}

func (n RepoNode) String() string {
	return fmt.Sprintf("node '%s'", n.path)
}

func (n RepoNode) MarshalJSON() ([]byte, error) {
	state, err := n.State()
	if err != nil {
		return nil, err
	}
	return json.Marshal(state)
}

func (n RepoNode) MarshalYAML() (interface{}, error) {
	return n.State()
}

func (n RepoNode) MarshalText() string {
	state, err := n.State()
	if err != nil {
		return fmt.Sprintf("path '%s' state cannot be read\n", n.path)
	}
	sb := bytes.NewBufferString("")
	if state.Exists {
		sb.WriteString(fmt.Sprintf("path '%s'\n", n.path))
		sb.WriteString(fmtx.TblProps(state.Properties))
	} else {
		sb.WriteString(fmt.Sprintf("path '%s' does not exist\n", n.path))
	}
	return sb.String()
}
