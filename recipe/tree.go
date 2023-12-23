package recipe

func LoadTree(rootPath string) (*RecipeTree, error) {
	t := &RecipeTree{
		recipes:        map[string]*Recipe{},
		traversalOrder: []string{},
	}

	r, err := t.loadRecipe(rootPath)
	if err != nil {
		return nil, err
	}

	t.Root = r
	return t, nil
}

func (t *RecipeTree) loadRecipe(path string) (*Recipe, error) {
	r, ok := t.recipes[path]
	if ok {
		return r, nil
	}

	r, err := loadRecipe(path)
	if err != nil {
		return nil, err
	}
	t.recipes[path] = r

	for _, layer := range r.Layers {
		_, err := t.loadRecipe(layer)
		if err != nil {
			return nil, err
		}

		t.linkLayer(r, t.recipes[layer])
	}

	t.traversalOrder = append(t.traversalOrder, path)
	return r, nil
}

func (t *RecipeTree) linkLayer(parent *Recipe, layer *Recipe) {
	switch layer.TargetType {
	case TARGET_STATIC_LIBRARY:
		parent.LinkedStaticLayers = append(parent.LinkedStaticLayers, layer.Target)
	case TARGET_SHARED_LIBRARY:
		parent.LinkedSharedLayers = append(parent.LinkedSharedLayers, layer.Target)
	}
}

type Visitor interface {
	Visit(r *Recipe, index int) bool
}

type RecipeTree struct {
	Root           *Recipe
	recipes        map[string]*Recipe
	traversalOrder []string
}

func (t RecipeTree) Traverse(v Visitor) bool {
	for index, key := range t.traversalOrder {
		if !v.Visit(t.recipes[key], index) {
			return false
		}
	}
	return true
}
