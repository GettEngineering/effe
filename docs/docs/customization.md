## Customization

You can define custom tasks for Effe.

We want to generate for making a post request by url.

```golang
effe.BuildFlow(
    effe.Step(step1),
    mytask.POST(
        "http://example.com",
    ),
)
```

For that, you should create the function `POST` in your package mytask.

```golang
package mytask

func POST(uri string) interface{} {
    panic("implementation is not generated, run myeffe")
}
```

Write functions for:

- Loader
- Strategy
- Drawer

```golang

// Loading and validating expressions.
func LoadPostRequestComponent(effeConditionCall *ast.CallExpr, f loaders.FlowLoader) (types.Component, error) {
    return nil, nil // type of your component and error
}

// Generates code.
func GenPostRequestComponent(f strategies.FlowGen, c types.Component) (strategies.ComponentCall, error) {
     return nil, nil
}

// Generates a statement for plantuml.
func DrawPostRequestComponent(drawer.Drawer, types.Component) (drawer.ComponentStmt, error) {
    return nil, nil
}
```

and register it

```golang
settings := generator.DefaultSettigs()
strategy := strategies.NewChain(strategies.WithServiceObjectName(settings.LocalInterfaceVarname()))
err := strategy.Register("POST", testcustomization.GenPostRequestComponent)
if err != nil {
    fmt.Println(err)
    os.Exit(1)
}

loader := loaders.NewLoader(loaders.WithPackages([]string{"effe", "testcustomization"}))
err = loader.Register("POST", testcustomization.LoadPostRequestComponent)
if err != nil {
    fmt.Println(err)
    os.Exit(1)
}

d := drawer.NewDrawer()
err = d.Register("POST", DrawPostRequestComponent)
if err != nil {
    fmt.Println(err)
    os.Exit(1)
}

gen := generator.NewGenerator(
    generator.WithSetttings(settings),
    generator.WithLoader(loader),
    generator.WithStrategy(strategy),
    generator.WithDrawer(d),
)

// Run generator
gen.Generate(context.Background(), d, os.Environ(), []string{"."}
```
