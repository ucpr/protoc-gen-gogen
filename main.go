package main

import (
	"fmt"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"google.golang.org/protobuf/proto"

	"github.com/ucpr/protoc-gen-gogen/gogen"
)

func main() {
	pgs.Init(
		pgs.DebugEnv("DEBUG"),
	).RegisterModule(
		New(),
	).RegisterPostProcessor(
		pgsgo.GoFmt(),
	).Render()
}

type Module struct {
	*pgs.ModuleBase
	pgs.DebuggerCommon
	ctx pgsgo.Context
}

func New() pgs.Module {
	return &Module{ModuleBase: &pgs.ModuleBase{}}
}

func (m *Module) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.ctx = pgsgo.InitContext(c.Parameters())
}

func (m *Module) Name() string {
	return "gogen"
}

func (m *Module) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, t := range targets {
		// gfname := m.ctx.OutputPath(t).SetExt(".go").String()

		// outdir := m.Parameters().Str("outdir")
		// filename := gfname
		// if outdir != "" {
		// 	filename = filepath.Join(outdir, gfname)
		// }
		opts := t.File().Descriptor().GetOptions()
		rawOpt := proto.GetExtension(opts, gogen.E_Generate)
		opt, ok := rawOpt.(string)
		if !ok {
			panic(fmt.Errorf("unable to convert extension to string"))
		}
		fmt.Println(opt)

		// buf := new(bytes.Buffer)
		// m.OverwriteGeneratorFile(filename, buf.String())
	}

	return m.Artifacts()
}
