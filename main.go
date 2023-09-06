package main

import (
	"bytes"
	"os"
	"path/filepath"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/ucpr/protoc-gen-gogen/gogen"
)

func main() {
	optional := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	pgs.Init(
		pgs.DebugEnv("DEBUG"), pgs.SupportedFeatures(&optional),
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

func (m *Module) Execute(files map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, f := range files {
		gfname := m.ctx.OutputPath(f).SetExt(".go").String()

		outdir := m.Parameters().Str("outdir")
		filename := gfname
		if outdir != "" {
			filename = filepath.Join(outdir, gfname)
		}
		opts := f.File().Descriptor().GetOptions()
		ropt := proto.GetExtension(opts, gogen.E_GoGenerate)
		opt, ok := ropt.(string)
		if !ok {
			continue
		}

		data, err := os.ReadFile(filename)
		_ = err

		buf := new(bytes.Buffer)
		buf.WriteString("//go:generate " + opt + "\n\n")
		buf.Write(data)

		m.OverwriteGeneratorFile(filename, buf.String())
	}

	return m.Artifacts()
}
