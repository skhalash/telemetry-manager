package dryrun

import (
	"context"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	telemetryv1alpha1 "github.com/kyma-project/telemetry-manager/apis/telemetry/v1alpha1"
	"github.com/kyma-project/telemetry-manager/internal/fluentbit/config/builder"
	"github.com/kyma-project/telemetry-manager/internal/resources/fluentbit"
)

type fileWriter interface {
	PreparePipelineDryRun(ctx context.Context, workDir string, pipeline *telemetryv1alpha1.LogPipeline) (func(), error)
	PrepareParserDryRun(ctx context.Context, workDir string, pipeline *telemetryv1alpha1.LogParser) (func(), error)
}

type fileWriterImpl struct {
	client client.Client
	config Config
}

func (f *fileWriterImpl) PrepareParserDryRun(ctx context.Context, workDir string, parser *telemetryv1alpha1.LogParser) (func(), error) {
	if err := makeDir(workDir); err != nil {
		return nil, err
	}
	if err := f.writeParsersWithParser(ctx, workDir, parser); err != nil {
		return nil, err
	}

	return func() { deleteWorkDir(ctx, workDir) }, nil
}

func (f *fileWriterImpl) PreparePipelineDryRun(ctx context.Context, workDir string, pipeline *telemetryv1alpha1.LogPipeline) (func(), error) {
	if err := makeDir(workDir); err != nil {
		return nil, err
	}
	if err := f.writeConfig(ctx, workDir); err != nil {
		return nil, err
	}
	if err := f.writeFiles(pipeline, workDir); err != nil {
		return nil, err
	}
	if err := f.writeSections(pipeline, workDir); err != nil {
		return nil, err
	}
	if err := f.writeParsers(ctx, workDir); err != nil {
		return nil, err
	}

	return func() { deleteWorkDir(ctx, workDir) }, nil
}

func (f *fileWriterImpl) writeConfig(ctx context.Context, basePath string) error {
	var cm corev1.ConfigMap
	var err error
	err = f.client.Get(ctx, f.config.FluentBitConfigMapName, &cm)
	if err != nil {
		if apierrors.IsNotFound(err) {
			newCm := fluentbit.MakeConfigMap(f.config.FluentBitConfigMapName)
			cm = *newCm
		} else {
			return err
		}
	}

	for key, data := range cm.Data {
		err = writeFile(filepath.Join(basePath, key), data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *fileWriterImpl) writeFiles(pipeline *telemetryv1alpha1.LogPipeline, basePath string) error {
	filesDir := filepath.Join(basePath, "files")
	if err := makeDir(filesDir); err != nil {
		return err
	}

	for _, file := range pipeline.Spec.Files {
		err := writeFile(filepath.Join(filesDir, file.Name), file.Content)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *fileWriterImpl) writeSections(pipeline *telemetryv1alpha1.LogPipeline, basePath string) error {
	dynamicDir := filepath.Join(basePath, "dynamic")
	if err := makeDir(dynamicDir); err != nil {
		return err
	}

	builderConfig := builder.BuilderConfig{
		PipelineDefaults: f.config.PipelineDefaults,
	}
	sectionsConfig, err := builder.BuildFluentBitConfig(pipeline, builderConfig)
	if err != nil {
		return err
	}

	return writeFile(filepath.Join(dynamicDir, pipeline.Name+".conf"), sectionsConfig)
}

func (f *fileWriterImpl) writeParsers(ctx context.Context, basePath string) error {
	dynamicParsersDir := filepath.Join(basePath, "dynamic-parsers")
	if err := makeDir(dynamicParsersDir); err != nil {
		return err
	}

	var logParsers telemetryv1alpha1.LogParserList
	if err := f.client.List(ctx, &logParsers); err != nil {
		return err
	}

	parsersConfig := builder.BuildFluentBitParsersConfig(&logParsers)
	return writeFile(filepath.Join(dynamicParsersDir, "parsers.conf"), parsersConfig)
}

func (f *fileWriterImpl) writeParsersWithParser(ctx context.Context, basePath string, parser *telemetryv1alpha1.LogParser) error {
	dynamicParsersDir := filepath.Join(basePath, "dynamic-parsers")
	if err := makeDir(dynamicParsersDir); err != nil {
		return err
	}

	var logParsers telemetryv1alpha1.LogParserList
	if err := f.client.List(ctx, &logParsers); err != nil {
		return err
	}

	appendOrReplace(&logParsers, parser)
	parsersConfig := builder.BuildFluentBitParsersConfig(&logParsers)

	return writeFile(filepath.Join(dynamicParsersDir, "parsers.conf"), parsersConfig)
}

func appendOrReplace(parsers *telemetryv1alpha1.LogParserList, parser *telemetryv1alpha1.LogParser) {
	for i := range parsers.Items {
		if parsers.Items[i].Name == parser.Name {
			parsers.Items[i] = *parser
			return
		}
	}
	parsers.Items = append(parsers.Items, *parser)
}

func deleteWorkDir(ctx context.Context, workDir string) {
	if err := os.RemoveAll(workDir); err != nil {
		log := logf.FromContext(ctx)
		log.Error(err, "Failed to remove Fluent Bit config directory")
	}
}

func makeDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func writeFile(path string, data string) error {
	return os.WriteFile(path, []byte(data), 0600)
}
